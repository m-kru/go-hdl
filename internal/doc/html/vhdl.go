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
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
)

func genVHDL() {
	vhdlLibs := vhdl.LibraryNames()

	if len(vhdlLibs) == 0 {
		return
	}

	err := os.MkdirAll(path.Join(htmlArgs.Path, "vhdl"), 0775)
	if err != nil {
		log.Fatalf("making vhdl directory: %v", err)
	}

	genVHDLIndex()
	genVHDLLibs()
}

func genVHDLIndex() {
	vhdlLibs := vhdl.LibraryNames()
	libList := strings.Builder{}

	libList.WriteString("<ul class=\"symbol-list\">")
	for _, l := range vhdlLibs {
		libList.WriteString(
			fmt.Sprintf("<li><a href=\"%[1]s/index.html\">%[1]s</a></li>", l),
		)
	}
	libList.WriteString("</ul>")

	langFmts := LangFormatters{
		Copyright:   htmlArgs.Copyright,
		Language:    "VHDL",
		LibraryList: libList.String(),
		Title:       htmlArgs.Title,
		Topbar:      topbar("vhdl", 1),
	}

	f, err := os.Create(path.Join(htmlArgs.Path, "vhdl", "index.html"))
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

func genVHDLLibs() {
	for _, l := range vhdl.LibraryNames() {
		genVHDLLib(l)
	}
}

func genVHDLLib(name string) {
	err := os.MkdirAll(path.Join(htmlArgs.Path, "vhdl", name), 0775)
	if err != nil {
		log.Fatalf("making vhdl/%s directory: %v", name, err)
	}

	genVHDLLibIndex(name)
	genVHDLLibSymbols(name)
}

func genVHDLLibIndex(name string) {
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

	ulStr := "<ul class=\"symbol-list\">"
	liStr := "<li><a href=\"%s.html\">%s%s</a></li>"
	ulEndStr := "</ul>"

	if len(ents) > 0 {
		smblList.WriteString(fmt.Sprintf("<h3>Entities (%d)</h3>", len(ents)))
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
		smblList.WriteString(fmt.Sprintf("<h3>Packages (%d)</h3>", len(pkgs)))
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
		smblList.WriteString(fmt.Sprintf("<h3>Testbenches (%d)</h3>", len(tbs)))
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
		Path:       lib.Path(),
		SymbolList: smblList.String(),
		Title:      htmlArgs.Title,
		Topbar:     topbar("vhdl", 2),
	}

	f, err := os.Create(path.Join(htmlArgs.Path, "vhdl", name, "index.html"))
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

func genVHDLLibSymbols(name string) {
	lib, ok := vhdl.GetLibrary(name)
	if !ok {
		panic("should never happen")
	}

	for _, k := range lib.InnerKeys() {
		genVHDLLibSymbol(lib, k)
	}
}

func genVHDLLibSymbol(lib *lib.Library, key string) {
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
			genVHDLEntityContent(s, details, &content)
		case vhdl.Package:
			genVHDLPkgContent(s, details, &content)
		case vhdl.PackageInstantiation:
			genVHDLPkgInstContent(s, details, &content)
		default:
			panic("should never happen")
		}
	}

	symPath := syms[0].Path()
	if len(syms) > 1 {
		symPath = fmt.Sprintf("%s.%s", lib.Path(), key)
	}

	symFmts := SymbolFormatters{
		Copyright: htmlArgs.Copyright,
		Path:      symPath,
		Content:   content.String(),
		Title:     htmlArgs.Title,
		Topbar:    topbar("vhdl", 2),
	}

	filePath := path.Join(htmlArgs.Path, "vhdl", lib.Name(), fmt.Sprintf("%s.html", key))
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("creating %s file: %v", filePath, err)
	}

	err = symbolTmpl.Execute(f, symFmts)
	if err != nil {
		log.Fatalf("generating %s file: %v", filePath, err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("closing %s file: %v", filePath, err)
	}
}

// If details is true, the content is put into the html <details> element.
func genVHDLEntityContent(ent sym.Symbol, details bool, b *strings.Builder) {
	if details {
		b.WriteString("<details>")
		b.WriteString(
			fmt.Sprintf("<summary class=\"filepath-summary\">%s</summary>", ent.Filepath()),
		)
		b.WriteString("<div class=\"details\">")
	} else {
		b.WriteString(fmt.Sprintf("<p>%s</p>", ent.Filepath()))
	}

	b.WriteString(fmt.Sprintf("<p class=\"doc\">%s</p>", utils.VHDLDeindentDecomment(ent.Doc())))
	b.WriteString(fmt.Sprintf("<p class=\"code\">%s</p>", utils.VHDLHTMLBold(ent.Code())))

	if details {
		b.WriteString("  </div></details>")
	}
}

