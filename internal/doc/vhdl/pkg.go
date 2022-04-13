package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/sym"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"os"
	"sort"
	"strings"
)

type Package struct {
	filepath string
	key      string
	name     string
	lineNum  uint32

	docStart uint32
	docEnd   uint32

	codeStart uint32
	codeEnd   uint32

	Consts   map[sym.ID]sym.Symbol
	Funcs    map[sym.ID]sym.Symbol
	Procs    map[sym.ID]sym.Symbol
	Prots    map[sym.ID]sym.Symbol
	Types    map[sym.ID]sym.Symbol
	Subtypes map[sym.ID]sym.Symbol
}

func (p Package) Filepath() string { return p.filepath }
func (p Package) Key() string      { return p.key }
func (p Package) Name() string     { return p.name }
func (p Package) Files() []string  { panic("should never happen") }
func (p Package) LineNum() uint32  { return p.lineNum }

func (p Package) Doc() string {
	f, err := os.ReadFile(p.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", p.filepath, err)
	}

	return string(f[p.docStart:p.docEnd])
}

func (p Package) Code() string {
	b := strings.Builder{}

	// Constants.
	consts := []string{}
	for id, _ := range p.Consts {
		consts = append(consts, id.Key)
	}
	sort.Strings(consts)
	for _, key := range consts {
		c := p.GetSymbol(key)[0]
		code := utils.Dewhitespace(c.Code())
		var s string
		if utils.IsSingleLine(code) {
			s = fmt.Sprintf("%s", code)
		} else {
			s = fmt.Sprintf("%s ...\n", utils.FirstLine(code))
		}
		b.WriteString(s)
	}

	// Functions.
	uniqueFuncs := map[string]bool{}
	for id, _ := range p.Funcs {
		uniqueFuncs[id.Key] = true
	}
	funcs := []string{}
	for name, _ := range uniqueFuncs {
		funcs = append(funcs, name)
	}
	sort.Strings(funcs)
	if len(funcs) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range funcs {
		fs := p.GetFunc(key)
		var s string
		if len(fs) == 1 {
			code := utils.Dewhitespace(fs[0].Code())
			if utils.IsSingleLine(code) {
				s = fmt.Sprintf("%s", code)
			} else {
				s = fmt.Sprintf("%s ...\n", utils.FirstLine(code))
			}
		} else {
			impureCount := 0
			for _, f := range fs {
				if f.(Function).impure {
					impureCount += 1
				}
			}
			impurePrefix := ""
			if impureCount == len(fs) {
				impurePrefix = "impure "
			} else if impureCount > 0 {
				impurePrefix = "(impure)? "
			}
			s = fmt.Sprintf("%sfunction %s ... (%d)\n", impurePrefix, fs[0].Name(), len(fs))
		}
		b.WriteString(s)
	}

	// Procedures.
	uniqueProcs := map[string]bool{}
	for id, _ := range p.Procs {
		uniqueProcs[id.Key] = true
	}
	procs := []string{}
	for name, _ := range uniqueProcs {
		procs = append(procs, name)
	}
	sort.Strings(procs)
	if len(procs) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range procs {
		ps := p.GetProc(key)
		var s string
		if len(ps) == 1 {
			code := utils.Dewhitespace(ps[0].Code())
			if utils.IsSingleLine(code) {
				s = fmt.Sprintf("%s", code)
			} else {
				s = fmt.Sprintf("procedure %s ...\n", ps[0].Name())
			}
		} else {
			s = fmt.Sprintf("procedure %s ... (%d)\n", ps[0].Name(), len(ps))
		}
		b.WriteString(s)
	}

	// Types.
	types := []string{}
	for id, _ := range p.Types {
		types = append(types, id.Key)
	}
	sort.Strings(types)
	if len(types) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range types {
		t := p.GetSymbol(key)[0]
		code := utils.Dewhitespace(t.Code())
		var s string
		if utils.IsSingleLine(code) {
			s = fmt.Sprintf("%s", code)
		} else {
			s = fmt.Sprintf("%s ...\n", utils.FirstLine(code))
		}
		b.WriteString(s)
	}

	// Subtypes.
	subtypes := []string{}
	for id, _ := range p.Subtypes {
		subtypes = append(subtypes, id.Key)
	}
	sort.Strings(subtypes)
	if len(subtypes) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	for _, key := range subtypes {
		t := p.GetSymbol(key)[0]
		code := utils.Dewhitespace(t.Code())
		var s string
		if utils.IsSingleLine(code) {
			s = fmt.Sprintf("%s", code)
		} else {
			s = fmt.Sprintf("%s ...\n", utils.FirstLine(code))
		}
		b.WriteString(s)
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
	case Constant:
		p.Consts[id] = s
	case Function:
		p.Funcs[id] = s
	case Procedure:
		p.Procs[id] = s
	case Protected:
		p.Prots[id] = s
	case Type:
		p.Types[id] = s
	case Subtype:
		p.Subtypes[id] = s
	default:
		panic("should never happen")
	}

	return nil
}

func (p Package) InnerKeys() []string {
	names := []string{}

	for id, _ := range p.Consts {
		names = append(names, id.Key)
	}
	for id, _ := range p.Funcs {
		names = append(names, id.Key)
	}
	for id, _ := range p.Procs {
		names = append(names, id.Key)
	}
	for id, _ := range p.Prots {
		names = append(names, id.Key)
	}
	for id, _ := range p.Types {
		names = append(names, id.Key)
	}
	for id, _ := range p.Subtypes {
		names = append(names, id.Key)
	}

	return names
}

func (p Package) GetSymbol(key string) []sym.Symbol {
	syms := []sym.Symbol{}

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
	for id, s := range p.Prots {
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
