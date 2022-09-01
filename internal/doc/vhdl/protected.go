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

type Protected struct {
	parent sym.Symbol

	filepath string
	key      string
	name     string
	lineNum  uint32

	docStart uint32
	docEnd   uint32

	codeStart uint32
	codeEnd   uint32

	Funcs map[sym.ID]sym.Symbol
	Procs map[sym.ID]sym.Symbol
}

func (p Protected) Filepath() string { return p.filepath }
func (p Protected) Key() string      { return p.key }
func (p Protected) Name() string     { return p.name }
func (p Protected) Files() []string  { panic("should never happen") }
func (p Protected) LineNum() uint32  { return p.lineNum }

func (p Protected) kind() string { return "protected" }

func (p Protected) AddSymbol(s sym.Symbol) error {
	id := sym.ID{Key: s.Key(), LineNum: s.LineNum()}

	switch s.(type) {
	case Function:
		p.Funcs[id] = s
	case Procedure:
		p.Procs[id] = s
	default:
		panic("should never happen")
	}

	return nil
}

func (p Protected) InnerKeys() []string {
	names := []string{}

	for id := range p.Funcs {
		names = append(names, id.Key)
	}
	for id := range p.Procs {
		names = append(names, id.Key)
	}

	return names
}

func (p Protected) GetSymbol(key string) []sym.Symbol {
	syms := []sym.Symbol{}

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

	return syms
}

func (p Protected) GetFunc(key string) []sym.Symbol {
	syms := []sym.Symbol{}

	for id, s := range p.Funcs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}

	return syms
}

func (p Protected) GetProc(key string) []sym.Symbol {
	syms := []sym.Symbol{}

	for id, s := range p.Procs {
		if id.Key == key {
			syms = append(syms, s)
		}
	}

	return syms
}

func (p Protected) Doc() string {
	f, err := os.ReadFile(p.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", p.filepath, err)
	}

	return string(f[p.docStart:p.docEnd])
}

func (p Protected) Code() string {
	b := strings.Builder{}

	// Functions.
	uniqueFuncs := map[string]bool{}
	for id := range p.Funcs {
		uniqueFuncs[id.Key] = true
	}
	funcs := []string{}
	for name := range uniqueFuncs {
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
				s = code
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
	for id := range p.Procs {
		uniqueProcs[id.Key] = true
	}
	procs := []string{}
	for name := range uniqueProcs {
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
				s = code
			} else {
				s = fmt.Sprintf("procedure %s ...\n", ps[0].Name())
			}
		} else {
			s = fmt.Sprintf("procedure %s ... (%d)\n", ps[0].Name(), len(ps))
		}
		b.WriteString(s)
	}

	return b.String()
}

func (p Protected) DocCode() (string, string) {
	f, err := os.ReadFile(p.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", p.filepath, err)
	}

	doc := string(f[p.docStart:p.docEnd])

	return doc, p.Code()
}

func (p Protected) OneLineSummary() string {
	f, err := os.ReadFile(p.filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", p.filepath, err)
	}

	code := utils.Dewhitespace(string(f[p.codeStart:p.codeEnd]))

	if utils.IsSingleLine(code) {
		return code
	}
	return fmt.Sprintf("%s ...\n", utils.FirstLine(code))
}

func (p Protected) Path() string {
	return p.parent.Path() + "." + p.name
}

func (p Protected) SortedFuncKeys() []string {
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

func (p Protected) SortedProcKeys() []string {
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
