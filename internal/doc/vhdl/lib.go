package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/sym"
	"github.com/m-kru/go-thdl/internal/utils"
	"strings"
)

func LibSummary(l *lib.Library) string {
	ents := []sym.Symbol{}
	tbEnts := []sym.Symbol{}
	pkgDecs := []sym.Symbol{}
	pkgInsts := []sym.Symbol{}

	for _, syms := range l.Symbols() {
		for _, s := range syms {
			switch s.(type) {
			case Entity:
				if utils.IsTestbench(s.Key()) {
					tbEnts = append(tbEnts, s)
				} else {
					ents = append(ents, s)
				}
			case Package:
				pkgDecs = append(pkgDecs, s)
			case PackageInstantiation:
				pkgInsts = append(pkgInsts, s)
			default:
				panic("should never happen")
			}
		}
	}

	sym.SortByName(ents)
	sym.SortByName(pkgDecs)
	sym.SortByName(tbEnts)

	b := strings.Builder{}

	// Entity Declarations
	entNameLen := 0
	for _, e := range ents {
		if len(e.Name()) > entNameLen {
			entNameLen = len(e.Name())
		}
	}
	if len(ents) > 0 {
		b.WriteString(fmt.Sprintf("Entity Declarations (%d):\n", len(ents)))
	}
	for _, e := range ents {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", entNameLen, e.Name(), e.Filepath()),
		)
	}

	// Package Declarations
	pkgNameLen := 0
	for _, p := range pkgDecs {
		if len(p.Name()) > pkgNameLen {
			pkgNameLen = len(p.Name())
		}
	}
	if len(pkgDecs) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	if len(pkgDecs) > 0 {
		b.WriteString(fmt.Sprintf("Package Declarations (%d):\n", len(pkgDecs)))
	}
	for _, p := range pkgDecs {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", pkgNameLen, p.Name(), p.Filepath()),
		)
	}

	// Package Instantiations
	pkgNameLen = 0
	for _, p := range pkgInsts {
		if len(p.Name()) > pkgNameLen {
			pkgNameLen = len(p.Name())
		}
	}
	if len(pkgInsts) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	if len(pkgInsts) > 0 {
		b.WriteString(fmt.Sprintf("Package Instantiations (%d):\n", len(pkgInsts)))
	}
	for _, p := range pkgInsts {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", pkgNameLen, p.Name(), p.Filepath()),
		)
	}

	// Testbench Entities
	tbEntNameLen := 0
	for _, e := range tbEnts {
		if len(e.Name()) > tbEntNameLen {
			tbEntNameLen = len(e.Name())
		}
	}
	if len(tbEnts) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	if len(tbEnts) > 0 {
		b.WriteString(fmt.Sprintf("Testbench Entities (%d):\n", len(tbEnts)))
	}
	for _, e := range tbEnts {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", tbEntNameLen, e.Name(), e.Filepath()),
		)
	}

	return b.String()
}
