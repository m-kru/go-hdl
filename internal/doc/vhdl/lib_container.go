package vhdl

import (
	"sync"
)

type libraryContainer map[string]*Library

var libContainer libraryContainer = libraryContainer{
	"_unknown_": &Library{name: "_unknown_"},
}
var libContainerMutex sync.Mutex

func (lc libraryContainer) Has(name string) bool {
	_, ok := lc[name]
	return ok
}

func (lc libraryContainer) Add(lib Library) {
	libContainerMutex.Lock()

	lc[lib.name] = &lib

	libContainerMutex.Unlock()
}

func LibraryNames() []string {
	names := []string{}

	for name, _ := range libContainer {
		names = append(names, name)
	}

	return names
}

func GetLibrary(name string) (*Library, bool) {
	if l, ok := libContainer[name]; ok {
		return l, true
	}
	return nil, false
}
