// Custom package for command line arguments parsing.
package args

import (
	"fmt"
	"os"
)

const Version string = "0.0.0"

func printVersion() {
	fmt.Println(Version)
	os.Exit(0)
}

func Parse() {
	argsLen := len(os.Args)
	if argsLen == 1 {
		printHelp()
	}
}
