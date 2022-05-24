package gen

// In case of VHDL GenDeclaration must generate item declaration (aka. specification).
// GenDefinition must generate item body.
//
// Gens parameter in the GenDeclarations and GenDefinitions is the map of generables
// within the same scope. In case of VHDL, the scope is limited to the design unit.
type Generable interface {
	ParseArgs(args []string) error
	GenDeclarations(gens map[string]Generable) string
	GenDefinitions(gens map[string]Generable) string
	Name() string
	Width() int
}
