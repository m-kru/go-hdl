package vhdl

import (
	"math"
	"strings"
)

type Enum struct {
	name   string
	values []string
}

func (e Enum) Name() string { return e.name }

func (e Enum) Width() uint {
	return uint(math.Ceil(math.Log2(float64(len(e.values)))))
}

func (e Enum) GenDeclaration(args []string) string {
	b := strings.Builder{}

	e.genToEnumDeclaration(&b)
	e.genToSlvDeclaration(&b)
	e.genToStrDeclaration(&b)

	return b.String()
}

func (e Enum) genToEnumDeclaration(b *strings.Builder) {
	b.WriteString("genToEnumDeclaration\n")
}

func (e Enum) genToSlvDeclaration(b *strings.Builder) {
	b.WriteString("genToSlvDeclaration\n")
}

func (e Enum) genToStrDeclaration(b *strings.Builder) {
	b.WriteString("genToStrDeclaration\n")
}

func (e Enum) GenDefinition(args []string) string {
	b := strings.Builder{}

	e.genToEnumDefinition(&b)
	e.genToSlvDefinition(&b)
	e.genToStrDefinition(&b)

	return b.String()
}

func (e Enum) genToEnumDefinition(b *strings.Builder) {
	b.WriteString("genToEnumDefinition\n")
}

func (e Enum) genToSlvDefinition(b *strings.Builder) {
	b.WriteString("genToSlvDefinition\n")
}

func (e Enum) genToStrDefinition(b *strings.Builder) {
	b.WriteString("genToStrDefinition")
}
