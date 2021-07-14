package main

import (
	"fmt"
	"os"
	"path"

	"github.com/maksimil/navgo/pkg/dterm"
)

const (
	PathTreeErr = iota
	PathTreeClosed
	PathTreeOpen
)

type PathTreePart interface{ isPathTreePart() }

type PathLeaf struct {
	path string
}

type PathTree struct {
	path     string
	state    int
	children []PathTreePart
}

func (*PathLeaf) isPathTreePart() {}
func (*PathTree) isPathTreePart() {}

func open(tree *PathTree) {
	if tree.state == PathTreeClosed {
		tree.state = PathTreeOpen
		de, err := os.ReadDir(tree.path)
		if err != nil {
			panic(err)
		}

		tree.children = make([]PathTreePart, len(de))

		for i := 0; i < len(de); i++ {
			if de[i].IsDir() {
				tree.children[i] = &PathTree{
					path.Join(tree.path, de[i].Name()),
					PathTreeClosed, make([]PathTreePart, 0)}
			} else {
				tree.children[i] = &PathLeaf{path.Join(tree.path, de[i].Name())}
			}
		}
	}
}

func drawPart(part PathTreePart, target *dterm.THandle, highlight []int) {
	highlighted := len(highlight) == 0
	switch part := part.(type) {
	case *PathLeaf:
		if highlighted {
			target.PutLinef("\x1b[47;30m%s\x1b[0m", path.Base(part.path))
		} else {
			target.PutLinef("%s", path.Base(part.path))
		}
		target.MoveBy(0, 1)
	case *PathTree:
		displayname := path.Base(part.path)
		if highlighted {
			displayname = fmt.Sprintf("\x1b[47;30m%s\x1b[0m", displayname)
		}
		switch part.state {
		case PathTreeClosed:
			target.PutLinef("%s", displayname)
			target.MoveBy(0, 1)
		case PathTreeErr:
			target.PutLinef("%s \x1b[31m- Error\x1b[0m", displayname)
			target.MoveBy(0, 1)
		case PathTreeOpen:
			target.PutLinef("%s", displayname)
			target.MoveBy(0, 1)
			for i, cpart := range part.children {
				if i+1 == len(part.children) {
					target.PutLine("\u2514\u2500 ")
				} else {
					target.PutLine("\u251C\u2500 ")
				}
				target.MoveBy(3, 0)
				x, y := target.Cpos()
				if highlighted || i != highlight[0] {
					drawPart(cpart, target, []int{-1})
				} else {
					drawPart(cpart, target, highlight[1:])
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
}
