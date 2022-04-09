package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Procedure struct {
	symbol
}

func (p Procedure) InnerKeys() []string               { return []string{} }
func (p Procedure) GetSymbol(key string) []sym.Symbol { return nil }
