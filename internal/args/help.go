package args

import (
	"fmt"
	"os"
)

var helpMsg string = `Hdl is a tool for easing the work with hardware description languages.
Version: %s

Usage:
  hdl <command> [arguments]

The commands are:
  check     check for extremely dumb mistakes
  doc       show or generate documentation
  generate  generate HDL files by processing sources

Use "hdl help <command>" for more information about a command.
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
	os.Exit(0)
}
