package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Function struct {
	symbol
	impure bool
}

func (f Function) InnerKeys() []string               { return []string{} }
func (f Function) GetSymbol(key string) []sym.Symbol { return nil }
