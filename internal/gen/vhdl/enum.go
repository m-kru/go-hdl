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
	name := e.toEnumName()
	b.WriteString(
		fmt.Sprintf(
			"   function %s(slv : std_logic_vector(%d downto 0)) return %s;\n",
			name, e.Width()-1, e.name,
		),
	)
}

func (e enum) genToSlvDeclaration(b *strings.Builder) {
	name := e.paramName()
	b.WriteString(
		fmt.Sprintf(
			"   function to_slv(%s : %s) return std_logic_vector;\n",
			name, e.name,
		),
	)
}

func (e enum) genToStrDeclaration(b *strings.Builder) {
	name := e.paramName()
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
	b.WriteRune('\n')
	e.genToSlvDefinition(&b)
	b.WriteRune('\n')
	e.genToStrDefinition(&b)

	return b.String()
}

func (e enum) genToEnumDefinition(b *strings.Builder) {
	name := e.toEnumName()
	b.WriteString(
		fmt.Sprintf(
			"   function %s(slv : std_logic_vector(%d downto 0)) return %s is\n",
			name, e.Width()-1, e.name,
		),
	)
	b.WriteString("   begin\n")
	b.WriteString("      case slv is\n")
	for i, v := range e.values {
		b.WriteString(
			fmt.Sprintf(
				"         when \"%0*b\" => return %s;\n", e.Width(), i, v,
			),
		)
	}
	b.WriteString("         when others => report \"invalid slv value \" & to_string(slv) severity failure;\n")
	b.WriteString("      end case;\n")
	b.WriteString("   end function;\n")
}

func (e enum) genToSlvDefinition(b *strings.Builder) {
	name := e.paramName()
	b.WriteString(
		fmt.Sprintf(
			"   function to_slv(%s : %s) return std_logic_vector is\n",
			name, e.name,
		),
	)
	b.WriteString("   begin\n")
	b.WriteString(fmt.Sprintf("      case %s is\n", name))
	for i, v := range e.values {
		b.WriteString(
			fmt.Sprintf(
				"         when %s => return \"%0*b\";\n", v, e.Width(), i,
			),
		)
	}
	b.WriteString("      end case;\n")
	b.WriteString("   end function;\n")
}

func (e enum) genToStrDefinition(b *strings.Builder) {
	name := e.paramName()
	b.WriteString(
		fmt.Sprintf(
			"   function to_str(%s : %s) return string is\n",
			name, e.name,
		),
	)
	b.WriteString("   begin\n")
	b.WriteString(fmt.Sprintf("      case %s is\n", name))
	for _, v := range e.values {
		b.WriteString(
			fmt.Sprintf(
				"         when %[1]s => return \"%[1]s\";\n", v,
			),
		)
	}
	b.WriteString("      end case;\n")
	b.WriteString("   end function;\n")
}

func (e enum) toEnumName() string {
	name := e.name
	if strings.HasPrefix(name, "t_") {
		name = name[2:]
	}
	return "to_" + name
}

// paramName returns the name of the parameter that should be used when enum
// type is passed to the function.
func (e enum) paramName() string {
	name := e.name
	if strings.HasPrefix(name, "t_") {
		name = name[2:]
	} else {
		name = name[0:1]
	}
	return name
}
