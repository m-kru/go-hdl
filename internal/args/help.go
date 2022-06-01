package args

import (
	"fmt"
)

var helpMsg string = `Thdl is a tool for easing the work with hardware description languages.
Version: %s

Usage:
------

  thdl <command> [arguments]

The commands are:
  doc   Show or generate documentation.
  gen   Generate HDL files by processing sources.
  help  Print more information about a specific command.
  ver   Print thdl version.
  vet   Check for likely mistakes.


Configuration
-------------

Thdl behavior can be configured using '.thdl.yml' file. This file must be placed
in the working directory (usually project's root directory) to be read.
Below snippet presents example configuration. It shows all currently supported
settings.

  # Libraries.
  # Key is the name of the library.
  # Value is a list of path patterns (strings).
  # If a file path contains any pattern from the list, then symbols
  # from this file will be put into this particular library.
  libs:
    my_lib:
      - gw/my-lib

  # Global ignore path patterns.
  ignore:
    - ignore/this/dir
    - foo.vhd

  # Doc command settings.
  doc:
    ignore: [bar.vhd]
    fusesoc: true
    no-bold: true
    html:
      # Copyright string placed in the left bottom corner.
      # If not set, then there is no copyright footer.
      copyright: "Copyright string"
      path: output/path/for/the/generated/html/documentation
      title: "HTML title string"

  # Gen command settings.
  gen:
    ignore:
      - baz.vhd

  # Vet command settings.
  vet:
    ignore:
      - some/ignored/dir
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
}
