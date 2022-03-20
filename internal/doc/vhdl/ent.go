package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/symbol"
)

type Entity struct {
	filepath *string
	name     string

	hasDoc   bool
	docStart uint32
	docEnd   uint32

	codeStart uint32
	codeEnd   uint32
}

func (e Entity) Name() string { return e.name }

func (e Entity) Doc() string {
	return "VHDL Entity Doc"
}

func (e Entity) Code() string {
	return "VHDL Entity Code"
}

func (e Entity) SymbolNames() []string {
	return []string{}
}

func (e Entity) GetSymbol(name string) (symbol.Symbol, bool) {
	panic("VHDL entity symbol isn't a symbol container")
}
