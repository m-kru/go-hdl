package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Procedure struct {
	symbol
}

func (p Procedure) SymbolNames() []string {
	return []string{}
}

func (p Procedure) GetSymbol(name string) []sym.Symbol {
	return nil
}
