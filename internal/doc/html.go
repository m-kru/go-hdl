package doc

import (
	_ "embed"
	"fmt"
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"log"
	"os"
	"strings"
	"text/template"
)

func generateHTML(htmlArgs args.HTMLArgs) {
	if err := os.MkdirAll(htmlArgs.Path, 0775); err != nil {
		log.Fatalf("making html directory: %v", err)
	}

	generateCSS(htmlArgs)

	generateIndex(htmlArgs)
	generateVHDL(htmlArgs)
}

//go:embed templates/style.css
var cssStyleStr string
var cssStyleTmpl = template.Must(template.New("style.css").Parse(cssStyleStr))

type cssFormatters struct{}

func generateCSS(htmlArgs args.HTMLArgs) {
	err := os.MkdirAll(htmlArgs.Path+"css", 0775)
	if err != nil {
		log.Fatalf("making css directory: %v", err)
	}

	f, err := os.Create(htmlArgs.Path + "css/style.css")
	if err != nil {
		log.Fatalf("creating style.css file: %v", err)
	}
	defer f.Close()

	cssFmts := cssFormatters{}
	err = cssStyleTmpl.Execute(f, cssFmts)
	if err != nil {
		log.Fatalf("generating style.css file: %v", err)
	}
}

//go:embed templates/index.html
var indexStr string
var indexTmpl = template.Must(template.New("index.html").Parse(indexStr))

type indexFormatters struct {
	Copyright   string
	LibraryList string
	Title       string
	Topbar      string
}

func generateIndex(htmlArgs args.HTMLArgs) {
	f, err := os.Create(htmlArgs.Path + "index.html")
	if err != nil {
		log.Fatalf("creating index.html file: %v", err)
	}
	defer f.Close()

	libList := strings.Builder{}

	vhdlLibs := vhdl.LibraryNames()
	if len(vhdlLibs) > 0 {
		libList.WriteString("    <h2>VHDL</h2>\n    <ul class=\"symbol-list\">\n")
		for _, l := range vhdlLibs {
			libList.WriteString(
				fmt.Sprintf(
					"      <li><a href=\"vhdl/libs/%[1]s/index.html\">%[1]s</a></li>\n", l,
				),
			)
		}
		libList.WriteString("    </ul>")
	}

	indexFmts := indexFormatters{
		Copyright:   "",
		LibraryList: libList.String(),
		Title:       "THDL Documentation",
		Topbar:      topbar("home", 0),
	}

	err = indexTmpl.Execute(f, indexFmts)
	if err != nil {
		log.Fatalf("generating index.html file: %v", err)
	}
}

//go:embed templates/lang_index.html
var langIndexStr string
var langIndexTmpl = template.Must(template.New("lang_index.html").Parse(langIndexStr))

type LangFormatters struct {
	Copyright   string
	Language    string
	LibraryList string
	Title       string
	Topbar      string
}

func generateVHDL(htmlArgs args.HTMLArgs) {
	vhdlLibs := vhdl.LibraryNames()

	if len(vhdlLibs) == 0 {
		return
	}

	err := os.MkdirAll(htmlArgs.Path+"vhdl", 0775)
	if err != nil {
		log.Fatalf("making vhdl directory: %v", err)
	}

	generateVHDLIndex(htmlArgs)
}

func generateVHDLIndex(htmlArgs args.HTMLArgs) {
	vhdlLibs := vhdl.LibraryNames()
	libList := strings.Builder{}

	libList.WriteString("    <ul class=\"symbol-list\">\n")
	for _, l := range vhdlLibs {
		libList.WriteString(
			fmt.Sprintf(
				"      <li><a href=\"vhdl/libs/%[1]s/index.html\">%[1]s</a></li>\n", l,
			),
		)
	}
	libList.WriteString("    </ul>")

	langFmts := LangFormatters{
		Copyright:   "",
		Language:    "VHDL",
		LibraryList: libList.String(),
		Title:       "THDL Documentation",
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

func topbar(active string, nestingLevel int) string {
	homeActive := ""
	vhdlActive := ""

	switch active {
	case "home":
		homeActive = " active"
	case "vhdl":
		vhdlActive = " active"
	default:
		panic("should never happen")
	}

	root := "./"
	if nestingLevel == 1 {
		root = "../"
	} else if nestingLevel == 2 {
		root = "../../"
	} else if nestingLevel == 3 {
		root = "../../../"
	}

	b := strings.Builder{}

	b.WriteString(
		fmt.Sprintf("  <div class=\"topbar\">\n"+
			"    <div class=\"dropdown\">\n"+
			"      <button class=\"dropbtn%s\"><a href=\"%sindex.html\">Home</a></button>\n"+
			"    </div>\n", homeActive, root,
		),
	)

	vhdlLibs := vhdl.LibraryNames()
	if len(vhdlLibs) > 0 {
		b.WriteString(
			fmt.Sprintf("    <div class=\"dropdown\">\n"+
				"      <button class=\"dropbtn%s\"><a href=\"%svhdl/index.html\">VHDL</a></button>\n"+
				"      <div class=\"dropdown-content\">\n", vhdlActive, root,
			),
		)
		for _, l := range vhdlLibs {
			b.WriteString(
				fmt.Sprintf("        <a href=\"vhdl/libs/%[1]s/index.html\">%[1]s</a>\n", l),
			)
		}
		b.WriteString(`      </div>
    </div>`)
	}

	b.WriteString(`  </div>`)

	return b.String()
}
