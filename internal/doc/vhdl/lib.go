package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"github.com/m-kru/go-thdl/internal/utils"
	"strings"
)

func LibSummary(l *lib.Library) string {
	entities := []symbol.Symbol{}
	tbEntities := []symbol.Symbol{}
	pkgs := []symbol.Symbol{}

	for _, syms := range l.Symbols() {
		for _, s := range syms {
			switch s.(type) {
			case Entity:
				if utils.IsTestbench(s.Name()) {
					tbEntities = append(tbEntities, s)
				} else {
					entities = append(entities, s)
				}
			case Package:
				pkgs = append(pkgs, s)
			default:
				panic("should never happen")
			}
		}
	}

	symbol.SortByName(entities)
	symbol.SortByName(pkgs)
	symbol.SortByName(tbEntities)

	b := strings.Builder{}

	entNameLen := 0
	for _, e := range entities {
		if len(e.Name()) > entNameLen {
			entNameLen = len(e.Name())
		}
	}
	if len(entities) > 0 {
		b.WriteString(fmt.Sprintf("Entities (%d):\n", len(entities)))
	}
	for _, e := range entities {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", entNameLen, e.Name(), e.Filepath()),
		)
	}

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

	tbEntNameLen := 0
	for _, e := range tbEntities {
		if len(e.Name()) > tbEntNameLen {
			tbEntNameLen = len(e.Name())
		}
	}
	if len(tbEntities) > 0 && b.Len() > 0 {
		b.WriteRune('\n')
	}
	if len(tbEntities) > 0 {
		b.WriteString(fmt.Sprintf("Testbench Entities (%d):\n", len(tbEntities)))
	}
	for _, e := range tbEntities {
		b.WriteString(
			fmt.Sprintf("  %-*s  %s\n", tbEntNameLen, e.Name(), e.Filepath()),
		)
	}

	return b.String()
}
