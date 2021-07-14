package main

import (
	"os"
	"strings"

	"github.com/maksimil/navgo/pkg/dterm"
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
	pt := &PathTree{dir, PathTreeClosed, make([]PathTreePart, 0)}
	open(pt)
	open(pt.children[0].(*PathTree))
	pt.children[0].(*PathTree).children[5].(*PathTree).state = PathTreeErr
	th := dterm.NewTHandle()
	drawPart(pt, &th, []int{0})
	th.CloseDirty()
	// // putting terminal in raw mode
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer term.Restore(int(os.Stdin.Fd()), oldState)
	// // drawing ui
	// th := dterm.NewTHandle()
	// th.PutLinef("\x1b[33m%s\x1b[0m", dir)
	// th.MoveBy(0, 1)
	// th.PutLinef("%s", path.Base(dir))

	// closechan := make(chan string)
	// uistatechan := make(chan func(*UIState), 16)

	// // raw user input goroutine
	// go func() {
	// 	for {
	// 		b := make([]byte, 1)
	// 		os.Stdin.Read(b)
	// 		c := b[0]
	// 		switch {
	// 		case c == 3:
	// 			closechan <- "\x1b[33mExit via ^C\x1b[31m"
	// 			close(closechan)
	// 		}
	// 	}
	// }()

	// // ui drawing goroutine
	// go func() {
	// 	uistate := UIState{}
	// 	for mutator := range uistatechan {
	// 		mutator(&uistate)
	// 		os.Stderr.WriteString(fmt.Sprintf("%v\n", uistate))
	// 	}
	// }()

	// th.Close(<-closechan)
}
