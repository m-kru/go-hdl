package vhdl

import (
	"github.com/m-kru/go-thdl/internal/doc/sym"
	"log"
	"os"
	"sort"
	"strings"
)

type Package struct {
	parent sym.Symbol

	filepath string
	key      string
	name     string
	lineNum  uint32

	docStart uint32
	docEnd   uint32

	codeStart uint32
	codeEnd   uint32

	Aliases   map[sym.ID]sym.Symbol
	Consts    map[sym.ID]sym.Symbol
	Funcs     map[sym.ID]sym.Symbol
	Procs     map[sym.ID]sym.Symbol
	Types     map[sym.ID]sym.Symbol
	Subtypes  map[sym.ID]sym.Symbol
	Variables map[sym.ID]sym.Symbol
}

func (p Package) Filepath() string       { return p.filepath }
func (p Package) Key() string            { return p.key }
func (p Package) Name() string           { return p.name }
func (p Package) Files() []string        { panic("should never happen") }
func (p Package) LineNum() uint32        { return p.lineNum }
func (p Package) OneLineSummary() string { panic("not yet implemented") }

func (p Package) Doc() string {
	f, err := os.ReadFile(p.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", p.filepath, err)
	}

	return string(f[p.docStart:p.docEnd])
}

// PkgSortedAliasKeys returns alias keys in alphabetical order.
func PkgSortedAliasKeys(p Package) []string {
	vars := []string{}
	for id := range p.Aliases {
		vars = append(vars, id.Key)
	}
	sort.Strings(vars)
	return vars
}

// PkgSortedConstKeys returns constant keys in alphabetical order.
func PkgSortedConstKeys(p Package) []string {
	consts := []string{}
	for id := range p.Consts {
		consts = append(consts, id.Key)
	}
	sort.Strings(consts)
	return consts
}

func (p Package) SortedFuncKeys() []string {
	uniqueFuncs := map[string]bool{}
	for id := range p.Funcs {
		uniqueFuncs[id.Key] = true
	}
	funcs := []string{}
	for name := range uniqueFuncs {
		funcs = append(funcs, name)
	}
	sort.Strings(funcs)
	return funcs
}

func (p Package) SortedProcKeys() []string {
	uniqueProcs := map[string]bool{}
	for id := range p.Procs {
		uniqueProcs[id.Key] = true
	}
	procs := []string{}
	for name := range uniqueProcs {
		procs = append(procs, name)
	}
	sort.Strings(procs)
	return procs
}

// PkgSortedTypeKeys returns type keys in alphabetical order.
func PkgSortedTypeKeys(p Package) []string {
	types := []string{}
	for id := range p.Types {
		types = append(types, id.Key)
	}
	sort.Strings(types)
	return types
}

// PkgSortedSubtypeKeys returns subtype keys in alphabetical order.
func PkgSortedSubtypeKeys(p Package) []string {
	subtypes := []string{}
	for id := range p.Subtypes {
		subtypes = append(subtypes, id.Key)
	}
	sort.Strings(subtypes)
	return subtypes
}

// PkgSortedVariableKeys returns variable keys in alphabetical order.
func PkgSortedVariableKeys(p Package) []string {
	vars := []string{}
	for id := range p.Variables {
		vars = append(vars, id.Key)
	}
	sort.Strings(vars)
	return vars
}

func (p Package) Code() string {
	b := strings.Builder{}

	// Aliases.
	aliases := PkgSortedAliasKeys(p)
	for _, key := range aliases {
		a := p.GetSymbol(key)[0]
		b.WriteString(a.OneLineSummary())
	}

	// Constants.
	consts := PkgSortedConstKeys(p)
	if len(consts) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range consts {
		c := p.GetSymbol(key)[0]
		b.WriteString(c.OneLineSummary())
	}

	// Functions.
	funcs := p.SortedFuncKeys()
	if len(funcs) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range funcs {
		fs := p.GetFunc(key)
		b.WriteString(FuncsCodeSummary(fs))
	}

	// Procedures.
	procs := p.SortedProcKeys()
	if len(procs) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range procs {
		ps := p.GetProc(key)
		b.WriteString(ProcsCodeSummary(ps))
	}

	// Types.
	types := PkgSortedTypeKeys(p)
	if len(types) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range types {
		t := p.GetSymbol(key)[0]
		b.WriteString(t.OneLineSummary())
	}

	// Subtypes.
	subtypes := PkgSortedSubtypeKeys(p)
	if len(subtypes) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range subtypes {
		s := p.GetSymbol(key)[0]
		b.WriteString(s.OneLineSummary())
	}

	// Variables.
	vars := PkgSortedVariableKeys(p)
	if len(vars) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range vars {
		v := p.GetSymbol(key)[0]
		b.WriteString(v.OneLineSummary())
	}

	return b.String()
}

func (p Package) DocCode() (string, string) {
	f, err := os.ReadFile(p.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", p.filepath, err)
	}

	doc := string(f[p.docStart:p.docEnd])

	return doc, p.Code()
}

func (p Package) AddSymbol(s sym.Symbol) error {
	id := sym.ID{Key: s.Key(), LineNum: s.LineNum()}

	switch s.(type) {
	case Alias:
		p.Aliases[id] = s
	case Constant:
		p.Consts[id] = s
	case Function:
		p.Funcs[id] = s
	case Procedure:
		p.Procs[id] = s
	case Type:
		p.Types[id] = s
	case Subtype:
		p.Subtypes[id] = s
	case Variable:
		p.Variables[id] = s
	default:
		panic("should never happen")
	}

	return nil
}

func (p Package) InnerKeys() []string {
	names := []string{}

	for id := range p.Aliases {
		names = append(names, id.Key)
	}
	for id := range p.Consts {
		names = append(names, id.Key)
	}
	for id := range p.Funcs {
		names = append(names, id.Key)
	}
	for id := range p.Procs {
		names = append(names, id.Key)
	}
	for id := range p.Types {
		names = append(names, id.Key)
	}
	for id := range p.Subtypes {
		names = append(names, id.Key)
	}
	for id := range p.Variables {
		names = append(names, id.Key)
	}

	return names
}

func (p Package) GetSymbol(key string) []sym.Symbol {
	syms := []sym.Symbol{}

	for id, s := range p.Aliases {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Consts {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Funcs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Procs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Types {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Subtypes {
		if id.Key == key {
			syms = append(syms, s)
		}
	}
	for id, s := range p.Variables {
		if id.Key == key {
			syms = append(syms, s)
		}
	}

	return syms
}

func (p Package) GetFunc(key string) []sym.Symbol {
	syms := []sym.Symbol{}

	for id, s := range p.Funcs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}

	return syms
}

func (p Package) GetProc(key string) []sym.Symbol {
	syms := []sym.Symbol{}

	for id, s := range p.Procs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}

	return syms
}

func (p Package) Path() string {
	return p.parent.Path() + "." + p.name
}
