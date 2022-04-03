package vhdl

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"strings"
)

func LibSummary(l *lib.Library) string {
	entities := []symbol.Symbol{}
	pkgs := []symbol.Symbol{}

	for _, s := range l.Symbols() {
		switch s.(type) {
		case Entity:
			entities = append(entities, s)
		case Package:
			pkgs = append(pkgs, s)
		default:
			panic("should never happen")
		}
	}

	symbol.SortByName(entities)
	symbol.SortByName(pkgs)

	entNameLen := 0
	for _, e := range entities {
		if len(e.Name()) > entNameLen {
			entNameLen = len(e.Name())
		}
	}

	b := strings.Builder{}

	if len(entities) > 0 {
		b.WriteString("Entities:\n")
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

	if len(pkgs) == 0 {
		return b.String()
	} else {
		b.WriteRune('\n')
	}
	b.WriteString("Packages:\n")

	for _, p := range pkgs {
		b.WriteString(fmt.Sprintf("  %-*s  %s\n", pkgNameLen, p.Name(), p.Filepath()))
	}

	return b.String()
}
