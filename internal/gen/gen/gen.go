package gen

// In case of VHDL GenDeclaration must generate item declaration (aka. specification).
// GenDefinition must generate item body.
type Generable interface {
	ParseArgs(args []string) error
	GenDeclarations() string
	GenDefinitions() string
	Name() string
	Width() int
}