func genVHDLPkgContent(pkg sym.Symbol, details bool, b *strings.Builder) {
	if details {
		b.WriteString("<details>")
		b.WriteString(fmt.Sprintf("<summary class=\"filepath-summary\">%s</summary>", pkg.Filepath()))
		b.WriteString("<div class=\"details\">")
	} else {
		b.WriteString(fmt.Sprintf("<p>%s</p>", pkg.Filepath()))
	}

	b.WriteString(fmt.Sprintf("<h3>Package %s</h3>", pkg.Name()))
	b.WriteString(fmt.Sprintf("<p class=\"doc\">%s</p>", utils.VHDLDeindentDecomment(pkg.Doc())))

	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Aliases", b)
	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Constants", b)
	genVHDLOverloadedSymbolsContent(pkg.(vhdl.Package), "Functions", b)
	genVHDLOverloadedSymbolsContent(pkg.(vhdl.Package), "Procedures", b)
	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Types", b)
	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Subtypes", b)
	genVHDLPkgUniqueSymbolsContent(pkg.(vhdl.Package), "Variables", b)

	if details {
		b.WriteString("</div></details>")
	}
}

func genVHDLPkgInstContent(pkg sym.Symbol, details bool, b *strings.Builder) {
	if details {
		b.WriteString("<details>")
		b.WriteString(fmt.Sprintf("<summary class=\"filepath-summary\">%s</summary>", pkg.Filepath()))
		b.WriteString("<div class=\"details\">")
	} else {
		b.WriteString(fmt.Sprintf("<p>%s</p>", pkg.Filepath()))
	}

	b.WriteString(fmt.Sprintf("<h3>Package %s</h3>", pkg.Name()))
	b.WriteString(fmt.Sprintf("<p class=\"doc\">%s</p>", utils.VHDLDeindentDecomment(pkg.Doc())))
	b.WriteString(fmt.Sprintf("<p class=\"code\">%s</p>", utils.VHDLHTMLBold(pkg.Code())))

	if details {
		b.WriteString("</div></details>")
	}
}

