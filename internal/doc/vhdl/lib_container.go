package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"sync"
)

type libraryContainer map[string]*lib.Library

var libContainer libraryContainer = libraryContainer{}

var libContainerMutex sync.Mutex

func (lc libraryContainer) Has(name string) bool {
	_, ok := lc[name]
	return ok
}

func (lc libraryContainer) Add(l *lib.Library) {
	libContainerMutex.Lock()

	lc[l.Name()] = l

	libContainerMutex.Unlock()
}

func LibraryNames() []string {
	names := []string{}

	for name, _ := range libContainer {
		names = append(names, name)
	}

	return names
}

func GetLibrary(name string) (*lib.Library, bool) {
	if l, ok := libContainer[name]; ok {
		return l, true
	}
	return nil, false
}
