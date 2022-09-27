package vhdl

import (
	"fmt"
	"github.com/m-kru/go-hdl/internal/doc/sym"
	"github.com/m-kru/go-hdl/internal/utils"
)

type Function struct {
	symbol
	impure bool
}

func (f Function) InnerKeys() []string               { return []string{} }
func (f Function) GetSymbol(key string) []sym.Symbol { return nil }

func FuncsCodeSummary(funcs []sym.Symbol) string {
	var s string
	if len(funcs) == 1 {
		code := utils.Dewhitespace(funcs[0].Code())
		if utils.IsSingleLine(code) {
			s = code
		} else {
			s = fmt.Sprintf("function %s ...\n", funcs[0].Name())
		}
	} else {
		impureCount := 0
		for _, f := range funcs {
			if f.(Function).impure {
				impureCount += 1
			}
		}
		impurePrefix := ""
		if impureCount == len(funcs) {
			impurePrefix = "impure "
		} else if impureCount > 0 {
			impurePrefix = "(impure)? "
		}
		s = fmt.Sprintf("%sfunction %s ... (%d)\n", impurePrefix, funcs[0].Name(), len(funcs))
	}
	return s
}
