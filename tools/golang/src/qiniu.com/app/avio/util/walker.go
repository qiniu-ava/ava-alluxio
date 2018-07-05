package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type WalkStatus int

var (
	NotWalked WalkStatus = 0
	Walking   WalkStatus = 1
	Walked    WalkStatus = 2
)

func (w *WalkStatus) increase() *WalkError {
	switch *w {
	case NotWalked:
		w = &Walking
		break
	case Walking:
		w = &Walked
		break
	default:
		return &WalkError{
			code: WalkStatusError,
			msg:  "invalid walk status to increase",
		}
	}
	return nil
}

type tree struct {
	info      os.FileInfo
	directory string
	status    WalkStatus
	// depth    uint
	// maxDepth uint
	children map[string]tree
}

func (t *tree) walk() error {
	if t.status != NotWalked {
		return &WalkError{
			code: WalkStatusError,
			msg:  "invalid walk status to start walk",
		}
	}

	if e := t.status.increase(); e != nil {
		return e
	}

	files, e := ioutil.ReadDir(path.Join(t.directory, t.info.Name()))
	if e != nil {
		return e
	}

	for _, fileInfo := range files {
		t.children[fileInfo.Name()] = tree{
			info:      fileInfo,
			directory: path.Join(t.directory, t.info.Name()),
			status:    NotWalked,
			children:  make(map[string]tree),
		}
	}

	return nil
}

type IWalker interface {
	Walk() *tree
}

// DFWalker depth-first walker
type DFWalker struct {
	root        string
	depth       uint
	full        bool
	tree        *tree
	currentPath []string
	generatorCh chan GeneratorMsg
}

func NewDFWalker(root string, depth uint, full bool) *DFWalker {
	fileInfo, e := os.Stat(root)
	if e != nil {
		panic(fmt.Sprintf("no directory or file named: %s", root))
	}

	director := path.Dir(root)

	return &DFWalker{
		root:  root,
		depth: depth,
		full:  full,
		tree: &tree{
			info:      fileInfo,
			directory: director,
			status:    NotWalked,
			children:  make(map[string]tree),
		},
		generatorCh: make(chan GeneratorMsg, 1024),
	}
}

func (b *DFWalker) Walk() (g PathGenerator, e error) {
	lastLen := 0
	for len(b.tree.children)%1000 == 0 && len(b.tree.children) == lastLen {
		lastLen = len(b.tree.children)
		e = b.tree.walk()
		if e != nil {
			return
		}
		if b.full {
			break
		}
	}

	g = PathGenerator{b.generatorCh}
	for _, childTree := range b.tree.children {
		b.postMsg(nil, false, path.Join(b.tree.directory, b.tree.info.Name(), childTree.info.Name()), childTree.info)
	}

	b.postMsg(nil, true, "", nil)

	return
}

func (b *DFWalker) postMsg(e error, eof bool, path string, info os.FileInfo) {
	b.generatorCh <- GeneratorMsg{
		Err:  e,
		EOF:  eof,
		Path: path,
		Info: info,
	}
}

func (b *DFWalker) Close() {
	close(b.generatorCh)
}

type GeneratorMsg struct {
	Err  error
	EOF  bool
	Path string
	Info os.FileInfo
}

type PathGenerator struct {
	ch chan GeneratorMsg
}

func (g *PathGenerator) Next() (msg GeneratorMsg) {
	msg = <-g.ch
	return
}

// @{TODO} breadth-first walker
type BFWalker struct {
}

func (b *BFWalker) Walk() (g PathGenerator) {
	return PathGenerator{}
}
