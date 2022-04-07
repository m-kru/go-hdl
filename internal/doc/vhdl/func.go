package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Function struct {
	symbol
}

func (f Function) SymbolNames() []string {
	return []string{}
}

func (f Function) GetSymbol(name string) []sym.Symbol {
	return nil
}
