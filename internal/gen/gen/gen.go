package gen

// In case of VHDL GenDeclaration must generate item declaration (aka. specification).
// GenDefinition must generate item body.
//
// Gens parameter in the GenDefinitions is the map of generables within the same scope.
// In case of VHDL, the scope is currently limited to the design unit.
type Generable interface {
	ParseArgs(args []string) error
	GenDeclarations() string
	GenDefinitions(gens Container) string
	Name() string
	Width() int
}
