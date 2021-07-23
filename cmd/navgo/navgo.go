package main

import (
	"os"
	"strings"

	"github.com/maksimil/navgo/pkg/dterm"
	"golang.org/x/term"
)

type UITreeState struct {
	tree     PathTree
	selected []int
	scroll   int
}

type UIState interface {
	isUIState()
}

func (state *UITreeState) isUIState() {}

func modulo(a, b int) int {
	return ((a % b) + b) % b
}

func main() {
	// getting current directory
	dir := func() string {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dir = strings.Replace(dir, "\\", "/", -1)
		return dir
	}()
	// putting terminal in raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// create dterm handle
	th := dterm.NewTHandle()
	th.HideCursor()

	closechan := make(chan func())
	uistatechan := make(chan func(UIState) (bool, UIState), 16)

	uistatechan <- func(u UIState) (bool, UIState) { return true, u }

	// raw user input goroutine
	go func() {
		for {
			b := make([]byte, 1)
			os.Stdin.Read(b)
			c := b[0]
			// fmt.Fprintf(os.Stderr, "%d\n", c)
			if c == 3 {
				closechan <- func() {
					th.Close("\x1b[33mExit via ^C\x1b[31m")
				}
				close(closechan)
			}
			switch {
			// h
			case c == 104:
				uistatechan <- func(u UIState) (bool, UIState) {
					switch u := u.(type) {
					case *UITreeState:
						if len(u.selected) > 0 {
							u.selected = u.selected[:len(u.selected)-1]
							return true, u
						} else {
							return false, u
						}
					default:
						return false, u
					}
				}
			// j
			case c == 106:
				uistatechan <- func(u UIState) (bool, UIState) {
					switch u := u.(type) {
					case *UITreeState:
						if len(u.selected) == 0 {
							return false, u
						} else {
							u.selected[len(u.selected)-1] += 1
							u.selected[len(u.selected)-1] = modulo(u.selected[len(u.selected)-1],
								len(
									u.tree.Get(u.selected[:len(u.selected)-1]).(*PathTree).children))
							return true, u
						}
					default:
						return false, u
					}
				}
			// k
			case c == 107:
				uistatechan <- func(u UIState) (bool, UIState) {
					switch u := u.(type) {
					case *UITreeState:
						if len(u.selected) == 0 {
							return false, u
						} else {
							u.selected[len(u.selected)-1] -= 1
							u.selected[len(u.selected)-1] = modulo(u.selected[len(u.selected)-1],
								len(
									u.tree.Get(u.selected[:len(u.selected)-1]).(*PathTree).children))
							return true, u
						}
					default:
						return false, u
					}
				}
			// l
			case c == 108:
				uistatechan <- func(u UIState) (bool, UIState) {
					switch u := u.(type) {
					case *UITreeState:
						if u.tree.Get(u.selected).Open() {
							switch part := u.tree.Get(u.selected).(type) {
							case *PathTree:
								if len(part.children) > 0 {
									u.selected = append(u.selected, 0)
								}
							}
							return true, u
						} else {
							return false, u
						}
					default:
						return false, u
					}
				}
			// d
			case c == 100:
				uistatechan <- func(u UIState) (bool, UIState) {
					switch u := u.(type) {
					case *UITreeState:
						if u.scroll > 0 {
							u.scroll -= 1
							return true, u
						} else {
							return false, u
						}
					default:
						return false, u
					}
				}
			// f
			case c == 102:
				uistatechan <- func(u UIState) (bool, UIState) {
					switch u := u.(type) {
					case *UITreeState:
						u.scroll += 1
						return true, u
					default:
						return false, u
					}
				}
			// i
			case c == 105:
				uistatechan <- func(u UIState) (bool, UIState) {
					switch u := u.(type) {
					case *UITreeState:
						closed := u.tree.Get(u.selected).Close()
						return closed, u
					default:
						return false, u
					}
				}
			}
		}
	}()

	// ui drawing goroutine
	go func() {
		uistate := UIState(new(UITreeState))
		switch state := uistate.(type) {
		case *UITreeState:
			state.tree.path = dir
			state.tree.state = PathTreeClosed

		}
		for mutator := range uistatechan {
			muted, uistate := mutator(uistate)
			for len(uistatechan) > 0 {
				var m bool
				m, uistate = (<-uistatechan)(uistate)
				muted = m || muted
			}
			if muted {
				// drawing
				th.Bufferize(func(handle *dterm.THandle) {
					switch state := uistate.(type) {
					case *UITreeState:
						handle.Clear()
						handle.PutLinef("\x1b[43;30m%d\x1b[0m \x1b[33m%s\x1b[0m",
							state.scroll, state.tree.path)
						handle.MoveBy(0, 1)
						handle.LockOffset(-state.scroll)
						state.tree.Draw(handle, state.selected)
						handle.Unlock()
					}
				})
			}
		}
	}()

	(<-closechan)()
}
