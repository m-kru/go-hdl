package doc

import (
	"fmt"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"github.com/m-kru/go-thdl/internal/doc/vhdl"
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"sync"
)

func Doc(cmdLineArgs map[string]string) uint8 {
	ScanFiles()

	symbolPaths := resolveSymbolPath(cmdLineArgs["symbolPath"])
	foundSymbols := map[symbolPath]symbol.Symbol{}

	for _, sp := range symbolPaths {
		fmt.Printf("looking for symbol %v", sp)
		paths, syms := findSymbol(sp)
		for i, _ := range paths {
			foundSymbols[paths[i]] = syms[i]
		}
	}

	foundCount := len(foundSymbols)

	if foundCount == 0 {
		log.Fatalf("no symbol matching path '%s' found", cmdLineArgs["symbolPath"])
	} else if foundCount == 1 {
		for _, sym := range foundSymbols {
			fmt.Printf(sym.Doc())
		}
	} else {
		msg := "provided path is ambiguous, found symbols with following paths:"
		for path, _ := range foundSymbols {
			msg = fmt.Sprintf("%s%s", msg, path)
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
	wg.Add(1)
	vhdl.ScanFiles(vhdlFiles, &wg)
}

func findSymbol(sp symbolPath) (paths []symbolPath, syms []symbol.Symbol) {
	var lib symbol.Symbol
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

		switch sp.language {
		case "vhdl":
			lib, ok = vhdl.GetLibrary(tmpSp.library)
		default:
			panic("should never happen")
		}
		if !ok {
			continue
		}

		if tmpSp.primary == "" && tmpSp.secondary == "" {
			panic("should never happen")
		} else if tmpSp.primary == "" {
			for _, primaryName := range lib.SymbolNames() {
				tmpSp.primary = primaryName
				pri, ok := lib.GetSymbol(tmpSp.primary)
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
			pri, ok := lib.GetSymbol(tmpSp.primary)
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
