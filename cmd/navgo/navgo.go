package main

import (
	"time"

	"github.com/maksimil/navgo/pkg/dterm"
)

func main() {
	th :=
		dterm.NewTHandle()
	defer th.Close("b msg")
	th.MoveTo(0, 2)
	th.PutLine("Hi")
	th.MoveTo(0, 0)
	th.PutLine("Bye")
	time.Sleep(time.Second)
}
