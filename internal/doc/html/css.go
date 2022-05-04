package html

import (
	_ "embed"
	"log"
	"os"
	"text/template"
)

//go:embed templates/style.css
var cssStyleStr string
var cssStyleTmpl = template.Must(template.New("style.css").Parse(cssStyleStr))

type cssFormatters struct{}

func generateCSS() {
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
