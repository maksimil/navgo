package main

import (
	"os"
	"strings"

	"github.com/maksimil/navgo/pkg/dterm"
	"golang.org/x/term"
)

type UIState struct {
	tree     PathTree
	selected []int
	scroll   int
}

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
	uistatechan := make(chan func(*UIState) bool, 16)

	uistatechan <- func(u *UIState) bool { return true }

	// raw user input goroutine
	go func() {
		for {
			b := make([]byte, 1)
			os.Stdin.Read(b)
			c := b[0]
			// fmt.Fprintf(os.Stderr, "%d\n", c)
			switch {
			case c == 3:
				closechan <- func() {
					th.Close("\x1b[33mExit via ^C\x1b[31m")
				}
				close(closechan)
			// h
			case c == 104:
				uistatechan <- func(u *UIState) bool {
					if len(u.selected) > 0 {
						u.selected = u.selected[:len(u.selected)-1]
						return true
					} else {
						return false
					}
				}
			// j
			case c == 106:
				uistatechan <- func(u *UIState) bool {
					if len(u.selected) == 0 {
						return false
					} else {
						u.selected[len(u.selected)-1] += 1
						u.selected[len(u.selected)-1] = modulo(u.selected[len(u.selected)-1],
							len(
								u.tree.Get(u.selected[:len(u.selected)-1]).(*PathTree).children))
						return true
					}
				}
			// k
			case c == 107:
				uistatechan <- func(u *UIState) bool {
					if len(u.selected) == 0 {
						return false
					} else {
						u.selected[len(u.selected)-1] -= 1
						u.selected[len(u.selected)-1] = modulo(u.selected[len(u.selected)-1],
							len(
								u.tree.Get(u.selected[:len(u.selected)-1]).(*PathTree).children))
						return true
					}
				}
			// l
			case c == 108:
				uistatechan <- func(u *UIState) bool {
					if u.tree.Get(u.selected).Open() {
						u.selected = append(u.selected, 0)
						return true
					} else {
						return false
					}
				}
			// d
			case c == 100:
				uistatechan <- func(u *UIState) bool {
					if u.scroll > 0 {
						u.scroll -= 1
						return true
					} else {
						return false
					}
				}
			// f
			case c == 102:
				uistatechan <- func(u *UIState) bool {
					u.scroll += 1
					return true
				}
			}
		}
	}()

	// ui drawing goroutine
	go func() {
		uistate := UIState{PathTree{dir, PathTreeClosed, []PathTreePart{}}, []int{}, 0}
		for mutator := range uistatechan {
			if mutator(&uistate) {
				// drawing
				th.Bufferize(func(handle *dterm.THandle) {
					handle.Clear()
					handle.PutLinef("\x1b[43;30m%d\x1b[0m \x1b[33m%s\x1b[0m",
						uistate.scroll, uistate.tree.path)
					handle.MoveBy(0, 1)
					handle.LockOffset(-uistate.scroll)
					uistate.tree.Draw(handle, uistate.selected)
					handle.Unlock()
				})
			}
		}
	}()

	(<-closechan)()
}
