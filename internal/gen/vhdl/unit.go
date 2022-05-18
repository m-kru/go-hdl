package vhdl

import (
	"github.com/m-kru/go-thdl/internal/gen/gen"
)

// Unit represents VHDL design unit.
// It is needed as single file can contain multiple design units.
//
// lineNum is useful when file needs to be regenerated, as there is no need to match
// against regex again and do the name comparison.
type unit struct {
	name    string
	lineNum uint
	typ     string
	gens    map[string]gen.Generable
}
