package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Function struct {
	Symbol
}

func (f Function) SymbolNames() []string {
	return []string{}
}

func (f Function) GetSymbol(name string) []symbol.Symbol {
	return nil
}
