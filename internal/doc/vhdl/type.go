package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Type struct {
	Symbol
}

func (t Type) SymbolNames() []string {
	return []string{}
}

func (t Type) GetSymbol(name string) []symbol.Symbol {
	return nil
}
