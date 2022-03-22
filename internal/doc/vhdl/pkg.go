package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Package struct {
	Symbol
}

func (p Package) SymbolNames() []string {
	return []string{}
}

func (p Package) GetSymbol(name string) (symbol.Symbol, bool) {
	return nil, false
}
