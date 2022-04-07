package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Entity struct {
	symbol
}

func (e Entity) SymbolNames() []string {
	return []string{}
}

func (e Entity) GetSymbol(name string) []sym.Symbol {
	return nil
}
