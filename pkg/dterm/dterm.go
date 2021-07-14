package dterm

import (
	"fmt"
)

type THandle struct {
	cx     int
	cy     int
	height int
}

func (handle *THandle) Cpos() (int, int) {
	return handle.cx, handle.cy
}

func NewTHandle() THandle {
	fmt.Print("\n\x1b[0A")
	sh := THandle{0, 0, 1}
	return sh
}

func moveby(x, y int) {
	if x > 0 {
		fmt.Printf("\x1b[%dC", x)
	} else if x != 0 {
		fmt.Printf("\x1b[%dD", -x)
	}
	if y > 0 {
		fmt.Printf("\x1b[%dB", y)
	} else if y != 0 {
		fmt.Printf("\x1b[%dA", -y)
	}
}

func (handle *THandle) Expand(by int) {
	// move to the end of the avaliable height
	moveby(-handle.cx, handle.height-handle.cy)
	for i := 0; i < by; i++ {
		fmt.Print("\n")
	}
	moveby(handle.cx, -(handle.height-handle.cy)-by)
	handle.height += by
}

func (handle *THandle) MoveBy(x, y int) {
	if handle.cy+y >= handle.height {
		handle.Expand(handle.cy + y - handle.height + 1)
	}
	handle.cx += x
	handle.cy += y

	moveby(x, y)
}

func (handle *THandle) MoveTo(x, y int) {
	handle.MoveBy(x-handle.cx, y-handle.cy)
}

func PutLine(line string) {
	fmt.Printf("\x1b7%s\x1b8", line)
}

func PutLinef(format string, a ...interface{}) {
	PutLine(fmt.Sprintf(format, a...))
}

func (handle *THandle) Clear() {
	handle.MoveTo(0, 0)
	fmt.Print("\x1b[0J")
}

func HideCursor() {
	fmt.Print("\x1b[?25l")
}

func ShowCursor() {
	fmt.Print("\x1b[?25h")
}

func (handle *THandle) Close(exmsg string) {
	handle.Clear()
	ShowCursor()
	fmt.Printf("%s", exmsg)

}

func (handle *THandle) CloseDirty() {
	handle.MoveTo(0, handle.height)
	ShowCursor()
}
