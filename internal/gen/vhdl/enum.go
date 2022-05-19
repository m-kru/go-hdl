package vhdl

import (
	"fmt"
	"math"
	"strings"
)

type enum struct {
	name   string
	values []string
}

func (e enum) Name() string { return e.name }

func (e enum) Width() uint {
	return uint(math.Ceil(math.Log2(float64(len(e.values)))))
}

func (e enum) GenDeclaration(args []string) string {
	b := strings.Builder{}

	e.genToEnumDeclaration(&b)
	e.genToSlvDeclaration(&b)
	e.genToStrDeclaration(&b)

	return b.String()
}

func (e enum) genToEnumDeclaration(b *strings.Builder) {
	name := e.name
	if strings.HasPrefix(name, "t_") {
		name = name[2:]
	}
	b.WriteString(
		fmt.Sprintf(
			"   function to_%s(slv : std_logic_vector(%d downto 0)) return %s;\n",
			name, e.Width()-1, e.name,
		),
	)
}

func (e enum) genToSlvDeclaration(b *strings.Builder) {
	name := e.name
	if strings.HasPrefix(name, "t_") {
		name = name[2:]
	} else {
		name = name[0:1]
	}
	b.WriteString(
		fmt.Sprintf(
			"   function to_slv(%s : %s) return std_logic_vector;\n",
			name, e.name,
		),
	)
}

func (e enum) genToStrDeclaration(b *strings.Builder) {
	name := e.name
	if strings.HasPrefix(name, "t_") {
		name = name[2:]
	} else {
		name = name[0:1]
	}
	b.WriteString(
		fmt.Sprintf(
			"   function to_str(%s : %s) return string;\n",
			name, e.name,
		),
	)
}

func (e enum) GenDefinition(args []string) string {
	b := strings.Builder{}

	e.genToEnumDefinition(&b)
	e.genToSlvDefinition(&b)
	e.genToStrDefinition(&b)

	return b.String()
}

func (e enum) genToEnumDefinition(b *strings.Builder) {
	b.WriteString("genToEnumDefinition\n")
}

func (e enum) genToSlvDefinition(b *strings.Builder) {
	b.WriteString("genToSlvDefinition\n")
}

func (e enum) genToStrDefinition(b *strings.Builder) {
	b.WriteString("genToStrDefinition")
}
