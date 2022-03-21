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

type CheckArgs struct {
	NoConfig bool
}

type DocArgs struct {
	Fusesoc    bool
	NoBold     bool
	NoConfig   bool
	SymbolPath string
}

type Args struct {
	Cmd       string
	CheckArgs CheckArgs
	DocArgs   DocArgs
}

func Parse() Args {
	args := Args{}

	argsLen := len(os.Args)
	if argsLen == 1 {
		printHelp()
		os.Exit(1)
	}

	args.Cmd = os.Args[1]

	switch args.Cmd {
	case "version":
		fmt.Printf("thdl version %s\n", Version)
		os.Exit(0)
	case "help":
		if argsLen < 3 {
			printHelp()
		} else if !isValidCommand(os.Args[2]) {
			printHelp()
			os.Exit(1)
		} else if os.Args[2] == "check" {
			fmt.Printf(checkHelpMsg)
		} else if os.Args[2] == "doc" {
			fmt.Printf(docHelpMsg)
		}
		os.Exit(0)
	case "check":
	case "doc":
		docArgs := DocArgs{}
		if argsLen < 3 {
			fmt.Printf("missing symbol path\n")
			os.Exit(1)
		}

		for _, arg := range os.Args[2 : argsLen-1] {
			switch arg {
			case "--fusesoc":
				docArgs.Fusesoc = true
			case "--no-bold":
				docArgs.NoBold = true
			default:
				fmt.Printf("invalid doc command flag '%s'\n", arg)
				os.Exit(1)
			}
		}

		// Path to symbol is always the last argument.
		symbolPath := os.Args[argsLen-1]
		if symbolPath[0] == '-' {
			if isValidDocFlag(symbolPath) {
				fmt.Printf("missing symbol path\n")
			} else {
				fmt.Printf("invalid symbol path '%s'\n", symbolPath)
			}
			os.Exit(1)
		}
		docArgs.SymbolPath = symbolPath

		args.DocArgs = docArgs
	default:
		printHelp()
		os.Exit(1)
	}

	return args
}
