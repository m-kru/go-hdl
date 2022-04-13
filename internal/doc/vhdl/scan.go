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
	"github.com/m-kru/go-thdl/internal/doc/sym"
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
			name := string(sCtx.actualLine[sm[2]:sm[3]])
			ent, err := scanEntityDeclaration(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, ent)
		} else if sm := packageInstantiation.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.actualLine[sm[2]:sm[3]])
			pkgInst, err := scanPackageInstantiation(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, pkgInst)
		} else if sm := packageDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.actualLine[sm[2]:sm[3]])
			pkg, err := scanPackageDeclaration(filepath, name, &sCtx)
			if err != nil {
				log.Fatalf("%s: %v", filepath, err)
			}
			libContainer.AddSymbol(libName, pkg)
		}
	}
}

func scanEntityDeclaration(filepath string, name string, sc *scanContext) (sym.Symbol, error) {
	e := Entity{
		symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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

	return e, fmt.Errorf("'%s' entity declaration end line not found", name)
}

func scanPackageDeclaration(filepath string, name string, sc *scanContext) (sym.Symbol, error) {
	pkg := Package{
		filepath:  filepath,
		key:       strings.ToLower(name),
		name:      name,
		codeStart: sc.startIdx,
		Consts:    map[sym.ID]sym.Symbol{},
		Funcs:     map[sym.ID]sym.Symbol{},
		Procs:     map[sym.ID]sym.Symbol{},
		Prots:     map[sym.ID]sym.Symbol{},
		Types:     map[sym.ID]sym.Symbol{},
		Subtypes:  map[sym.ID]sym.Symbol{},
	}

	if sc.docPresent {
		pkg.docStart = sc.docStart
		pkg.docEnd = sc.docEnd
	}

	for sc.proceed() {
		var syms []sym.Symbol
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
			name := string(sc.actualLine[sm[2]:sm[3]])
			syms, err = scanArrayTypeDeclaration(filepath, name, sc)
		} else if sm := enumTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.actualLine[sm[2]:sm[3]])
			syms, err = scanEnumTypeDeclaration(filepath, name, sc)
		} else if sm := functionDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			impure := false
			if sm[2] > 0 && string(sc.actualLine[sm[2]:sm[3]]) == "impure" {
				impure = true
			}
			name := string(sc.actualLine[sm[4]:sm[5]])
			syms, err = scanFunctionDeclaration(filepath, impure, name, sc)
		} else if sm := procedureDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.actualLine[sm[2]:sm[3]])
			syms, err = scanProcedureDeclaration(filepath, name, sc)
		} else if sm := protectedTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.actualLine[sm[2]:sm[3]])
			syms, err = scanProtectedTypeDeclaration(filepath, name, sc)
		} else if sm := recordTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.actualLine[sm[2]:sm[3]])
			syms, err = scanRecordTypeDeclaration(filepath, name, sc)
		} else if sm := subtypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.actualLine[sm[2]:sm[3]])
			syms, err = scanSubtypeDeclaration(filepath, name, sc)
		} else if sm := someTypeDeclaration.FindSubmatchIndex(sc.line); len(sm) > 0 {
			name := string(sc.actualLine[sm[2]:sm[3]])
			syms, err = scanSomeTypeDeclaration(filepath, name, sc)
		} else if (len(end.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(strings.ToLower(name)))) ||
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

	return pkg, fmt.Errorf("'%s' package declaration end line not found", name)
}

