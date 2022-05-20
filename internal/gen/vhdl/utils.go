package vhdl

import (
	"strings"
)

// toTypeFuncName returns name for the function converting std_logic_vector to particular type.
func toTypeFuncName(typeName string) string {
	name := typeName
	if strings.HasPrefix(name, "t_") {
		name = name[2:]
	}
	return "to_" + name
}

// funcParamName returns the name of the parameter that should be used when particular
// type is passed to the function.
func funcParamName(typeName string) string {
	name := typeName
	if strings.HasPrefix(name, "t_") {
		name = name[2:]
	} else {
		name = name[0:1]
	}
	return name
}
