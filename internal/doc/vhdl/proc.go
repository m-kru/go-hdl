package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Procedure struct {
	Symbol
}

func (p Procedure) SymbolNames() []string {
	return []string{}
}

func (p Procedure) GetSymbol(name string) []symbol.Symbol {
	return nil
}
