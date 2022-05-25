package html

import (
	_ "embed"
	"log"
	"os"
	"path"
	"text/template"
)

//go:embed templates/style.css
var cssStyleStr string
var cssStyleTmpl = template.Must(template.New("style.css").Parse(cssStyleStr))

type cssFormatters struct{}

func genCSS() {
	err := os.MkdirAll(path.Join(htmlArgs.Path, "css"), 0775)
	if err != nil {
		log.Fatalf("making css directory: %v", err)
	}

	f, err := os.Create(path.Join(htmlArgs.Path, "css/style.css"))
	if err != nil {
		log.Fatalf("creating style.css file: %v", err)
	}

	cssFmts := cssFormatters{}
	err = cssStyleTmpl.Execute(f, cssFmts)
	if err != nil {
		log.Fatalf("generating style.css file: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("closing style.css file: %v", err)
	}
}
