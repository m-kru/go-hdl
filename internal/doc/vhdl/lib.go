package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"sync"
)

type Library struct {
	name string

	files      []string
	filesMutex sync.Mutex

	symbols      map[string]symbol.Symbol
	symbolsMutex sync.Mutex
}

func MakeLibrary(name string) Library {
	return Library{
		files:   []string{},
		name:    name,
		symbols: map[string]symbol.Symbol{},
	}
}

func (l *Library) AddFile(f string) {
	l.filesMutex.Lock()
	l.files = append(l.files, f)
	l.filesMutex.Unlock()
}

func (l *Library) Name() string { return l.name }

func (l *Library) Doc() string {
	return "VHDL Library Doc"
}

func (l *Library) Code() string {
	return "VHDL Library Code"
}

func (l *Library) SymbolNames() []string {
	names := []string{}

	for name, _ := range l.symbols {
		names = append(names, name)
	}

	return names
}

func (l *Library) GetSymbol(name string) (symbol.Symbol, bool) {
	if s, ok := l.symbols[name]; ok {
		return s, true
	}
	return nil, false
}

func (l *Library) AddSymbol(s symbol.Symbol) {
	l.symbolsMutex.Lock()
	l.symbols[s.Name()] = s
	l.symbolsMutex.Unlock()
}
