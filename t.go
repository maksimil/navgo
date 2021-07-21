package main

import (
	"time"

	"github.com/maksimil/navgo/pkg/dterm"
)

func main() {
	th := dterm.NewTHandle()
	th.MoveBy(0, -1)
	th.PutLine("Hi")
	time.Sleep(time.Second)
	th.Close("Exited")
}
