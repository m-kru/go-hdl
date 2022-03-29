package args

import (
	"fmt"
)

var helpMsg string = `Thdl is a tool for easing the work with hardware description languages.
Version: %s

Usage:
  thdl <command> [arguments]

The commands are:
  doc   Show or generate documentation.
  gen   Generate HDL files by processing sources.
  help  Print more information about a specific command.
  ver   Print thdl version.
  vet   Check for extremely dumb mistakes.
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
}
