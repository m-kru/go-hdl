package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Alias struct {
	symbol
}

func (a Alias) InnerKeys() []string               { return []string{} }
func (a Alias) GetSymbol(key string) []sym.Symbol { return nil }
