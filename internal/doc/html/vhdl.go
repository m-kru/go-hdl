package html

import (
	_ "embed"
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/sym"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"golang.org/x/exp/maps"
	"log"
	"os"
	"sort"
	"strings"
)

func generateVHDL() {
	vhdlLibs := vhdl.LibraryNames()

	if len(vhdlLibs) == 0 {
		return
	}

	err := os.MkdirAll(htmlArgs.Path+"vhdl", 0775)
	if err != nil {
		log.Fatalf("making vhdl directory: %v", err)
	}

	generateVHDLIndex()
	generateVHDLLibs()
}

func generateVHDLIndex() {
	vhdlLibs := vhdl.LibraryNames()
	libList := strings.Builder{}

	libList.WriteString("    <ul class=\"symbol-list\">\n")
	for _, l := range vhdlLibs {
		libList.WriteString(
			fmt.Sprintf(
				"      <li><a href=\"vhdl/%[1]s/index.html\">%[1]s</a></li>\n", l,
			),
		)
	}
	libList.WriteString("    </ul>")

	langFmts := LangFormatters{
		Copyright:   htmlArgs.Copyright,
		Language:    "VHDL",
		LibraryList: libList.String(),
		Title:       htmlArgs.Title,
		Topbar:      topbar("vhdl", 1),
	}

	f, err := os.Create(htmlArgs.Path + "vhdl/index.html")
	if err != nil {
		log.Fatalf("creating vhdl/index.html file: %v", err)
	}

	err = langIndexTmpl.Execute(f, langFmts)
	if err != nil {
		log.Fatalf("generating vhdl/index.html file: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("closing vhdl/index.html file: %v", err)
	}
}

func generateVHDLLibs() {
	for _, l := range vhdl.LibraryNames() {
		generateVHDLLib(l)
	}
}

func generateVHDLLib(name string) {
	err := os.MkdirAll(htmlArgs.Path+"vhdl/"+name, 0775)
	if err != nil {
		log.Fatalf("making vhdl/%s directory: %v", name, err)
	}

	generateVHDLLibIndex(name)
}

func generateVHDLLibIndex(name string) {
	smblList := strings.Builder{}

	lib, ok := vhdl.GetLibrary(name)
	if !ok {
		panic("should never happen")
	}

	ents, pkgs, tbs := vhdl.LibSymbols(lib)

	entUniqueNames := sym.UniqueNames(ents)
	entNames := maps.Keys(entUniqueNames)
	sort.Strings(entNames)

	pkgUniqueNames := sym.UniqueNames(pkgs)
	pkgNames := maps.Keys(pkgUniqueNames)
	sort.Strings(pkgNames)

	tbUniqueNames := sym.UniqueNames(tbs)
	tbNames := maps.Keys(tbUniqueNames)
	sort.Strings(tbNames)

	ulStr := "    <ul class=\"symbol-list\">\n"
	liStr := "<li><a href=\"%s.html\">%s%s</a></li>\n"
	ulEndStr := "    </ul>\n"

	if len(ents) > 0 {
		smblList.WriteString(fmt.Sprintf("    <h3>Entities (%d)</h3>\n", len(ents)))
		smblList.WriteString(ulStr)
		for _, e := range entNames {
			count := ""
			if entUniqueNames[e] > 1 {
				count = fmt.Sprintf(" (%d)", entUniqueNames[e])
			}
			smblList.WriteString(fmt.Sprintf(liStr, strings.ToLower(e), e, count))
		}
		smblList.WriteString(ulEndStr)
	}

	if len(pkgs) > 0 {
		smblList.WriteString(fmt.Sprintf("    <h3>Packages (%d)</h3>\n", len(pkgs)))
		smblList.WriteString(ulStr)
		for _, p := range pkgNames {
			count := ""
			if pkgUniqueNames[p] > 1 {
				count = fmt.Sprintf(" (%d)", pkgUniqueNames[p])
			}
			smblList.WriteString(fmt.Sprintf(liStr, strings.ToLower(p), p, count))
		}
		smblList.WriteString(ulEndStr)
	}

	if len(tbs) > 0 {
		smblList.WriteString(fmt.Sprintf("    <h3>Testbenches (%d)</h3>\n", len(tbs)))
		smblList.WriteString(ulStr)
		for _, t := range tbNames {
			count := ""
			if tbUniqueNames[t] > 1 {
				count = fmt.Sprintf(" (%d)", tbUniqueNames[t])
			}
			smblList.WriteString(fmt.Sprintf(liStr, strings.ToLower(t), t, count))
		}
		smblList.WriteString(ulEndStr)
	}

	libFmts := LibFormatters{
		Copyright:  htmlArgs.Copyright,
		Path:       "vhdl:" + name,
		SymbolList: smblList.String(),
		Title:      htmlArgs.Title,
		Topbar:     topbar("vhdl", 2),
	}

	f, err := os.Create(htmlArgs.Path + "vhdl/" + name + "/index.html")
	if err != nil {
		log.Fatalf("creating vhdl/%s/index.html file: %v", name, err)
	}

	err = libIndexTmpl.Execute(f, libFmts)
	if err != nil {
		log.Fatalf("generating vhdl/%s/index.html file: %v", name, err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("closing vhdl/%s/index.html file: %v", name, err)
	}
}
