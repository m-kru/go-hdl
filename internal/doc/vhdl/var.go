package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Variable struct {
	symbol
}

func (v Variable) InnerKeys() []string               { return []string{} }
func (v Variable) GetSymbol(key string) []sym.Symbol { return nil }
