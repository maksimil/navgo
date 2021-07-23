package main

import (
	"os"
	"strings"

	"github.com/maksimil/navgo/pkg/dterm"
	"golang.org/x/term"
)

type UIState interface {
	Call(c byte) (bool, UIState)
	Draw(target *dterm.THandle)
}

// h j k l d f i
var KEYS = map[byte]interface{}{
	104: 0, 106: 0, 107: 0, 108: 0,
	100: 0, 102: 0, 105: 0,
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
	uistatechan := make(chan byte, 16)

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
			_, includes := KEYS[c]
			if includes {
				uistatechan <- c
			}
		}
	}()

	// ui drawing goroutine
	go func() {
		// initial state
		uistate := UIState(new(UITreeState))
		switch state := uistate.(type) {
		case *UITreeState:
			state.tree.path = dir
			state.tree.state = PathTreeClosed
		}
		// initial draw
		uistate.Draw(&th)
		// loop
		for c := range uistatechan {
			// mutating
			var muted bool
			muted, uistate = uistate.Call(c)
			for len(uistatechan) > 0 {
				var m bool
				m, uistate = uistate.Call(<-uistatechan)
				muted = m || muted
			}
			// drawing
			if muted {
				uistate.Draw(&th)
			}
		}
	}()

	(<-closechan)()
}
