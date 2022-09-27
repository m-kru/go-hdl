package vhdl

import (
	"github.com/m-kru/go-hdl/internal/doc/sym"
)

type Array struct {
	symbol
}

func (a Array) InnerKeys() []string               { return []string{} }
func (a Array) GetSymbol(key string) []sym.Symbol { return nil }

func (a Array) kind() string { return "array" }
