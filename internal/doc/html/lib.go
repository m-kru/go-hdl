package html

import (
	_ "embed"
	"text/template"
)

//go:embed templates/lib_index.html
var libIndexStr string
var libIndexTmpl = template.Must(template.New("lib_index.html").Parse(libIndexStr))

type LibFormatters struct {
	Copyright  string
	Path       string
	SymbolList string
	Title      string
	Topbar     string
}
