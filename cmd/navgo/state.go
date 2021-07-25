package main

import (
	"io/ioutil"
	"strings"

	"github.com/maksimil/navgo/pkg/dterm"
)

type UIState interface {
	Call(c byte) (bool, UIState)
	Draw(target *dterm.THandle)
}

type UITreeState struct {
	tree     PathTree
	selected []int
	scroll   int
}

func (state *UITreeState) Call(c byte) (bool, UIState) {
	switch c {
	// h
	case 104:
		if len(state.selected) > 0 {
			state.selected = state.selected[:len(state.selected)-1]
			return true, state
		} else {
			return false, state
		}

	// j
	case 106:
		if len(state.selected) == 0 {
			return false, state
		} else {
			state.selected[len(state.selected)-1] += 1
			state.selected[len(state.selected)-1] = modulo(state.selected[len(state.selected)-1],
				len(
					state.tree.Get(state.selected[:len(state.selected)-1]).(*PathTree).children))
			return true, state
		}

	// k
	case 107:
		if len(state.selected) == 0 {
			return false, state
		} else {
			state.selected[len(state.selected)-1] -= 1
			state.selected[len(state.selected)-1] = modulo(state.selected[len(state.selected)-1],
				len(
					state.tree.Get(state.selected[:len(state.selected)-1]).(*PathTree).children))
			return true, state

		}

	// l
	case 108:
		switch tp := state.tree.Get(state.selected).(type) {
		case *PathTree:
			tp.Open()
			if len(tp.children) > 0 {
				state.selected = append(state.selected, 0)
			}
			return true, state
		case *PathLeaf:
			// opening cat
			cat := OpenCat(tp.path, state)
			return true, &cat
		default:
			return false, state
		}

	// d
	case 100:
		if state.scroll > 0 {
			state.scroll -= 1
			return true, state
		} else {
			return false, state
		}

	// f
	case 102:
		state.scroll += 1
		return true, state

	// i
	case 105:
		closed := state.tree.Get(state.selected).Close()
		return closed, state

	// default
	default:
		return false, state
	}
}

func (state *UITreeState) Draw(target *dterm.THandle) {
	target.Bufferize(func(handle *dterm.THandle) {
		handle.Clear()
		handle.PutLinef("\x1b[43;30m%d\x1b[0m \x1b[33m%s\x1b[0m",
			state.scroll, state.tree.path)
		handle.MoveBy(0, 1)
		handle.LockOffset(-state.scroll)
		state.tree.Draw(handle, state.selected)
		handle.Unlock()
	})
}

type UICatState struct {
	tree   *UITreeState
	lines  []string
	scroll int
	path   string
}

func OpenCat(path string, tree *UITreeState) UICatState {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	return UICatState{tree, lines, 0, path}
}

func (state *UICatState) Call(c byte) (bool, UIState) {
	switch c {
	// h, i
	case 104, 105:
		return true, state.tree

	// j, f
	case 106, 102:
		state.scroll += 1
		return true, state

	// k, d
	case 107, 100:
		if state.scroll > 0 {
			state.scroll -= 1
			return true, state
		} else {
			return false, state
		}

	// l
	case 108:
		return false, state

	// default
	default:
		return false, state
	}
}

func (state *UICatState) Draw(target *dterm.THandle) {
	target.Bufferize(func(handle *dterm.THandle) {
		handle.Clear()
		handle.PutLinef("\x1b[43;30m%d\x1b[0m \x1b[33m%s\x1b[0m",
			state.scroll, state.path)
		handle.MoveBy(0, 1)
		for i := state.scroll; i < len(state.lines); i++ {
			handle.PutLine(state.lines[i])
			handle.MoveBy(0, 1)
		}
	})
}
