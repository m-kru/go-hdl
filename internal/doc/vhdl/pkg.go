package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Package struct {
	Symbol
	Consts map[symbol.ID]symbol.Symbol
	Funcs  map[symbol.ID]symbol.Symbol
	Procs  map[symbol.ID]symbol.Symbol
	Types  map[symbol.ID]symbol.Symbol
}

func (p Package) AddSymbol(s symbol.Symbol) error {
	id := symbol.ID{Name: s.Name(), LineNum: s.LineNum()}

	switch s.(type) {
	case Constant:
		if _, ok := p.Consts[id]; ok {
			return fmt.Errorf(
				"constant '%s' defined at least twice in package '%s'",
				s.Name(), p.Name(),
			)
		}
		p.Consts[id] = s
	case Function:
		p.Funcs[id] = s
	case Type:
		if _, ok := p.Types[id]; ok {
			return fmt.Errorf(
				"type '%s' defined at least twice in package '%s'",
				s.Name(), p.Name(),
			)
		}
		p.Types[id] = s
	default:
		panic("should never happen")
	}

	return nil
}

func (p Package) SymbolNames() []string {
	names := []string{}

	for id, _ := range p.Consts {
		names = append(names, id.Name)
	}
	for id, _ := range p.Funcs {
		names = append(names, id.Name)
	}
	for id, _ := range p.Procs {
		names = append(names, id.Name)
	}
	for id, _ := range p.Types {
		names = append(names, id.Name)
	}

	return names
}

func (p Package) GetSymbol(name string) []symbol.Symbol {
	syms := []symbol.Symbol{}

	for id, s := range p.Consts {
		if id.Name == name {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Funcs {
		if id.Name == name {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Types {
		if id.Name == name {
			syms = append(syms, s)
		}
	}

	return syms
}
