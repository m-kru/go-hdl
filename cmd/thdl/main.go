package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/check"
	"github.com/m-kru/go-thdl/internal/check/rprt"
	"github.com/m-kru/go-thdl/internal/doc"
)

var printDebug bool = false

type Logger struct{}

func (l Logger) Write(p []byte) (int, error) {
	print := true
	if len(p) > 4 && string(p)[:5] == "debug" {
		print = printDebug
	}
	if print {
		fmt.Fprintf(os.Stderr, string(p))
	}
	return len(p), nil
}

func main() {
	logger := Logger{}
	log.SetOutput(logger)
	log.SetFlags(0)

	args := args.Parse()

	printDebug = args.Debug

	if args.Cmd == "check" {
		check.Check(args.CheckArgs)
		if rprt.ViolationCount() > 0 {
			os.Exit(1)
		}
	} else if args.Cmd == "doc" {
		if doc.Doc(args.DocArgs) > 0 {
			os.Exit(1)
		}
	}
}
