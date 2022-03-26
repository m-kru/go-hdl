package doc

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"sync"
)

var docArgs args.DocArgs

func Doc(args args.DocArgs) uint8 {
	docArgs = args

	ScanFiles()

	symbolPaths := resolveSymbolPath(args.SymbolPath)
	foundSymbols := map[symbolPath]symbol.Symbol{}

	for _, sp := range symbolPaths {
		paths, syms := findSymbol(sp)
		for i, _ := range paths {
			foundSymbols[paths[i]] = syms[i]
		}
	}

	foundCount := len(foundSymbols)

	if foundCount == 0 {
		log.Fatalf("no symbol matching path '%s' found", args.SymbolPath)
	} else if foundCount == 1 {
		for path, sym := range foundSymbols {
			fmt.Printf("%s\n\n", path)
			fmt.Printf("%s\n\n", sym.Filepath())
			doc, code := sym.DocCode()
			fmt.Printf(utils.Deindent(doc))
			if !args.NoBold {
				code = utils.BoldCodeTerminal(path.language, code)
			}
			fmt.Printf(utils.Deindent(code))
		}
	} else {
		msg := "provided path is ambiguous, found symbols with following paths:"
		for path, _ := range foundSymbols {
			msg = fmt.Sprintf("%s\n  %s", msg, path)
		}
		log.Fatalf("%s", msg)
	}

	return 0
}

func ScanFiles() {
	var wg sync.WaitGroup
	defer wg.Wait()

	vhdlFiles, err := utils.GetFilePathsByExtension(".vhd", ".")
	if err != nil {
		log.Fatalf("%v", err)
	}
	vhdlFiles = docArgs.FilterIgnored(vhdlFiles)
	wg.Add(1)
	vhdl.ScanFiles(docArgs, vhdlFiles, &wg)
}

func findSymbol(sp symbolPath) (paths []symbolPath, syms []symbol.Symbol) {
	var ok bool

	libNames := []string{}

	if sp.library != "" {
		libNames = append(libNames, sp.library)
	} else {
		switch sp.language {
		case "vhdl":
			libNames = vhdl.LibraryNames()
		default:
			panic("should never happen")
		}
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

		if tmpSp.primary == "" && tmpSp.secondary == "" {
			panic("should never happen")
		} else if tmpSp.primary == "" {
			for _, primaryName := range l.SymbolNames() {
				tmpSp.primary = primaryName
				pri, ok := l.GetSymbol(tmpSp.primary)
				if !ok {
					continue
				}
				sec, ok := pri.GetSymbol(tmpSp.secondary)
				if !ok {
					continue
				}
				paths = append(paths, tmpSp)
				syms = append(syms, sec)
			}
		} else {
			pri, ok := l.GetSymbol(tmpSp.primary)
			if !ok {
				continue
			}
			if tmpSp.secondary == "" {
				paths = append(paths, tmpSp)
				syms = append(syms, pri)
			} else {
				sec, ok := pri.GetSymbol(tmpSp.secondary)
				if !ok {
					continue
				}
				paths = append(paths, tmpSp)
				syms = append(syms, sec)
			}
		}
	}

	return
}
