package dterm

import (
	"bytes"
)

type Buffer struct {
	handle    *THandle
	bufhandle THandle
	buf       *bytes.Buffer
}

func NewBuffer(handle *THandle) Buffer {
	buf := new(bytes.Buffer)
	bufhandle := NewTHandleStreamed(buf, handle.limit)
	bufhandle.cx = handle.cx
	bufhandle.cy = handle.cy
	bufhandle.height = handle.height
	return Buffer{handle, bufhandle, buf}
}

func (buffer *Buffer) Handle() *THandle {
	return &buffer.bufhandle
}

func (buffer Buffer) Write() {
	buffer.handle.Write(buffer.buf.String())
	buffer.handle.cx = buffer.bufhandle.cx
	buffer.handle.cy = buffer.bufhandle.cy
	buffer.handle.height = buffer.bufhandle.height
}

func (handle *THandle) Bufferize(draw func(*THandle)) {
	buf := NewBuffer(handle)
	draw(buf.Handle())
	buf.Write()
}
