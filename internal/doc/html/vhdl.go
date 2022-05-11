package html

import (
	_ "embed"
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/sym"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
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
				"      <li><a href=\"%[1]s/index.html\">%[1]s</a></li>\n", l,
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
	generateVHDLLibSymbols(name)
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

func generateVHDLLibSymbols(name string) {
	lib, ok := vhdl.GetLibrary(name)
	if !ok {
		panic("should never happen")
	}

	for _, k := range lib.InnerKeys() {
		generateVHDLLibSymbol(lib, k)
	}
}

func generateVHDLLibSymbol(lib *lib.Library, key string) {
	content := strings.Builder{}

	syms := lib.GetSymbol(key)
	sym.SortByFilepath(syms)
	details := false
	if len(syms) > 1 {
		details = true
	}
	for _, s := range syms {
		switch s.(type) {
		case vhdl.Entity:
			generateVHDLEntityContent(s, details, &content)
		case vhdl.Package:
			genVHDLPkgContent(s, details, &content)
		default:
			//panic("should never happen")
		}
	}

	symFmts := SymbolFormatters{
		Copyright: htmlArgs.Copyright,
		Path:      fmt.Sprintf("vhdl:%s:%s", lib.Key(), key),
		Content:   content.String(),
		Title:     htmlArgs.Title,
		Topbar:    topbar("vhdl", 2),
	}

	path := fmt.Sprintf("%svhdl/%s/%s.html", htmlArgs.Path, lib.Name(), key)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("creating %s file: %v", path, err)
	}

	err = symbolTmpl.Execute(f, symFmts)
	if err != nil {
		log.Fatalf("generating %s file: %v", path, err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("closing %s file: %v", path, err)
	}
}

// If details is true, the content is put into the html <details> element.
func generateVHDLEntityContent(ent sym.Symbol, details bool, content *strings.Builder) {
	if details {
		content.WriteString("  <details>")
		content.WriteString(fmt.Sprintf("<summary class=\"filepath-summary\">%s</summary>", ent.Filepath()))
		content.WriteString("  <div class=\"details1\">")
	} else {
		content.WriteString(fmt.Sprintf("<p>%s</p>", ent.Filepath()))
	}

	content.WriteString(fmt.Sprintf("  <p class=\"doc\">%s</p>", ent.Doc()))
	content.WriteString(fmt.Sprintf("  <p class=\"code\">%s</p>", utils.VHDLHTMLBold(ent.Code())))

	if details {
		content.WriteString("  </div>")
		content.WriteString("  </details>")
	}
}

func genVHDLPkgContent(pkg sym.Symbol, details bool, content *strings.Builder) {
	if details {
		content.WriteString("  <details>\n")
		content.WriteString(fmt.Sprintf("<summary class=\"summary\">%s</summary>\n", pkg.Filepath()))
		content.WriteString("  <div class=\"details1\">\n")
	} else {
		content.WriteString(fmt.Sprintf("<p>%s</p>", pkg.Filepath()))
	}

	content.WriteString(
		fmt.Sprintf("  <p class=\"doc\">%s</p>\n", pkg.Doc()),
	)

	detailsLevel := 1
	if details {
		detailsLevel = 2
	}

	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Constants", detailsLevel, content)
	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Types", detailsLevel, content)
	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Subtypes", detailsLevel, content)

	if details {
		content.WriteString("  </div>\n")
		content.WriteString("  </details>\n")
	}
}

func genVHDLUniqueSymbolContent(sym sym.Symbol, summary string, detailsLevel int, content *strings.Builder) {
	doc := utils.Deindent(sym.Doc())
	code := utils.Deindent(sym.Code())

	isSingleLine := utils.IsSingleLine(code)

	if len(doc) > 0 || !isSingleLine {
		content.WriteString("  <details>\n")
		content.WriteString(
			fmt.Sprintf("<summary class=\"code-summary\">%s</summary>\n", utils.VHDLHTMLBold(summary)),
		)
		content.WriteString(
			fmt.Sprintf("  <div class=\"details%d\">\n", detailsLevel),
		)
		if len(doc) > 0 {
			content.WriteString(fmt.Sprintf("  <p class=\"doc\">%s</p>", doc))
		}
		if !isSingleLine {
			content.WriteString(fmt.Sprintf("  <p class=\"code\">%s</p>", utils.VHDLHTMLBold(code)))
		}
		content.WriteString("  </div>\n")
		content.WriteString("  </details>\n")
	} else {
		content.WriteString(
			fmt.Sprintf("<p class=\"code-summary summary-align\">%s</p>\n", " "+utils.VHDLHTMLBold(summary)),
		)
	}
}

func genVHDLPkgUniqueSymbolsContent(pkg vhdl.Package, class string, detailsLevel int, content *strings.Builder) {
	var keys []string
	switch class {
	case "Constants":
		keys = vhdl.PkgSortedConstKeys(pkg)
	case "Types":
		keys = vhdl.PkgSortedTypeKeys(pkg)
	case "Subtypes":
		keys = vhdl.PkgSortedSubtypeKeys(pkg)
	default:
		panic("should never happen")
	}

	if len(keys) > 0 {
		content.WriteString(fmt.Sprintf("  <h3>%s</h3>\n", class))
	}

	for _, key := range keys {
		sym := pkg.GetSymbol(key)[0]
		code := utils.Dewhitespace(sym.Code())
		var s string
		if utils.IsSingleLine(code) {
			s = fmt.Sprintf("%s", code)
		} else {
			s = fmt.Sprintf("%s ...\n", utils.FirstLine(code))
		}

		genVHDLUniqueSymbolContent(sym, s, detailsLevel, content)
	}
}
