package vhdl

import (
	"github.com/m-kru/go-hdl/internal/doc/sym"
)

type Constant struct {
	symbol
}

func (c Constant) InnerKeys() []string               { return []string{} }
func (c Constant) GetSymbol(key string) []sym.Symbol { return nil }
