package main

import (
	_ "fmt"

	"github.com/m-kru/go-hdl/internal/args"
	"github.com/m-kru/go-hdl/internal/check"
)

func main() {
	cmdLineArgs := args.Parse()

	if cmdLineArgs["command"] == "check" {
		check.Check(cmdLineArgs)
	}
}
