package args

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// NoConfig returns true if '-no-config' flag is set.
func NoConfig() bool {
	if len(os.Args) <= 2 {
		return true
	} else {
		for _, a := range os.Args[2:] {
			if a == "-no-config" {
				return true
			}
		}
	}
	return false
}

func Parse() Args {
	fileCfg := FileCfg{}
	if !NoConfig() {
		thdlYml, err := os.ReadFile(".thdl.yml")
		if err == nil {
			err := yaml.UnmarshalStrict(thdlYml, &fileCfg)
			if err != nil {
				log.Fatalf("unmarshalling '.thdl.yml': %v", err)
			}
		}
	}
	fileCfg.propagateGlobalIgnore()
	log.Printf("debug: %s", fileCfg)

	args := Args{}
	setFileCfgArgs(fileCfg, &args)

	argsLen := len(os.Args)
	if argsLen == 1 {
		printHelp()
		os.Exit(1)
	}

	args.Cmd = os.Args[1]

	switch args.Cmd {
	case "ver":
		fmt.Printf("thdl version %s\n", Version)
		os.Exit(0)
	case "help":
		if argsLen < 3 {
			printHelp()
		} else if !isValidCommand(os.Args[2]) {
			printHelp()
			os.Exit(1)
		} else if os.Args[2] == "doc" {
			fmt.Printf(docHelpMsg)
		} else if os.Args[2] == "help" {
			fmt.Printf(helpHelpMsg)
		} else if os.Args[2] == "ver" {
			fmt.Printf(verHelpMsg)
		} else if os.Args[2] == "vet" {
			fmt.Printf(vetHelpMsg)
		}
		os.Exit(0)
	case "vet":
	case "doc":
		if argsLen < 3 {
			fmt.Printf("missing symbol path\n")
			os.Exit(1)
		}

		for _, arg := range os.Args[2 : argsLen-1] {
			switch arg {
			case "-debug":
				args.Debug = true
			case "-fusesoc":
				args.DocArgs.Fusesoc = true
			case "-no-bold":
				args.DocArgs.NoBold = true
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
		args.DocArgs.SymbolPath = symbolPath
	default:
		printHelp()
		os.Exit(1)
	}

	return args
}
