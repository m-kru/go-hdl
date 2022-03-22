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
var entityDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*entity\s+(\w*)\s+is`)
var packageDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*package\s+(\w*)\s+is`)
var endRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\b`)
var endPackageRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\s+package\b`)

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
	if ok := sc.scanner.Scan(); !ok {
		return false
	}

	sc.line = bytes.ToLower(sc.scanner.Bytes())

	sc.startIdx = sc.endIdx
	sc.endIdx += uint32(len(sc.line)) + 1

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
		if len(emptyLineRegExp.FindIndex(sCtx.line)) > 0 {
			sCtx.docPresent = false
		} else if len(commentLineRegExp.FindIndex(sCtx.line)) > 0 {
			sCtx.docEnd = sCtx.endIdx
			if !sCtx.docPresent {
				sCtx.docStart = sCtx.startIdx
				sCtx.docPresent = true
			}
		} else if submatches := entityDeclarationRegExp.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
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
		Symbol{
			filepath:  filepath,
			name:      name,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		pkg.hasDoc = true
		pkg.docStart = sc.docStart
		pkg.docEnd = sc.docEnd
	}

	for sc.proceed() {
		if (len(endRegExp.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(name))) ||
			(len(endPackageRegExp.FindIndex(sc.line)) > 0) {
			pkg.codeEnd = sc.endIdx
			return pkg, nil
		}
	}

	return pkg, fmt.Errorf("package declaration end line not found")
}
