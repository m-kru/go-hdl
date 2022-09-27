package vhdl

import (
	"fmt"
	"github.com/m-kru/go-hdl/internal/doc/sym"
	"github.com/m-kru/go-hdl/internal/utils"
)

type Procedure struct {
	symbol
}

func (p Procedure) InnerKeys() []string               { return []string{} }
func (p Procedure) GetSymbol(key string) []sym.Symbol { return nil }

func ProcsCodeSummary(procs []sym.Symbol) string {
	var s string
	if len(procs) == 1 {
		code := utils.Dewhitespace(procs[0].Code())
		if utils.IsSingleLine(code) {
			s = code
		} else {
			s = fmt.Sprintf("procedure %s ...\n", procs[0].Name())
		}
	} else {
		s = fmt.Sprintf("procedure %s ... (%d)\n", procs[0].Name(), len(procs))
	}
	return s
}
