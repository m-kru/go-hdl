package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Type struct {
	symbol
	kind string
}

func (t Type) InnerKeys() []string               { return []string{} }
func (t Type) GetSymbol(key string) []sym.Symbol { return nil }
