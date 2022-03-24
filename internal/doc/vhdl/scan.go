package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"github.com/m-kru/go-thdl/internal/utils"
)

func ScanFiles(filepaths []string, wg *sync.WaitGroup) {
	var filesWg sync.WaitGroup

	for _, fp := range filepaths {
		filesWg.Add(1)
		go scanFile(fp, &filesWg)
	}

	filesWg.Wait()
	wg.Done()
}

var commentLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*--`)
var emptyLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*$`)
var constantDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*constant\b`)
var entityDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*entity\s+(\w*)\s+is`)
var enumTypeDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*type\s+(\w+)\s+is\s*\(`)
var packageDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*package\s+(\w*)\s+is`)
var endRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\b`)
var endPackageRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\s+package\b`)
var endsWithSemicolonRegExp *regexp.Regexp = regexp.MustCompile(`;\s*$`)

type scanContext struct {
	scanner *bufio.Scanner
	line    []byte

	startIdx uint32 // Line start index.
	endIdx   uint32 // Line end index.

	docPresent bool
	docStart   uint32
	docEnd     uint32
}

func (sc *scanContext) proceed() bool {
GETLINE:
	if ok := sc.scanner.Scan(); !ok {
		return false
	}

	sc.line = bytes.ToLower(sc.scanner.Bytes())

	sc.startIdx = sc.endIdx
	sc.endIdx += uint32(len(sc.line)) + 1

	if len(emptyLineRegExp.FindIndex(sc.line)) > 0 {
		sc.docPresent = false
		goto GETLINE
	} else if len(commentLineRegExp.FindIndex(sc.line)) > 0 {
		sc.docEnd = sc.endIdx
		if !sc.docPresent {
			sc.docStart = sc.startIdx
			sc.docPresent = true
		}
	}

	return true
}

func scanFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if utils.IsIgnoredVHDLFile(filepath) {
		return
	}

	libName := "_unknown_"
	if !libContainer.Has(libName) {
		l := lib.MakeLibrary(libName)
		libContainer.Add(&l)
	}

	f, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", filepath, err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(f))

	sCtx := scanContext{scanner: scanner}

	for sCtx.proceed() {
		if submatches := entityDeclarationRegExp.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
			name := string(sCtx.line[submatches[2]:submatches[3]])
			ent, err := scanEntityDeclaration(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer[libName].AddSymbol(ent)
		} else if submatches := packageDeclarationRegExp.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
			name := string(sCtx.line[submatches[2]:submatches[3]])
			pkg, err := scanPackageDeclaration(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer[libName].AddSymbol(pkg)
		}
	}
}

func scanEntityDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	ent := Entity{
		Symbol{
			filepath:  filepath,
			name:      name,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		ent.hasDoc = true
		ent.docStart = sc.docStart
		ent.docEnd = sc.docEnd
	}

	for sc.proceed() {
		if len(endRegExp.FindIndex(sc.line)) > 0 {
			ent.codeEnd = sc.endIdx
			return ent, nil
		}
	}

	return ent, fmt.Errorf("entity declaration end line not found")
}

func scanPackageDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	pkg := Package{
		Symbol: Symbol{
			filepath:  filepath,
			name:      name,
			codeStart: sc.startIdx,
		},
		symbols: map[string]symbol.Symbol{},
	}

	if sc.docPresent {
		pkg.hasDoc = true
		pkg.docStart = sc.docStart
		pkg.docEnd = sc.docEnd
	}

	for sc.proceed() {
		/*
			if idxs := constantDeclarationRegExp.FindIndex(sc.line); len(idxs) > 0 {
				consts, err := scanConstantDeclaration(filepath, idxs[1], sc)
				if err != nil {
					return pkg, fmt.Errorf("package '%s': %v", name, err)
				}
				for _, c := range consts {
					err  = pkg.AddSymbol(c)
					if err != nil {
						return pkg, fmt.Errorf("package '%s': %v", name, err)
					}
				}
			}
		*/
		if submatches := enumTypeDeclarationRegExp.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanEnumTypeDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(t)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if (len(endRegExp.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(name))) ||
			(len(endPackageRegExp.FindIndex(sc.line)) > 0) {
			pkg.codeEnd = sc.endIdx
			return pkg, nil
		}
	}

	return pkg, fmt.Errorf("package declaration end line not found")
}

func scanEnumTypeDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	t := Type{
		Symbol{
			filepath:  filepath,
			name:      name,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		t.hasDoc = true
		t.docStart = sc.docStart
		t.docEnd = sc.docEnd
	}

	for {
		if len(endsWithSemicolonRegExp.FindIndex(sc.line)) > 0 {
			t.codeEnd = sc.endIdx
			return t, nil
		}

		sc.proceed()
	}

	return t, fmt.Errorf("enum type declaration line with ';' not found")
}

// endIdx is the index of 'constant' keyword end.
func scanConstantDeclaration(filepath string, endidx int, sc *scanContext) ([]symbol.Symbol, error) {
	const_ := Constant{
		Symbol{
			filepath:  filepath,
			codeStart: sc.startIdx,
		},
	}

	//names := []string
	syms := []symbol.Symbol{}

	if len(endsWithSemicolonRegExp.FindIndex(sc.line)) > 0 {
		const_.codeEnd = sc.endIdx
		return syms, nil
	}

	return syms, nil
}
