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

	args := args.Parse()

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
