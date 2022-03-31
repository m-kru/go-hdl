package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
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

type scanContext struct {
	scanner *bufio.Scanner
	line    []byte
	lineNum uint32

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

	sc.lineNum += 1
	sc.line = bytes.ToLower(sc.scanner.Bytes())

	sc.startIdx = sc.endIdx
	sc.endIdx += uint32(len(sc.line)) + 1

	if len(emptyLine.FindIndex(sc.line)) > 0 {
		sc.docPresent = false
		goto GETLINE
	} else if len(commentLine.FindIndex(sc.line)) > 0 {
		sc.docEnd = sc.endIdx
		if !sc.docPresent {
			sc.docStart = sc.startIdx
			sc.docPresent = true
		}
	} else if len(packageBodyDeclaration.FindIndex(sc.line)) > 0 ||
		len(architectureDeclaration.FindIndex(sc.line)) > 0 {
		sc.docPresent = false
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
		if submatches := entityDeclaration.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
			name := string(sCtx.line[submatches[2]:submatches[3]])
			ent, err := scanEntityDeclaration(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, ent)
		} else if submatches := packageInstantiation.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
			name := string(sCtx.line[submatches[2]:submatches[3]])
			pkgInst, err := scanPackageInstantiation(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, pkgInst)
		} else if submatches := packageDeclaration.FindSubmatchIndex(sCtx.line); len(submatches) > 0 {
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
	e := Entity{
		Symbol{
			filepath:  filepath,
			name:      name,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		e.docStart = sc.docStart
		e.docEnd = sc.docEnd
	}

	for sc.proceed() {
		if len(end.FindIndex(sc.line)) > 0 {
			e.codeEnd = sc.endIdx
			return e, nil
		}
	}

	return e, fmt.Errorf("entity declaration end line not found")
}

func scanPackageDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	pkg := Package{
		Symbol: Symbol{
			filepath:  filepath,
			name:      name,
			codeStart: sc.startIdx,
		},
		Consts: map[symbol.ID]symbol.Symbol{},
		Funcs:  map[symbol.ID]symbol.Symbol{},
		Procs:  map[symbol.ID]symbol.Symbol{},
		Types:  map[symbol.ID]symbol.Symbol{},
	}

	if sc.docPresent {
		pkg.docStart = sc.docStart
		pkg.docEnd = sc.docEnd
	}

	for sc.proceed() {
		if submatches := constantDeclaration.FindSubmatchIndex(sc.line); len(submatches) > 0 {
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
		} else if submatches := arrayTypeDeclaration.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanArrayTypeDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(t)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if submatches := enumTypeDeclaration.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanEnumTypeDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(t)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if submatches := functionDeclaration.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[4]:submatches[5]])
			f, err := scanFunctionDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(f)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if submatches := recordTypeDeclaration.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanRecordTypeDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(t)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if submatches := subtypeDeclaration.FindSubmatchIndex(sc.line); len(submatches) > 0 {
			name := string(sc.line[submatches[2]:submatches[3]])
			t, err := scanSubtypeDeclaration(filepath, name, sc)
			err = pkg.AddSymbol(t)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
		} else if (len(end.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(name))) ||
			(len(endPackage.FindIndex(sc.line)) > 0) ||
			(len(endWithSemicolon.FindIndex(sc.line)) > 0) {
			pkg.codeEnd = sc.endIdx
			return pkg, nil
		}
	}

	return pkg, fmt.Errorf("package declaration end line not found")
}

func scanPackageInstantiation(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	pi := PackageInstantiation{
		Symbol{
			filepath:  filepath,
			name:      name,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		pi.docStart = sc.docStart
		pi.docEnd = sc.docEnd
	}

	for sc.proceed() {
		if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
			pi.codeEnd = sc.endIdx
			return pi, nil
		}
	}

	return pi, fmt.Errorf("package instantiation line with ';' not found")
}

func scanEnumTypeDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	t := Type{
		Symbol{
			filepath:  filepath,
			name:      name,
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		t.docStart = sc.docStart
		t.docEnd = sc.docEnd
	}

	for {
		if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
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
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		t.docStart = sc.docStart
		t.docEnd = sc.docEnd
	}

	for {
		if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
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
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		t.docStart = sc.docStart
		t.docEnd = sc.docEnd
	}

	for {
		if (len(end.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(name))) ||
			(len(endRecord.FindIndex(sc.line)) > 0) ||
			(len(endWithSemicolon.FindIndex(sc.line)) > 0) {
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
	c := Constant{
		Symbol{
			filepath:  filepath,
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		c.docStart = sc.docStart
		c.docEnd = sc.docEnd
	}

	syms := []symbol.Symbol{}

	for {
		if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
			c.codeEnd = sc.endIdx
			for _, n := range names {
				c.name = n
				syms = append(syms, c)
			}
			return syms, nil
		}
		if !sc.proceed() {
			break
		}
	}

	return syms, fmt.Errorf("constant declaration line with ';' not found")
}

func scanFunctionDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	f := Function{
		Symbol{
			filepath:  filepath,
			name:      name,
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		f.docStart = sc.docStart
		f.docEnd = sc.docEnd
	}

	for {
		if len(endsWithReturn.FindIndex(sc.line)) > 0 {
			f.codeEnd = sc.endIdx
			return f, nil
		}
		if !sc.proceed() {
			break
		}
	}

	return f, fmt.Errorf("function declaration line with return not found")
}

func scanSubtypeDeclaration(filepath string, name string, sc *scanContext) (symbol.Symbol, error) {
	t := Type{
		Symbol{
			filepath:  filepath,
			name:      name,
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		t.docStart = sc.docStart
		t.docEnd = sc.docEnd
	}

	for {
		if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
			t.codeEnd = sc.endIdx
			return t, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return t, fmt.Errorf("subtype declaration line with ';' not found")
}
