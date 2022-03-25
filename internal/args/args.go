// Custom package for command line arguments parsing.
package args

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

const Version string = "0.0.0"

func printVersion() {
	fmt.Println(Version)
	os.Exit(0)
}

type IgnoreList struct {
	ignore []string
}

func (il IgnoreList) FilterIgnored(filepaths []string) []string {
	ret := []string{}

	for _, fp := range filepaths {
		ignore := false
		for _, i := range il.ignore {
			if strings.Contains(fp, i) {
				ignore = true
				log.Printf("debug: ignoring %s\n", fp)
				break
			}
		}
		if !ignore {
			ret = append(ret, fp)
		}
	}

	return ret
}

type CheckArgs struct {
	IgnoreList
}

type DocArgs struct {
	IgnoreList
	Fusesoc    bool
	NoBold     bool
	SymbolPath string
}

type GenArgs struct {
	IgnoreList
	Fusesoc    bool
	NoBold     bool
	SymbolPath string
}

type Args struct {
	Cmd       string
	Debug     bool
	CheckArgs CheckArgs
	DocArgs   DocArgs
	GenArgs   GenArgs
}

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

func setFileCfgArgs(fc FileCfg, args *Args) {
	args.CheckArgs.IgnoreList.ignore = fc.Check.Ignore

	args.DocArgs.IgnoreList.ignore = fc.Doc.Ignore
	args.DocArgs.Fusesoc = fc.Doc.Fusesoc
	args.DocArgs.NoBold = fc.Doc.NoBold

	args.GenArgs.IgnoreList.ignore = fc.Gen.Ignore
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
	fmt.Printf("debug: %s", fileCfg)

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
		} else if os.Args[2] == "check" {
			fmt.Printf(checkHelpMsg)
		} else if os.Args[2] == "doc" {
			fmt.Printf(docHelpMsg)
		}
		os.Exit(0)
	case "check":
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
