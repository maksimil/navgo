package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/maksimil/navgo/pkg/dterm"
	"golang.org/x/term"
)

type UIState struct {
	tree     PathTree
	selected []int
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
	dterm.HideCursor()

	closechan := make(chan string)
	uistatechan := make(chan func(*UIState) bool, 16)

	uistatechan <- func(u *UIState) bool { return true }

	// raw user input goroutine
	go func() {
		for {
			b := make([]byte, 1)
			os.Stdin.Read(b)
			c := b[0]
			os.Stderr.WriteString(fmt.Sprintln(c))
			switch {
			case c == 3:
				closechan <- "\x1b[33mExit via ^C\x1b[31m"
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
			// k
			case c == 107:
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
			}
		}
	}()

	// ui drawing goroutine
	go func() {
		uistate := UIState{PathTree{dir, PathTreeClosed, []PathTreePart{}}, []int{}}
		for mutator := range uistatechan {
			if mutator(&uistate) {
				// drawing
				th.Clear()
				dterm.PutLinef("\x1b[33m%s\x1b[0m", uistate.tree.path)
				th.MoveBy(0, 1)
				uistate.tree.Draw(&th, uistate.selected)
			}
		}
	}()

	th.Close(<-closechan)
}
