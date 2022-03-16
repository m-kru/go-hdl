package main

import (
	"os"

	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/check"
	"github.com/m-kru/go-thdl/internal/check/rprt"
)

func main() {
	cmdLineArgs := args.Parse()

	if cmdLineArgs["command"] == "check" {
		check.Check(cmdLineArgs)
		if rprt.ViolationCount() > 0 {
			os.Exit(1)
		}
	}
}
