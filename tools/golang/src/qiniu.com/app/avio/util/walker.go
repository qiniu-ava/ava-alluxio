package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type WalkStatus int

var (
	NotWalked WalkStatus = 0
	Walking   WalkStatus = 1
	Walked    WalkStatus = 2
)

func (w *WalkStatus) increase() *AvioError {
	switch *w {
	case NotWalked:
		w = &Walking
		break
	case Walking:
		w = &Walked
		break
	default:
		return &AvioError{
			code: WALK_ERROR_INVALID_STATUS,
			msg:  ErrorMessage[WALK_ERROR_INVALID_STATUS],
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
		return &AvioError{
			code: WALK_ERROR_INVALID_STATUS,
			msg:  ErrorMessage[WALK_ERROR_INVALID_STATUS],
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

var wg = sync.WaitGroup{}

type IWalker interface {
	Walk() *tree
}

// DFWalker depth-first walker
type DFWalker struct {
	root  string
	depth uint
	full  bool
	tree  *tree
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
	}
}

type WalkProgressHandler func(msg *WalkerMsg)

func (b *DFWalker) Walk(handler WalkProgressHandler) error {
	lastLen := 0
	for len(b.tree.children)%1000 == 0 && len(b.tree.children) == lastLen {
		lastLen = len(b.tree.children)
		e := b.tree.walk()
		if e != nil {
			return e
		}
		if b.full {
			break
		}
	}

	for _, childTree := range b.tree.children {
		handler(&WalkerMsg{
			Err:  nil,
			EOF:  false,
			Path: path.Join(b.tree.directory, b.tree.info.Name(), childTree.info.Name()),
			Info: childTree.info,
		})
		if b.depth > 1 && childTree.info.IsDir() {
			innerWalker := NewDFWalker(path.Join(b.tree.directory, b.tree.info.Name(), childTree.info.Name()), b.depth-1, b.full)
			innerWalker.Walk(func(msg *WalkerMsg) {
				if !msg.EOF {
					handler(msg)
				}
			})
		}
	}

	handler(&WalkerMsg{
		Err:  nil,
		EOF:  true,
		Path: "",
		Info: nil,
	})

	return nil
}

type WalkerMsg struct {
	Err  error
	EOF  bool
	Path string
	Info os.FileInfo
}
