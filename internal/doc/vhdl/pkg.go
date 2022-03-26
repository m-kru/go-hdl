package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Package struct {
	Symbol
	Consts map[string]symbol.Symbol
	Funcs  map[string]symbol.Symbol
	Procs  map[string]symbol.Symbol
	Types  map[string]symbol.Symbol
}

func (p Package) AddSymbol(s symbol.Symbol) error {
	switch s.(type) {
	case Type:
		if _, ok := p.Consts[s.Name()]; ok {
			return fmt.Errorf(
				"type '%s' defined at least twice in package '%s'",
				s.Name(), p.Name(),
			)
		}
		p.Consts[s.Name()] = s
	default:
		panic("should never happen")
	}

	return nil
}

func (p Package) SymbolNames() []string {
	names := []string{}

	for name, _ := range p.Consts {
		names = append(names, name)
	}
	for name, _ := range p.Funcs {
		names = append(names, name)
	}
	for name, _ := range p.Procs {
		names = append(names, name)
	}
	for name, _ := range p.Types {
		names = append(names, name)
	}

	return names
}

func (p Package) GetSymbol(name string) (symbol.Symbol, bool) {
	if sym, ok := p.Consts[name]; ok {
		return sym, true
	}
	if sym, ok := p.Types[name]; ok {
		return sym, true
	}
	// TODO: Thinks how to handle functions and procedures.
	// Their names can be overloaded.
	return nil, false
}
