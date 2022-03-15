package args

import (
	"fmt"
)

var helpMsg string = `Hdl is a tool for easing the work with hardware description languages.
Version: %s

Usage:
  thdl <command> [arguments]

The commands are:
  check     check for extremely dumb mistakes
  doc       show or generate documentation
  generate  generate HDL files by processing sources
  help      print more information about a specific command
  version   print thdl version
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
}
