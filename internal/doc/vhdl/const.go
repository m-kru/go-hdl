package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Constant struct {
	Symbol
}

func (c Constant) SymbolNames() []string {
	return []string{}
}

func (c Constant) GetSymbol(name string) (symbol.Symbol, bool) {
	return nil, false
}
