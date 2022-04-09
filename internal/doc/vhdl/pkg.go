package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Package struct {
	symbol
	Consts map[sym.ID]sym.Symbol
	Funcs  map[sym.ID]sym.Symbol
	Procs  map[sym.ID]sym.Symbol
	Types  map[sym.ID]sym.Symbol
}

func (p Package) AddSymbol(s sym.Symbol) error {
	id := sym.ID{Key: s.Key(), LineNum: s.LineNum()}

	switch s.(type) {
	case Constant:
		if _, ok := p.Consts[id]; ok {
			return fmt.Errorf(
				"constant '%s' defined at least twice in package '%s'",
				s.Key(), p.Key(),
			)
		}
		p.Consts[id] = s
	case Function:
		p.Funcs[id] = s
	case Procedure:
		p.Procs[id] = s
	case Type:
		if _, ok := p.Types[id]; ok {
			return fmt.Errorf(
				"type '%s' defined at least twice in package '%s'",
				s.Key(), p.Key(),
			)
		}
		p.Types[id] = s
	default:
		panic("should never happen")
	}

	return nil
}

func (p Package) InnerKeys() []string {
	names := []string{}

	for id, _ := range p.Consts {
		names = append(names, id.Key)
	}
	for id, _ := range p.Funcs {
		names = append(names, id.Key)
	}
	for id, _ := range p.Procs {
		names = append(names, id.Key)
	}
	for id, _ := range p.Types {
		names = append(names, id.Key)
	}

	return names
}

func (p Package) GetSymbol(key string) []sym.Symbol {
	syms := []sym.Symbol{}

	for id, s := range p.Consts {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Funcs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Procs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Types {
		if id.Key == key {
			syms = append(syms, s)
		}
	}

	return syms
}
