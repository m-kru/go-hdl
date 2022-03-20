package main

import (
	"log"
	"os"

	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/check"
	"github.com/m-kru/go-thdl/internal/check/rprt"
	"github.com/m-kru/go-thdl/internal/doc"
)

func main() {
	log.SetFlags(0)

	cmdLineArgs := args.Parse()

	if cmdLineArgs["command"] == "check" {
		check.Check(cmdLineArgs)
		if rprt.ViolationCount() > 0 {
			os.Exit(1)
		}
	} else if cmdLineArgs["command"] == "doc" {
		if doc.Doc(cmdLineArgs) > 0 {
			os.Exit(1)
		}
	}
}
