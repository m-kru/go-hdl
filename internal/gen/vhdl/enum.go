package vhdl

import (
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
	b.WriteString("genToEnumDeclaration\n")
}

func (e enum) genToSlvDeclaration(b *strings.Builder) {
	b.WriteString("genToSlvDeclaration\n")
}

func (e enum) genToStrDeclaration(b *strings.Builder) {
	b.WriteString("genToStrDeclaration\n")
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
