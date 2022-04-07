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

func scanFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if utils.IsIgnoredVHDLFile(filepath) {
		return
	}

	libName := docArgs.Lib(filepath)
	if libName == "" {
		libName = "work"
	}
	l := lib.MakeLibrary("vhdl", libName, LibSummary)
	libContainer.Add(&l)
	libContainer.Get(libName).AddFile(filepath)

	f, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", filepath, err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(f))

	sCtx := scanContext{scanner: scanner}

	for sCtx.proceed() {
		if sm := entityDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			ent, err := scanEntityDeclaration(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, ent)
		} else if sm := packageInstantiation.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			pkgInst, err := scanPackageInstantiation(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, pkgInst)
		} else if sm := packageDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
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
		var syms []symbol.Symbol
		var err error

		syms = nil

		if sm := constantDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			names := []string{}
			for i := 1; i < len(sm)/2; i++ {
				if sm[2*i] < 0 {
					continue
				}
				name := string(sc.line[sm[2*i]:sm[2*i+1]])
				if name[0] == ',' {
					name = strings.TrimSpace(name[1:])
				}
				names = append(names, name)
			}
			syms, err = scanConstantDeclaration(filepath, names, sc)
		} else if sm := arrayTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.line[sm[2]:sm[3]])
			syms, err = scanArrayTypeDeclaration(filepath, name, sc)
		} else if sm := enumTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.line[sm[2]:sm[3]])
			syms, err = scanEnumTypeDeclaration(filepath, name, sc)
		} else if sm := functionDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.line[sm[4]:sm[5]])
			syms, err = scanFunctionDeclaration(filepath, name, sc)
		} else if sm := procedureDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.line[sm[2]:sm[3]])
			syms, err = scanProcedureDeclaration(filepath, name, sc)
		} else if sm := recordTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.line[sm[2]:sm[3]])
			syms, err = scanRecordTypeDeclaration(filepath, name, sc)
		} else if sm := subtypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.line[sm[2]:sm[3]])
			syms, err = scanSubtypeDeclaration(filepath, name, sc)
		} else if sm := someTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.line[sm[2]:sm[3]])
			syms, err = scanSomeTypeDeclaration(filepath, name, sc)
		} else if (len(end.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(name))) ||
			(len(endPackage.FindIndex(sc.line)) > 0) ||
			(len(endWithSemicolon.FindIndex(sc.line)) > 0) {
			pkg.codeEnd = sc.endIdx
			return pkg, nil
		}

		if err != nil {
			return pkg, fmt.Errorf("package '%s': %v", name, err)
		}
		for _, s := range syms {
			err = pkg.AddSymbol(s)
			if err != nil {
				return pkg, fmt.Errorf("package '%s': %v", name, err)
			}
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

func scanEnumTypeDeclaration(filepath string, name string, sc *scanContext) ([]symbol.Symbol, error) {
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
			return []symbol.Symbol{t}, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("enum declaration line with ';' not found")
}

func scanArrayTypeDeclaration(filepath string, name string, sc *scanContext) ([]symbol.Symbol, error) {
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
			return []symbol.Symbol{t}, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("array declaration end line not found")
}

func scanRecordTypeDeclaration(filepath string, name string, sc *scanContext) ([]symbol.Symbol, error) {
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
			return []symbol.Symbol{t}, nil
		}
		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("record declaration line with ';' not found")
}

func scanSomeTypeDeclaration(filepath string, name string, sc *scanContext) ([]symbol.Symbol, error) {
	if !sc.lookahead() {
		return nil, fmt.Errorf("some type declaration line with type kind not found")
	}

	if len(startsWithRecord.FindIndex(sc.lookaheadLine)) > 0 {
		return scanRecordTypeDeclaration(filepath, name, sc)
	} else if len(startsWithRoundBracket.FindIndex(sc.lookaheadLine)) > 0 {
		return scanEnumTypeDeclaration(filepath, name, sc)
	} else if len(startsWithArray.FindIndex(sc.lookaheadLine)) > 0 {
		return scanArrayTypeDeclaration(filepath, name, sc)
	}

	return nil, nil
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

func scanFunctionDeclaration(filepath string, name string, sc *scanContext) ([]symbol.Symbol, error) {
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
			return []symbol.Symbol{f}, nil
		}
		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("function declaration line with return not found")
}

func scanProcedureDeclaration(filepath string, name string, sc *scanContext) ([]symbol.Symbol, error) {
	p := Procedure{
		Symbol{
			filepath:  filepath,
			name:      name,
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		p.docStart = sc.docStart
		p.docEnd = sc.docEnd
	}

	hasParams := false
	if bytes.Contains(sc.line, []byte("(")) {
		hasParams = true
	}

	for {
		if hasParams {
			if len(endsWithRoundBracketAndSemicolon.FindIndex(sc.line)) > 0 {
				p.codeEnd = sc.endIdx
				return []symbol.Symbol{p}, nil
			}
		} else {
			if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
				p.codeEnd = sc.endIdx
				return []symbol.Symbol{p}, nil
			}
		}
		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("function declaration line with return not found")
}

func scanSubtypeDeclaration(filepath string, name string, sc *scanContext) ([]symbol.Symbol, error) {
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
			return []symbol.Symbol{t}, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("subtype declaration line with ';' not found")
}
