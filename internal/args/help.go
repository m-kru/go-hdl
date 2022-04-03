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
  vet   Check for extremely dumb mistakes.

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

  # Gen command settings.
  gen:
    ignore:
      - baz.vhd

  # Vet command settings.
  vet:
    ignore:
      - some/ignored/dir

Library documentation
---------------------
To document a library provide 'doc.<langugeExtension>' file within the library.
For example, to document a VHDL library provide 'doc.vhd' file. Each library
can have only one doc file. If more than one doc file is found per library,
then the error is reported.
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
}
