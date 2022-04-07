package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type PackageInstantiation struct {
	symbol
}

func (pi PackageInstantiation) SymbolNames() []string {
	return []string{}
}

func (pi PackageInstantiation) GetSymbol(name string) []sym.Symbol {
	return nil
}
