package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Constant struct {
	symbol
}

func (c Constant) SymbolNames() []string {
	return []string{}
}

func (c Constant) GetSymbol(name string) []sym.Symbol {
	return nil
}
