package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

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
var endRegExp *regexp.Regexp = regexp.MustCompile(`\bend\b`)
var entityDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*entity\s+(\w*)\s+is`)

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

	sc.startIdx = sc.endIdx + 1
	sc.endIdx += uint32(len(sc.line))

	return true
}

func scanFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if utils.IsIgnoredVHDLFile(filepath) {
		return
	}

	lib := "_unknown_"

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
			if sCtx.docPresent {
				sCtx.docEnd = sCtx.endIdx
			} else {
				sCtx.docStart = sCtx.startIdx
				sCtx.docPresent = true
			}
		} else if submatches := entityDeclarationRegExp.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
			ent, err := scanEntityDeclaration(
				&filepath, string(sCtx.line[submatches[2]:submatches[3]]), &sCtx,
			)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			if !libContainer.Has(lib) {
				libContainer.Add(Library{name: lib})
			}
			libContainer[lib].AddSymbol(ent)
		}
	}
}

func scanEntityDeclaration(filepath *string, name string, sc *scanContext) (symbol.Symbol, error) {
	ent := Entity{
		filepath:  filepath,
		name:      name,
		codeStart: sc.startIdx,
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
