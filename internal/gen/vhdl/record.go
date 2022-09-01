package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/gen/gen"
	"strings"
)

type field struct {
	name   string
	typ    string
	width  int
	toType string
	toSlv  string
	toStr  string
}

type record struct {
	name    string
	fields  []field
	noToStr bool
}

func (r *record) Name() string { return r.name }

func (r *record) Width() int {
	width := 0
	for _, f := range r.fields {
		width += f.width
	}
	return width
}

func (r *record) GenDeclarations() string {
	b := strings.Builder{}

	r.genToRecordDeclaration(&b)
	r.genToSlvDeclaration(&b)
	if !r.noToStr {
		r.genToStrDeclaration(&b)
	}

	return b.String()
}

func (r *record) genToRecordDeclaration(b *strings.Builder) {
	funcName := toTypeFuncName(r.name)
	b.WriteString(
		fmt.Sprintf(
			"   function %s(slv : std_logic_vector(%d downto 0)) return %s;\n",
			funcName, r.Width()-1, r.name,
		),
	)
}

func (r *record) genToSlvDeclaration(b *strings.Builder) {
	paramName := funcParamName(r.name)
	b.WriteString(
		fmt.Sprintf(
			"   function to_slv(%s : %s) return std_logic_vector;\n",
			paramName, r.name,
		),
	)
}

func (r *record) genToStrDeclaration(b *strings.Builder) {
	paramName := funcParamName(r.name)
	b.WriteString(
		fmt.Sprintf(
			"   function to_str(%s : %s; add_names : boolean := false) return string;\n",
			paramName, r.name,
		),
	)
}

func (r *record) GenDefinitions(gens gen.Container) string {
	b := strings.Builder{}

	r.genToRecordDefinition(gens, &b)
	b.WriteRune('\n')
	r.genToSlvDefinition(gens, &b)
	if !r.noToStr {
		b.WriteRune('\n')
		r.genToStrDefinition(gens, &b)
	}

	return b.String()
}

func (r *record) genToRecordDefinition(gens gen.Container, b *strings.Builder) {
	funcName := toTypeFuncName(r.name)
	varName := funcParamName(r.name)
	width := r.Width() - 1

	b.WriteString(
		fmt.Sprintf(
			"   function %[1]s(slv : std_logic_vector(%[2]d downto 0)) return %[3]s is\n"+
				"      variable %[4]s : %[3]s;\n"+
				"   begin\n",
			funcName, width, r.name, varName,
		),
	)

	for i := range r.fields {
		width = r.slvToField(i, gens, b, width)
	}

	b.WriteString(
		fmt.Sprintf(
			"      return %s;\n"+
				"   end function;\n",
			varName,
		),
	)
}

func (r *record) genToSlvDefinition(gens gen.Container, b *strings.Builder) {
	paramName := funcParamName(r.name)
	width := r.Width() - 1

	b.WriteString(
		fmt.Sprintf(
			"   function to_slv(%s : %s) return std_logic_vector is\n"+
				"      variable slv : std_logic_vector(%d downto 0);\n"+
				"   begin\n",
			paramName, r.name, width,
		),
	)

	for i := range r.fields {
		width = r.fieldToSlv(i, gens, b, width)
	}

	b.WriteString("      return slv;\n   end function;\n")
}

func (r *record) slvToField(idx int, gens gen.Container, b *strings.Builder, width int) int {
	varName := funcParamName(r.name)

	f := r.fields[idx]
	typ := f.typ

	switch typ {
	case "std_logic", "std_ulogic":
		b.WriteString(
			fmt.Sprintf("      %s.%s := slv(%d);\n", varName, f.name, width),
		)
	case "bit", "boolean":
		one := "'1'"
		zero := "'0'"
		if typ == "boolean" {
			one = "true"
			zero = "false"
		}
		b.WriteString(
			fmt.Sprintf(
				"      if slv(%[1]d) = '1' then\n"+
					"         %[2]s.%[3]s := %[4]s;\n"+
					"      elsif slv(%[1]d) = '0' then\n"+
					"         %[2]s.%[3]s := %[5]s;\n"+
					"      else\n"+
					"         report \"bit %[1]d: cannot convert \" & to_string(slv(%[1]d)) & \" to %[6]s type\" severity failure;\n"+
					"      end if;\n",
				width, varName, f.name, one, zero, typ,
			),
		)
	case "integer":
		b.WriteString(
			fmt.Sprintf(
				"      %s.%s := to_integer(signed(slv(%d downto %d)));\n",
				varName, f.name, width, width-f.width+1,
			),
		)
	case "natural", "positive":
		b.WriteString(
			fmt.Sprintf(
				"      %s.%s := to_integer(unsigned(slv(%d downto %d)));\n",
				varName, f.name, width, width-f.width+1,
			),
		)
	case "std_logic_vector", "std_ulogic_vector":
		b.WriteString(
			fmt.Sprintf(
				"      %s.%s := slv(%d downto %d);\n",
				varName, f.name, width, width-f.width+1,
			),
		)
	case "signed", "unsigned":
		b.WriteString(
			fmt.Sprintf(
				"      %s.%s := %s(slv(%d downto %d));\n",
				varName, f.name, typ, width, width-f.width+1,
			),
		)
	default:
		if g, ok := gens.Get(typ); ok {
			b.WriteString(
				fmt.Sprintf(
					"      %s.%s := %s(slv(%d downto %d));\n",
					varName, f.name, toTypeFuncName(g.Name()), width, width-f.width+1,
				),
			)
		} else if f.width != 0 {
			funcName := toTypeFuncName(typ)
			if f.toType != "" {
				funcName = f.toType
			}
			b.WriteString(
				fmt.Sprintf(
					"      %s.%s := %s(slv(%d downto %d));\n",
					varName, f.name, funcName, width, width-f.width+1,
				),
			)
		} else {
			panic("should never happen")
		}
	}

	width -= f.width

	return width
}

