package symbol

type Symbol interface {
	Name() string
	SymbolNames() []string                // Get names of all inner symbols.
	GetSymbol(name string) (Symbol, bool) // Get inner symbol.
	Doc() string
	Code() string
}
