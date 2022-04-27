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
	tbs := []sym.Symbol{}
	pkgs := []sym.Symbol{}

	for _, syms := range l.Symbols() {
		for _, s := range syms {
			switch s.(type) {
			case Entity:
				if utils.IsTestbench(s.Key()) {
					tbs = append(tbs, s)
				} else {
					ents = append(ents, s)
				}
			case Package:
				pkgs = append(pkgs, s)
			case PackageInstantiation:
				pkgs = append(pkgs, s)
			default:
				panic("should never happen")
			}
		}
	}

	sym.SortByName(ents)
	sym.SortByName(pkgs)
	sym.SortByName(tbs)

	b := strings.Builder{}

	// Entities
	entNameLen := 0
	for _, e := range ents {
		if len(e.Name()) > entNameLen {
			entNameLen = len(e.Name())
		}
	}
	if len(ents) > 0 {
		b.WriteString(fmt.Sprintf("Entities (%d):\n", len(ents)))
	}
	for _, e := range ents {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", entNameLen, e.Name(), e.Filepath()),
		)
	}

	// Packages
	pkgNameLen := 0
	for _, p := range pkgs {
		if len(p.Name()) > pkgNameLen {
			pkgNameLen = len(p.Name())
		}
	}
	if len(pkgs) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	if len(pkgs) > 0 {
		b.WriteString(fmt.Sprintf("Packages (%d):\n", len(pkgs)))
	}
	for _, p := range pkgs {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", pkgNameLen, p.Name(), p.Filepath()),
		)
	}

	// Testbenches
	tbEntNameLen := 0
	for _, e := range tbs {
		if len(e.Name()) > tbEntNameLen {
			tbEntNameLen = len(e.Name())
		}
	}
	if len(tbs) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	if len(tbs) > 0 {
		b.WriteString(fmt.Sprintf("Testbenches (%d):\n", len(tbs)))
	}
	for _, e := range tbs {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", tbEntNameLen, e.Name(), e.Filepath()),
		)
	}

	return b.String()
}
