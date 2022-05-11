package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/sym"
	"sort"
	"sync"
)

type libraryContainer map[string]*lib.Library

var libContainer libraryContainer = libraryContainer{}

var libContainerMutex sync.Mutex

// Add adds library in atomic way. If library with given name
// already exists, then it is not overwritten.
func (lc libraryContainer) Add(l *lib.Library) {
	libContainerMutex.Lock()
	if _, ok := lc[l.Key()]; !ok {
		lc[l.Key()] = l
	}
	libContainerMutex.Unlock()
}

func (lc libraryContainer) Get(name string) *lib.Library {
	if _, ok := lc[name]; !ok {
		panic("should never happen")
	}
	return lc[name]
}

func (lc libraryContainer) AddSymbol(libName string, s sym.Symbol) {
	libContainerMutex.Lock()
	lc[libName].AddSymbol(s)
	libContainerMutex.Unlock()
}

// LibraryNames returns library names sorted in alphabetical order.
func LibraryNames() []string {
	names := []string{}

	for name, _ := range libContainer {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func GetLibrary(name string) (*lib.Library, bool) {
	if l, ok := libContainer[name]; ok {
		return l, true
	}
	return nil, false
}
