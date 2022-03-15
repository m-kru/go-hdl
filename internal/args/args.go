// Custom package for command line arguments parsing.
package args

import (
	"fmt"
	"os"

	"github.com/m-kru/go-thdl/internal/utils"
)

const Version string = "0.0.0"

func printVersion() {
	fmt.Println(Version)
	os.Exit(0)
}

func Parse() map[string]string {
	args := map[string]string{}

	argsLen := len(os.Args)
	if argsLen == 1 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	if !utils.IsValidCommand(cmd) {
		printHelp()
		os.Exit(1)
	}
	args["command"] = os.Args[1]

	if cmd == "version" {
		fmt.Printf("thdl version %s\n", Version)
		os.Exit(0)
	} else if cmd == "help" {
		if argsLen < 3 {
			printHelp()
		} else if !utils.IsValidCommand(os.Args[2]) {
			printHelp()
			os.Exit(1)
		} else if os.Args[2] == "check" {
			fmt.Printf(checkHelpMsg)
		}
		os.Exit(0)
	}

	return args
}
