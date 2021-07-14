package dterm

import (
	"fmt"
	"io"
	"os"
)

type THandle struct {
	cx     int
	cy     int
	height int
	stream io.Writer
}

func (handle *THandle) Cpos() (int, int) {
	return handle.cx, handle.cy
}

func NewTHandle() THandle {
	return NewTHandleStreamed(os.Stdout)
}

func NewTHandleStreamed(stream io.Writer) THandle {
	sh := THandle{0, 0, 1, stream}
	sh.Write("\n\x1b[0A")
	return sh
}

func (handle *THandle) Write(s string) {
	fmt.Fprint(handle.stream, s)
}

func (handle *THandle) Writef(format string, a ...interface{}) {
	fmt.Fprintf(handle.stream, format, a...)
}

func (handle *THandle) moveby_raw(x, y int) {
	if x > 0 {
		handle.Writef("\x1b[%dC", x)
	} else if x != 0 {
		handle.Writef("\x1b[%dD", -x)
	}
	if y > 0 {
		handle.Writef("\x1b[%dB", y)
	} else if y != 0 {
		handle.Writef("\x1b[%dA", -y)
	}
}

func (handle *THandle) Expand(by int) {
	// move to the end of the avaliable height
	handle.moveby_raw(-handle.cx, handle.height-handle.cy)
	for i := 0; i < by; i++ {
		handle.Write("\n")
	}
	handle.moveby_raw(handle.cx, -(handle.height-handle.cy)-by)
	handle.height += by
}

func (handle *THandle) MoveBy(x, y int) {
	if handle.cy+y >= handle.height {
		handle.Expand(handle.cy + y - handle.height + 1)
	}
	handle.cx += x
	handle.cy += y

	handle.moveby_raw(x, y)
}

func (handle *THandle) MoveTo(x, y int) {
	handle.MoveBy(x-handle.cx, y-handle.cy)
}

func (handle *THandle) PutLine(line string) {
	handle.Writef("\x1b7%s\x1b8", line)
}

func (handle *THandle) PutLinef(format string, a ...interface{}) {
	handle.PutLine(fmt.Sprintf(format, a...))
}

func (handle *THandle) Clear() {
	handle.MoveTo(0, 0)
	handle.Write("\x1b[0J")
}

func (handle *THandle) HideCursor() {
	handle.Write("\x1b[?25l")
}

func (handle *THandle) ShowCursor() {
	handle.Write("\x1b[?25h")
}

func (handle *THandle) Close(exmsg string) {
	handle.Clear()
	handle.ShowCursor()
	handle.Writef("%s", exmsg)

}

func (handle *THandle) CloseDirty() {
	handle.MoveTo(0, handle.height)
	handle.ShowCursor()
}
