package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/m-kru/go-thdl/internal/args"
	"github.com/m-kru/go-thdl/internal/doc/lib"
	"github.com/m-kru/go-thdl/internal/doc/symbol"
	"github.com/m-kru/go-thdl/internal/utils"
)

var docArgs args.DocArgs

func ScanFiles(args args.DocArgs, filepaths []string, wg *sync.WaitGroup) {
	docArgs = args

	var filesWg sync.WaitGroup

	for _, fp := range filepaths {
		filesWg.Add(1)
		go scanFile(fp, &filesWg)
	}

	filesWg.Wait()
	wg.Done()
}

var emptyLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*$`)
var commentLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*--`)

var entityDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*entity\s+(\w*)\s+is`)

var packageDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*package\s+(\w*)\s+is`)
var endPackageRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\s+package\b`)
var arrayTypeDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*type\s+(\w+)\s+is\s+array\b`)
var enumTypeDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*type\s+(\w+)\s+is\s*\(`)
var recordTypeDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*type\s+(\w+)\s+is\s+record\b`)
var endRecordRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\s+record\b`)
var constantDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*constant\s+(\w+)\s*(,\s*\w+)?\s*(,\s*\w+)?`)

var endRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\b`)
var endWithSemicolonRegExp *regexp.Regexp = regexp.MustCompile(`^\s*end\s*;`)
var endsWithSemicolonRegExp *regexp.Regexp = regexp.MustCompile(`;\s*$`)

var packageBodyDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*package\s+body\s+\w+\s+is\b`)
var architectureDeclarationRegExp *regexp.Regexp = regexp.MustCompile(`^\s*architecture\s+\w+\s+of\s*\w+\s+is\b`)

type scanContext struct {
	scanner *bufio.Scanner
	line    []byte

	startIdx uint32 // Line start index.
	endIdx   uint32 // Line end index.

	docPresent bool
	docStart   uint32
	docEnd     uint32
}

// proceed returns false on EOF, architecture declaration or package
// body declaration. There is no point in scanning architecture
// declarations and package bodies, as they either contain private symbols
// or they implement symbols declared in the package declaration.
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
	} else if len(packageBodyDeclarationRegExp.FindIndex(sc.line)) > 0 ||
		len(architectureDeclarationRegExp.FindIndex(sc.line)) > 0 {
		return false
	}

	return true
}

func scanFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if utils.IsIgnoredVHDLFile(filepath) {
		return
	}

	libName := docArgs.Lib(filepath)
	if libName == "" {
		libName = "work"
	}
	l := lib.MakeLibrary(libName)
	libContainer.Add(&l)

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
			libContainer.AddSymbol(libName, ent)
		} else if submatches := packageDeclarationRegExp.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
			name := string(sCtx.line[submatches[2]:submatches[3]])
			pkg, err := scanPackageDeclaration(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, pkg)
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
		Consts: map[string]symbol.Symbol{},
		Funcs:  map[string]symbol.Symbol{},
		Procs:  map[string]symbol.Symbol{},
		Types:  map[string]symbol.Symbol{},
	}

	if sc.docPresent {
		pkg.hasDoc = true
		pkg.docStart = sc.docStart
		pkg.docEnd = sc.docEnd
	}

	for sc.proceed() {
		if submatches := constantDeclarationRegExp.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			names := []string{}
			for i := 1; i < len(submatches)/2; i++ {
				if submatches[2*i] < 0 {
					continue
				}
				name := string(sc.line[submatches[2*i]:submatches[2*i+1]])
				if name[0] == ',' {
					name = strings.TrimSpace(name[1:])
				}
				names = append(names, name)
			}
			consts, err := scanConstantDeclaration(filepath, names, sc)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
			for _, c := range consts {
				err = pkg.AddSymbol(c)
				if err != nil {
					return pkg, fmt.Errorf("package '%s': %v", name, err)
				}
			}
		} else if submatches := arrayTypeDeclarationRegExp.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanArrayTypeDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(t)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if submatches := enumTypeDeclarationRegExp.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanEnumTypeDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(t)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if submatches := recordTypeDeclarationRegExp.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanRecordTypeDeclaration(filepath, name, sc)
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

		if !sc.proceed() {
			break
		}
	}

	return t, fmt.Errorf("enum declaration line with ';' not found")
}

func scanArrayTypeDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
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

		if !sc.proceed() {
			break
		}
	}

	return t, fmt.Errorf("array declaration end line not found")
}

func scanRecordTypeDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
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
		if (len(endRegExp.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(name))) ||
			(len(endRecordRegExp.FindIndex(sc.line)) > 0) ||
			(len(endWithSemicolonRegExp.FindIndex(sc.line)) > 0) {
			t.codeEnd = sc.endIdx
			return t, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return t, fmt.Errorf("record declaration line with ';' not found")
}

func scanConstantDeclaration(filepath string, names []string, sc *scanContext) ([]symbol.Symbol, error) {
	const_ := Constant{
		Symbol{
			filepath:  filepath,
			codeStart: sc.startIdx,
		},
	}

	syms := []symbol.Symbol{}

	if len(endsWithSemicolonRegExp.FindIndex(sc.line)) > 0 {
		const_.codeEnd = sc.endIdx

		for _, n := range names {
			const_.name = n
			syms = append(syms, const_)
		}

		return syms, nil
	}

	return syms, fmt.Errorf("constant declaration line with ';' not found")
}
