package vhdl

import (
	"github.com/m-kru/go-hdl/internal/doc/sym"
)

type SubprogramsContainer interface {
	GetFunc(key string) []sym.Symbol
	GetProc(key string) []sym.Symbol
	SortedFuncKeys() []string // SortedFuncKeys must return function keys in alphabetical order.
	SortedProcKeys() []string // SortedProcKeys must return procedures keys in alphabetical order.
}
