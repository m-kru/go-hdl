package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Entity struct {
	Symbol
}

func (e Entity) SymbolNames() []string {
	return []string{}
}

func (e Entity) GetSymbol(name string) []symbol.Symbol {
	return nil
}
