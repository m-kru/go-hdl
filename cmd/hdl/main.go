package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m-kru/go-hdl/internal/args"
	"github.com/m-kru/go-hdl/internal/doc"
	"github.com/m-kru/go-hdl/internal/gen"
	"github.com/m-kru/go-hdl/internal/vet"
	"github.com/m-kru/go-hdl/internal/vet/rprt"
)

var printDebug bool = false

type Logger struct{}

func (l Logger) Write(p []byte) (int, error) {
	print := true
	if len(p) > 4 && string(p)[:5] == "debug" {
		print = printDebug
	}
	if print {
		fmt.Fprint(os.Stderr, string(p))
	}
	return len(p), nil
}

func main() {
	logger := Logger{}
	log.SetOutput(logger)
	log.SetFlags(0)

	args := args.Parse()

	printDebug = args.Debug

	switch args.Cmd {
	case "doc":
		doc.Doc(args.DocArgs)
	case "gen":
		gen.Gen(args.GenArgs)
	case "vet":
		vet.Vet(args.VetArgs)
		if rprt.ViolationCount() > 0 {
			os.Exit(1)
		}
	}
}