func scanPackageInstantiation(filepath string, name string, sc *scanContext) (sym.Symbol, error) {
	pi := PackageInstantiation{
		symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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

	return pi, fmt.Errorf("'%s' package instantiation line with ';' not found", name)
}

func scanEnumTypeDeclaration(filepath string, name string, sc *scanContext) ([]sym.Symbol, error) {
	t := Type{
		kind: "enum",
		symbol: symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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
			return []sym.Symbol{t}, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' enum declaration line with ';' not found", name)
}

func scanArrayTypeDeclaration(filepath string, name string, sc *scanContext) ([]sym.Symbol, error) {
	t := Type{
		kind: "array",
		symbol: symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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
			return []sym.Symbol{t}, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' array declaration end line not found", name)
}

func scanProtectedTypeDeclaration(filepath string, name string, sc *scanContext) ([]sym.Symbol, error) {
	t := Protected{
		symbol: symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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
		if (len(end.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(strings.ToLower(name)))) ||
			(len(endProtected.FindIndex(sc.line)) > 0) ||
			(len(endWithSemicolon.FindIndex(sc.line)) > 0) {
			t.codeEnd = sc.endIdx
			return []sym.Symbol{t}, nil
		}
		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' protected declaration line with ';' not found", name)
}

func scanRecordTypeDeclaration(filepath string, name string, sc *scanContext) ([]sym.Symbol, error) {
	t := Type{
		kind: "record",
		symbol: symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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
		if (len(end.FindIndex(sc.line)) > 0 && bytes.Contains(sc.line, []byte(strings.ToLower(name)))) ||
			(len(endRecord.FindIndex(sc.line)) > 0) ||
			(len(endWithSemicolon.FindIndex(sc.line)) > 0) {
			t.codeEnd = sc.endIdx
			return []sym.Symbol{t}, nil
		}
		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' record declaration line with ';' not found", name)
}

func scanSomeTypeDeclaration(filepath string, name string, sc *scanContext) ([]sym.Symbol, error) {
	if !sc.lookahead() {
		return nil, fmt.Errorf("some type declaration line with type kind not found")
	}

	if len(startsWithProtected.FindIndex(sc.lookaheadLine)) > 0 {
		return scanProtectedTypeDeclaration(filepath, name, sc)
	} else if len(startsWithRecord.FindIndex(sc.lookaheadLine)) > 0 {
		return scanRecordTypeDeclaration(filepath, name, sc)
	} else if len(startsWithRoundBracket.FindIndex(sc.lookaheadLine)) > 0 {
		return scanEnumTypeDeclaration(filepath, name, sc)
	} else if len(startsWithArray.FindIndex(sc.lookaheadLine)) > 0 {
		return scanArrayTypeDeclaration(filepath, name, sc)
	}

	return nil, nil
}

func scanConstantDeclaration(filepath string, names []string, sc *scanContext) ([]sym.Symbol, error) {
	c := Constant{
		symbol{
			filepath:  filepath,
			lineNum:   sc.lineNum,
			codeStart: sc.startIdx,
		},
	}

	if sc.docPresent {
		c.docStart = sc.docStart
		c.docEnd = sc.docEnd
	}

	syms := []sym.Symbol{}

	for {
		if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
			c.codeEnd = sc.endIdx
			for _, n := range names {
				c.key = strings.ToLower(n)
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

func scanFunctionDeclaration(filepath string, impure bool, name string, sc *scanContext) ([]sym.Symbol, error) {
	f := Function{
		impure: impure,
		symbol: symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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
			return []sym.Symbol{f}, nil
		}
		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' function declaration line with return not found", name)
}

func scanProcedureDeclaration(filepath string, name string, sc *scanContext) ([]sym.Symbol, error) {
	p := Procedure{
		symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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
				return []sym.Symbol{p}, nil
			}
		} else {
			if len(endsWithSemicolon.FindIndex(sc.line)) > 0 {
				p.codeEnd = sc.endIdx
				return []sym.Symbol{p}, nil
			}
		}
		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' procedure declaration line with return not found", name)
}

func scanSubtypeDeclaration(filepath string, name string, sc *scanContext) ([]sym.Symbol, error) {
	t := Subtype{
		symbol{
			filepath:  filepath,
			key:       strings.ToLower(name),
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
			return []sym.Symbol{t}, nil
		}

		if !sc.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("subtype declaration line with ';' not found")
}
