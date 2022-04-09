package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
)

type Entity struct {
	symbol
}

func (e Entity) InnerKeys() []string               { return []string{} }
func (e Entity) GetSymbol(key string) []sym.Symbol { return nil }
