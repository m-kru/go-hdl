package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/enc"
	"math"
	"strings"
)

type enum struct {
	name     string
	values   []string
	encoding string
}

func (e enum) Name() string { return e.name }

func (e enum) Width() int {
	switch e.encoding {
	case "one-hot":
		return len(e.values)
	case "gray":
		panic("not yet implemented")
	case "sequential":
		return int(math.Ceil(math.Log2(float64(len(e.values)))))
	default:
		panic("should never happen")
	}
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
		b.WriteString(fmt.Sprintf("         when %s => return %s;\n", e.slv(i), v))
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
				"         when %s => return %s;\n", v, e.slv(i),
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

// toEnumName returns name for the function converting std_logic_vector to enum type.
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

func (e *enum) parseArgs(args []string) error {
	validFlags := map[string]bool{
		"-gray": true, "-one-hot": true,
	}

	encoding := ""

	for _, a := range args {
		if a[0] != '-' {
			return fmt.Errorf("invalid argument '%s'", a)
		} else {
			if _, ok := validFlags[a]; !ok {
				return fmt.Errorf("invalid argument '%s'", a)
			}
			switch a {
			case "-one-hot", "-gray":
				if encoding != "" {
					return fmt.Errorf(
						"cannot set '%s' encoding, as '%s' encoding is already set", a[1:], encoding,
					)
				}
				encoding = a[1:]
			}
		}
	}

	if encoding == "" {
		encoding = "sequential"
	}
	e.encoding = encoding

	return nil
}

// slv returns std_logic_vector value for enum value of given index.
func (e enum) slv(idx int) string {
	var s string
	switch e.encoding {
	case "one-hot":
		s = enc.OneHot(idx, e.Width())
	case "gray":
		panic("not yet implemented")
	case "sequential":
		s = fmt.Sprintf("%0*b", e.Width(), idx)
	default:
		panic("should never happen")
	}

	return "\"" + s + "\""
}
