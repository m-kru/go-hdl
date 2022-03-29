package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type PackageInstantiation struct {
	Symbol
}

func (pi PackageInstantiation) SymbolNames() []string {
	return []string{}
}

func (pi PackageInstantiation) GetSymbol(name string) []symbol.Symbol {
	return nil
}
