package html

import (
	_ "embed"
	"github.com/m-kru/go-thdl/internal/args"
	"log"
	"os"
)

func Generate(htmlArgs args.HTMLArgs) {
	if err := os.MkdirAll(htmlArgs.Path, 0775); err != nil {
		log.Fatalf("making html directory: %v", err)
	}

	generateCSS(htmlArgs)

	if htmlArgs.Copyright != "" {
		htmlArgs.Copyright = "&copy; " + htmlArgs.Copyright
	}
	if htmlArgs.Title == "" {
		htmlArgs.Title = "THDL Documentation"
	}

	generateIndex(htmlArgs)
	generateVHDL(htmlArgs)
}
