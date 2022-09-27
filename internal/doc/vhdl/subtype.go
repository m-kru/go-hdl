package vhdl

import (
	"github.com/m-kru/go-hdl/internal/doc/sym"
)

type Subtype struct {
	symbol
}

func (s Subtype) InnerKeys() []string               { return []string{} }
func (s Subtype) GetSymbol(key string) []sym.Symbol { return nil }
