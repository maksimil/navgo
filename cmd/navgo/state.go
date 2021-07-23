package main

import "github.com/maksimil/navgo/pkg/dterm"

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
		if state.tree.Get(state.selected).Open() {
			switch part := state.tree.Get(state.selected).(type) {
			case *PathTree:
				if len(part.children) > 0 {
					state.selected = append(state.selected, 0)
				}
			}
			return true, state
		} else {
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
