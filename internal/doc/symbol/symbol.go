package symbol

import (
	"sort"
)

type Symbol interface {
	Filepath() string
	Name() string
	LineNum() uint32
	SymbolNames() []string          // Get names of all inner symbols.
	GetSymbol(name string) []Symbol // Get inner symbol.
	Doc() string
	Code() string
	DocCode() (string, string) // Get Doc and Code in one call, no need to read file twice.
}

// ID is a unique symbol identifier.
// It is assumed, that multiple symbols with the same name
// can't be declared in the same line.
type ID struct {
	Name    string
	LineNum uint32
}

// SortByLineNum sorts Symbol slice by line number in increasing order.
func SortByLineNum(s []Symbol) {
	sortFunc := func(i, j int) bool {
		if s[i].LineNum() < s[j].LineNum() {
			return true
		}
		return false
	}

	sort.Slice(s, sortFunc)
}

// SortByName sorts Symbol slice by name in increasing order.
func SortByName(s []Symbol) {
	sortFunc := func(i, j int) bool {
		if s[i].Name() < s[j].Name() {
			return true
		}
		return false
	}

	sort.Slice(s, sortFunc)
}
