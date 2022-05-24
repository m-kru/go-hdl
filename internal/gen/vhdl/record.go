package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/gen/gen"
	"strings"
)

type field struct {
	name         string
	typ          string
	width        int
	toRecordFunc string
	toSlvFunc    string
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

func (r *record) GenDeclarations(gens map[string]gen.Generable) string {
	b := strings.Builder{}

	r.genToRecordDeclaration(&b)
	r.genToSlvDeclaration(&b)
	if !r.noToStr {
		r.genToStrDeclaration(&b)
	}

	return b.String()
}

func (r *record) genToRecordDeclaration(b *strings.Builder) {
	name := toTypeFuncName(r.name)
	b.WriteString(
		spf(
			"   function %s(slv : std_logic_vector(%d downto 0)) return %s;\n",
			name, r.Width()-1, r.name,
		),
	)
}

func (r *record) genToSlvDeclaration(b *strings.Builder) {
	name := funcParamName(r.name)
	b.WriteString(
		spf(
			"   function to_slv(%s : %s) return std_logic_vector;\n",
			name, r.name,
		),
	)
}

func (r *record) genToStrDeclaration(b *strings.Builder) {
	name := funcParamName(r.name)
	b.WriteString(
		spf(
			"   function to_str(%s : %s) return string;\n",
			name, r.name,
		),
	)
}

func (r *record) GenDefinitions(gens map[string]gen.Generable) string {
	b := strings.Builder{}

	r.genToRecordDefinition(&b)
	/*
			e.genToSlvDefinition(&b)
		if !r.noToStr {
			e.genToStrDefinition(&b)
		}
	*/

	return b.String()
}

func (r *record) genToRecordDefinition(b *strings.Builder) {
	name := toTypeFuncName(r.name)
	width := r.Width() - 1

	bws := b.WriteString

	bws(
		spf(
			"   function %s(slv : std_logic_vector(%d downto 0)) return %s is\n",
			name, width, r.name,
		),
	)
	varName := funcParamName(r.name)
	bws(spf("      variable %s : %s;\n", varName, r.name))
	bws("   begin\n")

	for i, _ := range r.fields {
		width = r.slvToField(i, b, width)
	}

	bws(spf("      return %s;\n", varName))
	bws("   end function;\n")
}

func (r *record) slvToField(idx int, b *strings.Builder, width int) int {
	bws := b.WriteString
	varName := funcParamName(r.name)

	f := r.fields[idx]
	typ := f.typ

	switch typ {
	case "std_logic", "std_ulogic":
		bws(spf("      %s.%s := slv(%d);\n", varName, f.name, width))
	case "bit", "boolean":
		one := "'1'"
		zero := "'0'"
		if typ == "boolean" {
			one = "true"
			zero = "false"
		}
		bws(spf("      if slv(%d) = '1' then\n", width))
		bws(spf("         %s.%s := %s;\n", varName, f.name, one))
		bws(spf("      elsif slv(%d) = '0' then\n", width))
		bws(spf("         %s.%s := %s;\n", varName, f.name, zero))
		bws("      else\n")
		bws(
			spf(
				"         report \"bit %[1]d: cannot convert \" & to_string(slv(%[1]d)) & \" to %[2]s type\" severity failure;\n", width, typ,
			),
		)
		bws("      end if;\n")
	case "integer":
		bws(spf("      %s.%s := to_integer(signed(slv(%d downto %d)));\n", varName, f.name, width, width-f.width+1))
	case "natural", "positive":
		bws(spf("      %s.%s := to_integer(unsigned(slv(%d downto %d)));\n", varName, f.name, width, width-f.width+1))
	case "std_logic_vector", "std_ulogic_vector":
		bws(spf("      %s.%s := slv(%d downto %d);\n", varName, f.name, width, width-f.width+1))
	case "signed", "unsigned":
		bws(spf("      %s.%s := %s(slv(%d downto %d));\n", varName, f.name, typ, width, width-f.width+1))
	default:
		panic("not yet implemented")
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
