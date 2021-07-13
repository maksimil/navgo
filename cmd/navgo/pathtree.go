package main

import (
	"os"
	"path"
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
