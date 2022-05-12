package args

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// isPresent returns true if given argument is present in the argument list.
func isPresent(arg string) bool {
	if len(os.Args) <= 2 {
		return true
	} else {
		for _, a := range os.Args[2:] {
			if a == arg {
				return true
			}
		}
	}
	return false
}

func Parse() Args {
	fileCfg := FileCfg{}
	if !isPresent("-no-config") {
		thdlYml, err := os.ReadFile(".thdl.yml")
		if err == nil {
			err := yaml.UnmarshalStrict(thdlYml, &fileCfg)
			if err != nil {
				log.Fatalf("unmarshalling '.thdl.yml': %v", err)
			}
		}
	}
	fileCfg.propagateGlobalIgnore()

	args := Args{}
	setFileCfgArgs(fileCfg, &args)

	argsLen := len(os.Args)
	if argsLen == 1 {
		log.Fatalf("missing command, run 'thdl help' for more information")
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

		iterLen := argsLen - 1
		if isPresent("-html") {
			iterLen = argsLen
			args.DocArgs.GenHTML = true
		}
		for i := 2; i < iterLen; i++ {
			arg := os.Args[i]

			switch arg {
			case "-debug":
				args.Debug = true
			case "-fusesoc":
				args.DocArgs.Fusesoc = true
			case "-no-bold":
				args.DocArgs.NoBold = true
			// HTML arguments.
			case "-html":
			case "-html-copyright":
				args.DocArgs.HTML.Copyright = os.Args[i+1]
				i++
			case "-html-title":
				args.DocArgs.HTML.Title = os.Args[i+1]
				i++
			case "-html-path":
				args.DocArgs.HTML.Path = os.Args[i+1]
				i++
			default:
				fmt.Printf("invalid doc command flag '%s'\n", arg)
				os.Exit(1)
			}
		}

		// Path to symbol is always the last argument, if there is no '-html' flag.
		if args.DocArgs.GenHTML {
			if args.DocArgs.HTML.Path == "" {
				args.DocArgs.HTML.Path = "./"
			} else {
				args.DocArgs.HTML.Path += "/"
			}
		} else {
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
		}
	default:
		log.Fatalf(fmt.Sprintf("invalid command '%s', run 'thdl help' for more information", args.Cmd))
	}

	return args
}
