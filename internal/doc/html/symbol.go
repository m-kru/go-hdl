package html

import (
	_ "embed"
	"text/template"
)

//go:embed templates/symbol.html
var symbolStr string
var symbolTmpl = template.Must(template.New("symbol.html").Parse(symbolStr))

type SymbolFormatters struct {
	Copyright string
	Path      string
	Content   string
	Title     string
	Topbar    string
}
