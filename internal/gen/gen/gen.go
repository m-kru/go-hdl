package gen

// In case of VHDL GenDeclaration must generate item declaration (aka. specification).
// GenDefinition must generate item body.
type Generable interface {
	GenDeclaration(args []string) string
	GenDefinition(args []string) string
	Name() string
	Width() int
}