func (r *record) fieldToSlv(idx int, gens gen.Container, b *strings.Builder, width int) int {
	varName := funcParamName(r.name)

	f := r.fields[idx]
	typ := f.typ

	switch typ {
	case "std_logic", "std_ulogic":
		b.WriteString(
			fmt.Sprintf("      slv(%d) := %s.%s;\n", width, varName, f.name),
		)
	case "bit":
		b.WriteString(
			fmt.Sprintf(
				"      if %[1]s.%[2]s = '1' then slv(%[3]d) := '1'; else slv(%[3]d) := '0'; end if;\n",
				varName, f.name, width,
			),
		)
	case "boolean":
		b.WriteString(
			fmt.Sprintf(
				"      if %[1]s.%[2]s then slv(%[3]d) := '1'; else slv(%[3]d) := '0'; end if;\n",
				varName, f.name, width,
			),
		)
	case "integer":
		b.WriteString(
			fmt.Sprintf(
				"      slv(%d downto %d) := std_logic_vector(to_signed(%s.%s, 32));\n",
				width, width-f.width+1, varName, f.name,
			),
		)
	case "natural", "positive":
		b.WriteString(
			fmt.Sprintf(
				"      slv(%d downto %d) := std_logic_vector(to_unsigned(%s.%s, 32));\n",
				width, width-f.width+1, varName, f.name,
			),
		)
	case "std_logic_vector", "std_ulogic_vector":
		b.WriteString(
			fmt.Sprintf(
				"      slv(%d downto %d) := %s.%s;\n",
				width, width-f.width+1, varName, f.name,
			),
		)
	case "signed", "unsigned":
		b.WriteString(
			fmt.Sprintf(
				"      slv(%d downto %d) := std_logic_vector(%s.%s);\n",
				width, width-f.width+1, varName, f.name,
			),
		)
	default:
		if _, ok := gens.Get(typ); ok {
			b.WriteString(
				fmt.Sprintf(
					"      slv(%d downto %d) := to_slv(%s.%s);\n",
					width, width-f.width+1, varName, f.name,
				),
			)
		} else if f.width != 0 {
			funcName := "to_slv"
			if f.toSlv != "" {
				funcName = f.toSlv
			}
			b.WriteString(
				fmt.Sprintf(
					"      slv(%d downto %d) := %s(%s.%s);\n",
					width, width-f.width+1, funcName, varName, f.name,
				),
			)
		} else {
			panic("should never happen")
		}
	}

	width -= f.width

	return width
}

func (r *record) ParseArgs(args []string) error {
	validFlags := map[string]bool{
		"no-to-str": true,
	}

	for _, arg := range args {
		if _, ok := validFlags[arg]; !ok {
			return fmt.Errorf("invalid flag '%s'", arg)
		}

		switch arg {
		case "no-to-str":
			r.noToStr = true
		}
	}

	return nil
}

func (r *record) genToStrDefinition(gens gen.Container, b *strings.Builder) {
	paramName := funcParamName(r.name)

	b.WriteString(
		fmt.Sprintf(
			"   function to_str(%s : %s; add_names : boolean := false) return string is\n"+
				"   begin\n"+
				"      if add_names then\n"+
				"         return \"(\" &",
			paramName, r.name,
		),
	)
	for i := range r.fields {
		r.fieldToStr(i, gens, true, b)
	}
	b.WriteString(" & \")\";\n      end if;\n      return \"(\" &")
	for i := range r.fields {
		r.fieldToStr(i, gens, false, b)
	}

	b.WriteString(" & \")\";\n   end function;\n")
}

func (r *record) fieldToStr(idx int, gens gen.Container, withName bool, b *strings.Builder) {
	paramName := funcParamName(r.name)

	f := r.fields[idx]
	typ := f.typ

	if idx != 0 {
		b.WriteString(" & \", \" &")
	}

	if withName {
		b.WriteString(fmt.Sprintf(" \"%s => \" &", f.name))
	}

	switch typ {
	case "bit", "boolean", "std_logic", "std_ulogic", "std_logic_vector", "integer", "natural", "positive", "signed", "unsigned":
		b.WriteString(fmt.Sprintf(" to_string(%s.%s)", paramName, f.name))
	default:
		if _, ok := gens.Get(typ); ok {
			b.WriteString(fmt.Sprintf(" to_str(%s.%s)", paramName, f.name))
		} else if f.width != 0 {
			toStr := f.toStr
			if toStr == "" {
				toStr = "to_str"
			}
			b.WriteString(fmt.Sprintf(" %s(%s.%s)", f.toStr, paramName, f.name))
		} else {
			panic("should never happen")
		}
	}
}
