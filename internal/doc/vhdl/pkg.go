package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Package struct {
	Symbol
	symbols map[string]symbol.Symbol
}

func (p Package) AddSymbol(s symbol.Symbol) error {
	if _, ok := p.symbols[s.Name()]; ok {
		return fmt.Errorf(
			"symbol '%s' defined at least twice in package '%s'",
			s.Name(), p.Name(),
		)
	}

	p.symbols[s.Name()] = s

	return nil
}

func (p Package) SymbolNames() []string {
	names := []string{}

	for name, _ := range p.symbols {
		names = append(names, name)
	}

	return names
}

func (p Package) GetSymbol(name string) (symbol.Symbol, bool) {
	if sym, ok := p.symbols[name]; ok {
		return sym, true
	}
	return nil, false
}