func genVHDLProtectedType(prot vhdl.Protected, rand uint32) {
	b := strings.Builder{}

	b.WriteString(
		fmt.Sprintf(
			"<h3>Protected %s</h3>"+
				"<p class=\"doc\">%s</p>",
			prot.Name(), utils.VHDLDeindentDecomment(prot.Doc()),
		),
	)

	genVHDLOverloadedSymbolsContent(prot, "Functions", &b)
	genVHDLOverloadedSymbolsContent(prot, "Procedures", &b)

	symFmts := SymbolFormatters{
		Copyright: htmlArgs.Copyright,
		Path:      prot.Path(),
		Content:   b.String(),
		Title:     htmlArgs.Title,
		Topbar:    topbar("vhdl", 2),
	}

	filePath := prot.Path()
	filePath = strings.Replace(filePath, ":", ".", -1)
	elems := strings.Split(filePath, ".")
	elems = elems[0 : len(elems)-2]
	elems = append(elems, prot.Key())
	filePath = path.Join(elems...)

	filePath = fmt.Sprintf("%s_%d.html", filePath, rand)
	f, err := os.Create(path.Join(htmlArgs.Path, filePath))
	if err != nil {
		log.Fatalf("creating %s file: %v", filePath, err)
	}

	err = symbolTmpl.Execute(f, symFmts)
	if err != nil {
		log.Fatalf("generating %s file: %v", filePath, err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("closing %s file: %v", filePath, err)
	}
}

func genVHDLUniqueSymbolContent(pkg vhdl.Package, key string, b *strings.Builder) {
	sym := pkg.GetSymbol(key)[0]
	summary := sym.OneLineSummary()

	doc := utils.Deindent(sym.Doc())
	code := utils.Deindent(sym.Code())

	isSingleLine := utils.IsSingleLine(code)

	aPrefix := ""
	aSuffix := ""
	if _, ok := sym.(vhdl.Protected); ok {
		rand := rand.Uint32()
		genVHDLProtectedType(sym.(vhdl.Protected), rand)
		aPrefix = fmt.Sprintf("<a href=\"%s_%d.html\">", sym.Key(), rand)
		aSuffix = "</a>"
	}

	if len(doc) > 0 || !isSingleLine {
		b.WriteString("<details>")
		b.WriteString(
			fmt.Sprintf(
				"<summary class=\"code-summary\">%s%s%s</summary>",
				aPrefix, utils.VHDLHTMLBold(summary), aSuffix,
			),
		)
		b.WriteString("<div class=\"details\">")
		if len(doc) > 0 {
			b.WriteString(fmt.Sprintf("<p class=\"doc\">%s</p>", utils.VHDLDecomment(doc)))
		}
		if !isSingleLine {
			b.WriteString(fmt.Sprintf("<p class=\"code\">%s</p>", utils.VHDLHTMLBold(code)))
		}
		b.WriteString("</div></details>")
	} else {
		b.WriteString(
			fmt.Sprintf("<p class=\"code-summary summary-align\">%s</p>", " "+utils.VHDLHTMLBold(summary)),
		)
	}
}

func genVHDLPkgUniqueSymbolsContent(pkg vhdl.Package, class string, content *strings.Builder) {
	var keys []string
	switch class {
	case "Aliases":
		keys = vhdl.PkgSortedAliasKeys(pkg)
	case "Constants":
		keys = vhdl.PkgSortedConstKeys(pkg)
	case "Types":
		keys = vhdl.PkgSortedTypeKeys(pkg)
	case "Subtypes":
		keys = vhdl.PkgSortedSubtypeKeys(pkg)
	case "Variables":
		keys = vhdl.PkgSortedVariableKeys(pkg)
	default:
		panic("should never happen")
	}

	if len(keys) > 0 {
		content.WriteString(fmt.Sprintf("<h4>%s</h4>", class))
	}

	for _, key := range keys {
		genVHDLUniqueSymbolContent(pkg, key, content)
	}
}

func genVHDLOverloadedSymbolContent(syms []sym.Symbol, summary string, b *strings.Builder) {
	if len(syms) == 1 && utils.IsSingleLine(syms[0].Code()) && len(syms[0].Doc()) == 0 {
		b.WriteString(
			fmt.Sprintf("<p class=\"code-summary summary-align\">%s</p>", " "+utils.VHDLHTMLBold(summary)),
		)
	} else if len(syms) == 1 && utils.IsSingleLine(syms[0].Code()) {
		doc := syms[0].Doc()

		b.WriteString("<details>")
		b.WriteString(fmt.Sprintf("<summary class=\"code-summary\">%s</summary>", utils.VHDLHTMLBold(summary)))
		b.WriteString("<div class=\"details\">")
		if len(doc) > 0 {
			b.WriteString(fmt.Sprintf("  <p class=\"doc\">%s</p>", utils.VHDLDeindentDecomment(doc)))
		}
		b.WriteString("</div>")
		b.WriteString("</details>")
	} else {
		sym.SortByLineNum(syms)
		b.WriteString("<details>")
		b.WriteString(fmt.Sprintf("<summary class=\"code-summary\">%s</summary>", utils.VHDLHTMLBold(summary)))
		b.WriteString("<div class=\"details\">")

		for _, sym := range syms {
			doc := sym.Doc()
			code := utils.Deindent(sym.Code())
			if len(doc) > 0 {
				b.WriteString(fmt.Sprintf("<p class=\"doc\">%s</p>", utils.VHDLDeindentDecomment(doc)))
			}
			b.WriteString(fmt.Sprintf("<p class=\"code\">%s</p>", utils.VHDLHTMLBold(code)))
		}

		b.WriteString("</div></details>")
	}
}

func genVHDLOverloadedSymbolsContent(sc vhdl.SubprogramsContainer, class string, content *strings.Builder) {
	var keys []string
	switch class {
	case "Functions":
		keys = sc.SortedFuncKeys()
	case "Procedures":
		keys = sc.SortedProcKeys()
	default:
		panic("should never happen")
	}

	if len(keys) > 0 {
		content.WriteString(fmt.Sprintf("<h4>%s</h4>", class))
	}

	var summary string
	var syms []sym.Symbol
	for _, key := range keys {
		switch class {
		case "Functions":
			syms = sc.GetFunc(key)
			summary = vhdl.FuncsCodeSummary(syms)
		case "Procedures":
			syms = sc.GetProc(key)
			summary = vhdl.ProcsCodeSummary(syms)
		default:
			panic("should never happen")
		}

		genVHDLOverloadedSymbolContent(syms, summary, content)
	}
}
