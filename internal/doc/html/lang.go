package html

import (
	_ "embed"
	"text/template"
)

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
