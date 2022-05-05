package html

import (
	_ "embed"
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/index.html
var indexStr string
var indexTmpl = template.Must(template.New("index.html").Parse(indexStr))

type indexFormatters struct {
	Copyright   string
	LibraryList string
	Title       string
	Topbar      string
}

func generateIndex() {
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
					"      <li><a href=\"vhdl/%[1]s/index.html\">%[1]s</a></li>\n", l,
				),
			)
		}
		libList.WriteString("    </ul>")
	}

	indexFmts := indexFormatters{
		Copyright:   htmlArgs.Copyright,
		LibraryList: libList.String(),
		Title:       htmlArgs.Title,
		Topbar:      topbar("home", 0),
	}

	err = indexTmpl.Execute(f, indexFmts)
	if err != nil {
		log.Fatalf("generating index.html file: %v", err)
	}
}
