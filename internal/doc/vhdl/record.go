package vhdl

import (
	"github.com/m-kru/go-hdl/internal/doc/sym"
)

type Record struct {
	symbol
}

func (r Record) InnerKeys() []string               { return []string{} }
func (r Record) GetSymbol(key string) []sym.Symbol { return nil }

func (r Record) kind() string { return "record" }
