package sym

import (
	"sort"
)

// Symbol interface represents generic symbol.
//
// Key is the key by which symbol must be searched in a symbol container.
// Key can differ from Name if language is case insensitive.
type Symbol interface {
	Filepath() string
	Files() []string
	Key() string
	Name() string
	LineNum() uint32
	Path() string                  // Must return full path to the symbol.
	InnerKeys() []string           // List of inner symbols keys.
	GetSymbol(key string) []Symbol // Get inner symbol.
	Doc() string
	Code() string
	DocCode() (string, string) // Get Doc and Code in one call, no need to read file twice.
	OneLineSummary() string
}

// ID is a unique symbol identifier.
// It is assumed, that multiple symbols with the same name
// can't be declared in the same line.
type ID struct {
	Key     string
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
		} else if s[i].Name() == s[j].Name() {
			if s[i].Filepath() < s[j].Filepath() {
				return true
			}
		}
		return false
	}

	sort.Slice(s, sortFunc)
}

// SortByName sorts Symbol slice by filepath in increasing order.
func SortByFilepath(s []Symbol) {
	sortFunc := func(i, j int) bool {
		if s[i].Filepath() < s[j].Filepath() {
			return true
		}
		return false
	}

	sort.Slice(s, sortFunc)
}

// UniqueNames returns a list containing only unique names with their count from the symbol list.
func UniqueNames(syms []Symbol) map[string]int {
	names := map[string]int{}
	for _, s := range syms {
		if _, ok := names[s.Name()]; !ok {
			names[s.Name()] = 1
		} else {
			names[s.Name()] += 1
		}
	}
	return names
}
