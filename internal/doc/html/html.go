package html

import (
	_ "embed"
	"github.com/m-kru/go-hdl/internal/args"
	"log"
	"os"
)

var htmlArgs args.HTMLArgs

func Generate(args args.HTMLArgs) {
	htmlArgs = args

	if err := os.MkdirAll(htmlArgs.Path, 0775); err != nil {
		log.Fatalf("making html directory: %v", err)
	}

	genCSS()

	if htmlArgs.Copyright != "" {
		htmlArgs.Copyright = "&copy; " + htmlArgs.Copyright
	}
	if htmlArgs.Title == "" {
		htmlArgs.Title = "HDL Documentation"
	}

	genIndex()
	genVHDL()
}
