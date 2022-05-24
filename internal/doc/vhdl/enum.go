package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Enum struct {
	symbol
}

func (e Enum) InnerKeys() []string               { return []string{} }
func (e Enum) GetSymbol(key string) []sym.Symbol { return nil }

func (e Enum) kind() string { return "enum" }
