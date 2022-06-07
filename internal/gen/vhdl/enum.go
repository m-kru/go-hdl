package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/enc"
	"github.com/m-kru/go-thdl/internal/gen/gen"
	"math"
	"strings"
)

type enum struct {
	name     string
	values   []string
	encoding string
}

func (e *enum) Name() string { return e.name }

func (e *enum) Width() int {
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

func (e *enum) GenDeclarations() string {
	b := strings.Builder{}

	e.genToEnumDeclaration(&b)
	e.genToSlvDeclaration(&b)
	e.genToStrDeclaration(&b)

	return b.String()
}

func (e *enum) genToEnumDeclaration(b *strings.Builder) {
	name := toTypeFuncName(e.name)
	b.WriteString(
		fmt.Sprintf(
			"   function %s(slv : std_logic_vector(%d downto 0)) return %s;\n",
			name, e.Width()-1, e.name,
		),
	)
}

func (e *enum) genToSlvDeclaration(b *strings.Builder) {
	name := funcParamName(e.name)
	b.WriteString(
		fmt.Sprintf(
			"   function to_slv(%s : %s) return std_logic_vector;\n",
			name, e.name,
		),
	)
}

func (e *enum) genToStrDeclaration(b *strings.Builder) {
	name := funcParamName(e.name)
	b.WriteString(
		fmt.Sprintf(
			"   function to_str(%s : %s) return string;\n",
			name, e.name,
		),
	)
}

func (e *enum) GenDefinitions(gens gen.Container) string {
	b := strings.Builder{}

	e.genToEnumDefinition(&b)
	b.WriteRune('\n')
	e.genToSlvDefinition(&b)
	b.WriteRune('\n')
	e.genToStrDefinition(&b)

	return b.String()
}

func (e *enum) genToEnumDefinition(b *strings.Builder) {
	name := toTypeFuncName(e.name)
	b.WriteString(
		fmt.Sprintf(
			"   function %s(slv : std_logic_vector(%d downto 0)) return %s is\n"+
				"   begin\n"+
				"      case slv is\n",
			name, e.Width()-1, e.name,
		),
	)

	for i, v := range e.values {
		b.WriteString(fmt.Sprintf("         when %s => return %s;\n", e.slv(i), v))
	}

	b.WriteString(
		"         when others => report \"invalid slv value \" & to_string(slv) severity failure;\n" +
			"      end case;\n" +
			"   end function;\n",
	)
}

func (e *enum) genToSlvDefinition(b *strings.Builder) {
	paramName := funcParamName(e.name)

	b.WriteString(
		fmt.Sprintf(
			"   function to_slv(%[1]s : %[2]s) return std_logic_vector is\n"+
				"   begin\n"+
				"      case %[1]s is\n",
			paramName, e.name,
		),
	)
	for i, v := range e.values {
		b.WriteString(fmt.Sprintf("         when %s => return %s;\n", v, e.slv(i)))
	}
	b.WriteString("      end case;\n   end function;\n")
}

func (e *enum) genToStrDefinition(b *strings.Builder) {
	paramName := funcParamName(e.name)

	b.WriteString(
		fmt.Sprintf(
			"   function to_str(%[1]s : %[2]s) return string is\n"+
				"   begin\n"+
				"      case %[1]s is\n",
			paramName, e.name,
		),
	)
	for _, v := range e.values {
		b.WriteString(fmt.Sprintf("         when %[1]s => return \"%[1]s\";\n", v))
	}
	b.WriteString("      end case;\n   end function;\n")
}

func (e *enum) ParseArgs(args []string) error {
	validParams := map[string]bool{
		"encoding": true,
	}
	validEncodings := map[string]bool{
		"gray": true, "one-hot": true, "sequential": true,
	}

	encoding := ""

	for _, arg := range args {
		splits := strings.Split(arg, "=")
		param := splits[0]

		if _, ok := validParams[param]; !ok {
			return fmt.Errorf("invalid parameter '%s'", param)
		}

		if len(splits) == 1 {
			return fmt.Errorf("missing argument for '%s' parameter", param)
		}
		a := splits[1]

		switch param {
		case "encoding":
			if _, ok := validEncodings[a]; !ok {
				return fmt.Errorf("invalid argument '%s' for 'encoding' parameter", a)
			}
			encoding = a
		}
	}

	if encoding == "" {
		encoding = "sequential"
	}
	e.encoding = encoding

	return nil
}

// slv returns std_logic_vector value for enum value of given index.
func (e *enum) slv(idx int) string {
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
