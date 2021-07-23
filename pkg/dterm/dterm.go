package dterm

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type THandle struct {
	cx     int
	cy     int
	height int
	stream io.Writer

	limit int
	lock  int
}

func (handle *THandle) Cpos() (int, int) {
	return handle.cx, handle.cy
}

func (handle *THandle) act(y int) int {
	if y >= handle.lock && y < handle.limit {
		return y
	} else if y >= handle.lock {
		return handle.limit - 1
	} else {
		return handle.lock
	}
}

func (handle *THandle) acty() int {
	return handle.act(handle.cy)
}

func NewTHandle() THandle {
	_, limit, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	return NewTHandleStreamed(os.Stdout, limit-1)
}

func NewTHandleStreamed(stream io.Writer, limit int) THandle {
	sh := THandle{0, 0, 1, stream, limit, 0}
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
	if handle.height+by > handle.limit {
		handle.Expand(handle.limit - handle.height)
	} else {
		// move to the end of the avaliable height
		handle.moveby_raw(-handle.cx, handle.height-handle.acty())
		for i := 0; i < by; i++ {
			handle.Write("\n")
		}
		handle.moveby_raw(handle.cx, -(handle.height-handle.acty())-by)
		handle.height += by
	}
}

func (handle *THandle) MoveBy(x, y int) {
	if handle.cy+y >= handle.height {
		handle.Expand(handle.cy + y - handle.height + 1)
	}
	handle.cx += x
	handle.cy += y

	dy := handle.act(handle.cy) - handle.act(handle.cy-y)

	handle.moveby_raw(x, dy)
}

func (handle *THandle) MoveTo(x, y int) {
	handle.MoveBy(x-handle.cx, y-handle.cy)
}

func (handle *THandle) PutLine(line string) {
	if handle.cy < handle.limit && handle.cy >= handle.lock {
		handle.Writef("\x1b7%s\x1b8", line)
	}
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

func (handle *THandle) LockOffset(offset int) {
	handle.lock = handle.cy
	handle.MoveBy(0, offset)
}

func (handle *THandle) Unlock() {
	if handle.cy < handle.lock {
		handle.MoveTo(handle.cx, handle.lock)
	}
	handle.lock = 0
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
