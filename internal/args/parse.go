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

	if len(os.Args) == 1 {
		log.Fatalf("missing command, run 'thdl help' for more information")
	}

	args.Cmd = os.Args[1]

	switch args.Cmd {
	case "doc":
		parseDocArgs(&args)
	case "gen":
		parseGenArgs(&args)
	case "help":
		if len(os.Args) < 3 {
			printHelp()
		} else if !isValidCommand(os.Args[2]) {
			printHelp()
			os.Exit(1)
		} else if os.Args[2] == "doc" {
			fmt.Printf(docHelpMsg)
		} else if os.Args[2] == "gen" {
			fmt.Printf(genHelpMsg)
		} else if os.Args[2] == "help" {
			fmt.Printf(helpHelpMsg)
		} else if os.Args[2] == "ver" {
			fmt.Printf(verHelpMsg)
		} else if os.Args[2] == "vet" {
			fmt.Printf(vetHelpMsg)
		}
		os.Exit(0)
	case "ver":
		fmt.Printf("thdl version %s\n", Version)
		os.Exit(0)
	case "vet":
		parseVetArgs(&args)
	default:
		log.Fatalf(fmt.Sprintf("invalid command '%s', run 'thdl help' for more information", args.Cmd))
	}

	return args
}

func parseDocArgs(args *Args) {
	var param string
	var expectArg bool

	for i, a := range os.Args[2:] {
		if expectArg {
			switch param {
			case "-html-copyright":
				args.DocArgs.HTML.Copyright = a
			case "-html-path":
				args.DocArgs.HTML.Path = a
			case "-html-title":
				args.DocArgs.HTML.Title = a
			}
			expectArg = false
			continue
		}
		switch a {
		case "-debug":
			args.Debug = true
		case "-fusesoc":
			args.DocArgs.Fusesoc = true
		case "-no-bold":
			args.DocArgs.NoBold = true
		// HTML arguments.
		case "-html":
			args.DocArgs.GenHTML = true
		case "-html-copyright", "-html-path", "-html-title":
			param = a
			expectArg = true
		default:
			if i == len(os.Args)-3 {
				args.DocArgs.SymbolPath = a
			} else {
				log.Fatalf("invalid doc command flag '%s'\n", a)
			}
		}
	}

	if expectArg {
		log.Fatalf("missing argument for parameter '%s'", param)
	}

	// HTML path post-processing.
	if args.DocArgs.GenHTML {
		if args.DocArgs.HTML.Path == "" {
			args.DocArgs.HTML.Path = "./"
		} else {
			args.DocArgs.HTML.Path += "/"
		}
	} else {
		if args.DocArgs.SymbolPath == "" {
			log.Fatalf("missing symbol path\n")
		}
	}
}

func parseGenArgs(args *Args) {
	for i, a := range os.Args[2:] {
		switch a {
		case "-to-stdout":
			args.GenArgs.ToStdout = true
		default:
			if i == len(os.Args)-3 {
				args.GenArgs.Filepath = a
			} else {
				log.Fatalf("invalid gen command flag '%s'\n", a)
			}
		}
	}
}

func parseVetArgs(args *Args) {
	for i, a := range os.Args[2:] {
		switch a {
		default:
			if i == len(os.Args)-3 {
				args.VetArgs.Filepath = a
			} else {
				log.Fatalf("invalid vet command flag '%s'\n", a)
			}
		}
	}
}
