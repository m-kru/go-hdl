package html

import (
	_ "embed"
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"log"
	"os"
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
	defer f.Close()

	err = langIndexTmpl.Execute(f, langFmts)
	if err != nil {
		log.Fatalf("generating index.html file: %v", err)
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

	ents, pkgs, tbs := vhdl.LibSortedSymbols(lib)

	tbNames := map[string]int{}
	for _, tb := range tbs {
		if _, ok := tbNames[tb.Name()]; !ok {
			tbNames[tb.Name()] = 1
		} else {
			tbNames[tb.Name()] += 1
		}
	}

	ulStr := "    <ul class=\"symbol-list\">\n"
	liStr := "<li><a href=\"%[1]s.html\">%[1]s</a></li>\n"
	ulEndStr := "    </ul>\n"

	if len(ents) > 0 {
		smblList.WriteString(fmt.Sprintf("    <h3>Entities (%d)</h3>\n", len(ents)))
		smblList.WriteString(ulStr)
		for _, e := range ents {
			smblList.WriteString(fmt.Sprintf(liStr, e.Name()))
		}
		smblList.WriteString(ulEndStr)
	}

	if len(pkgs) > 0 {
		smblList.WriteString(fmt.Sprintf("    <h3>Packages (%d)</h3>\n", len(pkgs)))
		smblList.WriteString(ulStr)
		for _, p := range pkgs {
			smblList.WriteString(fmt.Sprintf(liStr, p.Name()))
		}
		smblList.WriteString(ulEndStr)
	}

	if len(tbs) > 0 {
		smblList.WriteString(fmt.Sprintf("    <h3>Testbenches (%d)</h3>\n", len(tbs)))
		smblList.WriteString(ulStr)
		for _, t := range tbs {
			smblList.WriteString(fmt.Sprintf(liStr, t.Name()))
		}
		smblList.WriteString(ulEndStr)
	}

	libFmts := LibFormatters{
		Copyright:  htmlArgs.Copyright,
		Path:       "vhdl:" + name,
		SymbolList: smblList.String(),
		Title:      htmlArgs.Title,
		Topbar:     topbar("vhdl", 3),
	}

	f, err := os.Create(htmlArgs.Path + "vhdl/" + name + "/index.html")
	if err != nil {
		log.Fatalf("creating vhdl/%s/index.html file: %v", name, err)
	}
	defer f.Close()

	err = libIndexTmpl.Execute(f, libFmts)
	if err != nil {
		log.Fatalf("generating vhdl/%s/index.html file: %v", name, err)
	}
}
