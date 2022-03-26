package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"sync"
)

type libraryContainer map[string]*lib.Library

var libContainer libraryContainer = libraryContainer{}

var libContainerMutex sync.Mutex

// Add adds library in atomic way. If library with given name
// already exists, then it is not overwritten.
func (lc libraryContainer) Add(l *lib.Library) {
	libContainerMutex.Lock()
	if _, ok := lc[l.Name()]; !ok {
		lc[l.Name()] = l
	}
	libContainerMutex.Unlock()
}

func (lc libraryContainer) AddSymbol(libName string, s symbol.Symbol) {
	libContainerMutex.Lock()
	lc[libName].AddSymbol(s)
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
