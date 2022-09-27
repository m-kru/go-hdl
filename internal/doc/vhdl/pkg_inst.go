package vhdl

import (
	"github.com/m-kru/go-hdl/internal/doc/sym"
)

type PackageInstantiation struct {
	symbol
}

func (pi PackageInstantiation) InnerKeys() []string               { return []string{} }
func (pi PackageInstantiation) GetSymbol(key string) []sym.Symbol { return nil }
