package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Protected struct {
	symbol
}

func (p Protected) InnerKeys() []string               { return []string{} }
func (p Protected) GetSymbol(key string) []sym.Symbol { return nil }
