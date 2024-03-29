package doc

import (
	"fmt"
	"github.com/m-kru/go-hdl/internal/args"
	"github.com/m-kru/go-hdl/internal/doc/html"
	"github.com/m-kru/go-hdl/internal/doc/lib"
	"github.com/m-kru/go-hdl/internal/doc/sym"
	"github.com/m-kru/go-hdl/internal/doc/vhdl"
	"github.com/m-kru/go-hdl/internal/utils"
	"log"
	"sync"
)

var docArgs args.DocArgs

func Doc(args args.DocArgs) {
	docArgs = args

	ScanFiles()

	if docArgs.GenHTML {
		html.Generate(docArgs.HTML)
		return
	}

	symbolPaths := resolveSymbolPath(args.SymbolPath)
	log.Printf("debug: looking for following symbol paths:\n")
	for _, p := range symbolPaths {
		log.Printf("debug:   %s", p.DebugString())
	}

	foundSymbols := map[symbolPath][]sym.Symbol{}

	for _, sp := range symbolPaths {
		paths, syms := findSymbol(sp)
		for i := range paths {
			foundSymbols[paths[i]] = syms[i]
		}
	}

	foundCount := len(foundSymbols)

	if foundCount == 0 {
		log.Fatalf("found no symbol matching path '%s'", args.SymbolPath)
	} else if foundCount == 1 {
		for path, syms := range foundSymbols {
			fmt.Printf("%s\n", path)
			prevFilepath := ""
			sym.SortByLineNum(syms)
			for _, s := range syms {
				fp := s.Filepath()
				if fp != "" && fp != prevFilepath {
					fmt.Printf("\n%s\n", fp)
					prevFilepath = fp
				}
				fmt.Print("\n")
				doc, code := s.DocCode()
				fmt.Print(utils.Deindent(doc))
				if !args.NoBold {
					code = utils.BoldCodeTerminal(path.language, code)
				}
				if _, ok := s.(vhdl.Package); ok {
					if doc != "" {
						fmt.Println()
					}
				}
				fmt.Print(utils.Deindent(code))
			}
		}
	} else {
		msg := "provided path is ambiguous, found symbols with following paths:"
		for path := range foundSymbols {
			msg = fmt.Sprintf("%s\n  %s", msg, path)
		}
		log.Fatalf("%s", msg)
	}
}

func ScanFiles() {
	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles := utils.GetVHDLFilePaths()
	vhdlFiles = docArgs.FilterIgnored(vhdlFiles)

	wg.Add(1)
	vhdl.ScanFiles(docArgs, vhdlFiles, &wg)
}

func findSymbol(sp symbolPath) (paths []symbolPath, syms [][]sym.Symbol) {
	var ok bool

	if sp.isLibrary() {
		return findLibrary(sp)
	}

	libNames := []string{}
	if sp.library == "*" {
		switch sp.language {
		case "vhdl":
			libNames = vhdl.LibraryNames()
		default:
			panic("should never happen")
		}
	} else {
		libNames = append(libNames, sp.library)
	}

	for _, libName := range libNames {
		tmpSp := sp
		tmpSp.library = libName

		var l *lib.Library

		switch sp.language {
		case "vhdl":
			l, ok = vhdl.GetLibrary(tmpSp.library)
		default:
			panic("should never happen")
		}
		if !ok {
			continue
		}

		priNames := []string{}
		if tmpSp.primary == "*" {
			priNames = l.InnerKeys()
		} else {
			priNames = append(priNames, tmpSp.primary)
		}

		for _, priName := range priNames {
			tmpSp.primary = priName
			tmpSp.secondary = ""
			tmpSp.tertiary = ""
			pri := l.GetSymbol(priName)
			if len(pri) == 0 {
				continue
			}

			secNames := []string{}

			if sp.secondary == "" {
				paths = append(paths, tmpSp)
				syms = append(syms, pri)
				continue
			} else if sp.secondary == "*" {
				secNames = pri[0].InnerKeys()
			} else {
				secNames = append(secNames, sp.secondary)
			}

			for _, secName := range secNames {
				tmpSp.secondary = secName
				tmpSp.tertiary = ""
				sec := pri[0].GetSymbol(secName)
				if len(sec) == 0 {
					continue
				}

				terName := sp.tertiary

				if terName == "" {
					paths = append(paths, tmpSp)
					syms = append(syms, sec)
					continue
				} else if terName == "*" {
					log.Fatalf("tertiary sym can't be '*' wildcard")
				}

				ter := sec[0].GetSymbol(terName)
				if len(ter) == 0 {
					continue
				}
				tmpSp.tertiary = terName
				paths = append(paths, tmpSp)
				syms = append(syms, ter)
			}
		}
	}

	return
}

func findLibrary(sp symbolPath) (paths []symbolPath, syms [][]sym.Symbol) {
	var ok bool
	var l *lib.Library

	switch sp.language {
	case "vhdl":
		l, ok = vhdl.GetLibrary(sp.library)
	default:
		panic("should never happen")
	}
	if !ok {
		return
	}

	paths = append(paths, sp)
	syms = append(syms, []sym.Symbol{l})

	return
}
