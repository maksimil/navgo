package main

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/maksimil/navgo/pkg/dterm"
)

const (
	PathTreeErr = iota
	PathTreeClosed
	PathTreeOpen
)

type PathTreePart interface {
	Path() string
	Get(idx []int) PathTreePart
	Draw(target *dterm.THandle, highlight []int)
	Open() bool
	Close() bool
}

type PathLeaf struct {
	path string
}

type PathTree struct {
	path     string
	state    int
	children []PathTreePart
}

func (part *PathTree) Path() string {
	return part.path
}

func (part *PathTree) Get(idx []int) PathTreePart {
	if len(idx) == 0 {
		return part
	} else {
		return part.children[idx[0]].Get(idx[1:])
	}
}

func (part *PathTree) Draw(target *dterm.THandle, highlight []int) {
	dname := path.Base(part.path)
	if len(highlight) == 0 {
		dname = fmt.Sprintf("\x1b[104;30m%s\x1b[0m", dname)
	} else {
		dname = fmt.Sprintf("\x1b[40;94m%s\x1b[0m", dname)
	}
	switch part.state {
	case PathTreeClosed:
		target.PutLinef("%s", dname)
		target.MoveBy(0, 1)
	case PathTreeErr:
		target.PutLinef("%s \x1b[31m- Error\x1b[0m", dname)
		target.MoveBy(0, 1)
	case PathTreeOpen:
		target.PutLinef("%s", dname)
		target.MoveBy(0, 1)
		for i, cpart := range part.children {
			if i+1 == len(part.children) {
				target.PutLine("\u2514\u2500 ")
			} else {
				target.PutLine("\u251C\u2500 ")
			}
			target.MoveBy(3, 0)
			x, y := target.Cpos()
			if len(highlight) == 0 || i != highlight[0] {
				cpart.Draw(target, []int{-1})
			} else {
				cpart.Draw(target, highlight[1:])
			}
			_, y1 := target.Cpos()
			if i+1 != len(part.children) {
				target.MoveTo(x-3, y+1)
				for j := 0; j < y1-y-1; j++ {
					target.PutLine("\u2502")
					target.MoveBy(0, 1)
				}
			} else {
				target.MoveBy(-3, 0)
			}
		}
	}
}

func (part *PathTree) Open() bool {
	if part.state == PathTreeClosed {
		part.state = PathTreeOpen
		de, err := os.ReadDir(part.path)
		if err != nil {
			panic(err)
		}

		part.children = make([]PathTreePart, len(de))

		for i := 0; i < len(de); i++ {
			if de[i].IsDir() {
				part.children[i] = &PathTree{
					path.Join(part.path, de[i].Name()),
					PathTreeClosed, make([]PathTreePart, 0)}
			} else {
				part.children[i] = &PathLeaf{path.Join(part.path, de[i].Name())}
			}
		}

		sort.Slice(part.children, func(i, j int) bool {
			switch ci := part.children[i].(type) {
			case *PathLeaf:
				switch cj := part.children[j].(type) {
				case *PathLeaf:
					return ci.path < cj.path
				case *PathTree:
					return false
				}
			case *PathTree:
				switch cj := part.children[j].(type) {
				case *PathLeaf:
					return true
				case *PathTree:
					return ci.path < cj.path
				}
			}
			return true
		})
	}
	return true
}

func (part *PathTree) Close() bool {
	if part.state == PathTreeOpen {
		part.state = PathTreeClosed
		part.children = nil
		return true
	} else {
		return false
	}
}

func (part *PathLeaf) Path() string {
	return part.path
}

func (part *PathLeaf) Get(idx []int) PathTreePart {
	if len(idx) == 0 {
		return part
	} else {
		panic("Index out of bounds")
	}
}

func (part *PathLeaf) Draw(target *dterm.THandle, highlight []int) {
	dname := path.Base(part.path)
	if len(highlight) == 0 {
		dname = fmt.Sprintf("\x1b[2;47;30m%s\x1b[0m", dname)
	}
	target.PutLine(dname)
	target.MoveBy(0, 1)
}

func (part *PathLeaf) Open() bool {
	return false
}

func (part *PathLeaf) Close() bool {
	return false
}
