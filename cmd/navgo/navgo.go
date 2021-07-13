package main

import (
	"os"
	"path"
	"strings"

	"github.com/maksimil/navgo/pkg/dterm"
	"golang.org/x/term"
)

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
	// drawing ui
	th := dterm.NewTHandle()
	th.PutLinef("\x1b[33m%s\x1b[0m", dir)
	th.MoveBy(0, 1)
	th.PutLinef("%s", path.Base(dir))

	closechan := make(chan string)

	// raw user input goroutine
	go func() {
		for {
			b := make([]byte, 1)
			os.Stdin.Read(b)
			c := b[0]
			switch {
			case c == 3:
				closechan <- "\x1b[33mExit via ^C\x1b[31m"
				close(closechan)
			}
		}
	}()

	th.Close(<-closechan)
}
