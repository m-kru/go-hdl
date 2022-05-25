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
	"github.com/m-kru/go-thdl/internal/vhdl/re"
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
	lib := lib.MakeLibrary("vhdl", libName, LibSummary)
	libContainer.Add(lib)
	lib = libContainer.Get(libName)
	lib.AddFile(filepath)

	f, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", filepath, err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(f))

	sCtx := scanContext{scanner: scanner}

	for sCtx.proceed() {
		var err error
		var ent Entity
		var sym sym.Symbol

		if sm := re.EntityDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			ent, err = scanEntityDeclaration(lib, filepath, name, &sCtx)
			ent.parent = lib
			sym = ent
		} else if sm := re.PackageInstantiation.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			sym, err = scanPackageInstantiation(lib, filepath, name, &sCtx)
		} else if sm := re.PackageDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			sym, err = scanPackageDeclaration(lib, filepath, name, &sCtx)
		}

		if err != nil {
			log.Fatalf("%s: %v", filepath, err)
		}
		if sym != nil {
			libContainer.AddSymbol(libName, sym)
			sCtx.symbolAdded()
		}
	}
}

func scanEntityDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) (Entity, error) {
	e := Entity{
		symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		e.docStart = sCtx.docStart
		e.docEnd = sCtx.docEnd
	}

	for sCtx.proceed() {
		if len(re.End.FindIndex(sCtx.line)) > 0 {
			e.codeEnd = sCtx.endIdx
			return e, nil
		}
	}

	return e, fmt.Errorf("'%s' entity declaration end line not found", name)
}

func scanPackageDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) (sym.Symbol, error) {
	pkg := Package{
		parent:    parent,
		filepath:  filepath,
		key:       strings.ToLower(name),
		name:      name,
		codeStart: sCtx.startIdx,
		Consts:    map[sym.ID]sym.Symbol{},
		Funcs:     map[sym.ID]sym.Symbol{},
		Procs:     map[sym.ID]sym.Symbol{},
		Types:     map[sym.ID]sym.Symbol{},
		Subtypes:  map[sym.ID]sym.Symbol{},
	}

	if sCtx.docPresent {
		pkg.docStart = sCtx.docStart
		pkg.docEnd = sCtx.docEnd
	}

	for sCtx.proceed() {
		var err error
		var syms []sym.Symbol

		if sm := re.ConstantDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			names := []string{}
			for i := 1; i < len(sm)/2; i++ {
				if sm[2*i] < 0 {
					continue
				}
				name := string(sCtx.line[sm[2*i]:sm[2*i+1]])
				if name[0] == ',' {
					name = strings.TrimSpace(name[1:])
				}
				names = append(names, name)
			}
			syms, err = scanConstantDeclaration(pkg, filepath, names, sCtx)
		} else if sm := re.ArrayTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanArrayTypeDeclaration(pkg, filepath, name, sCtx)
		} else if sm := re.EnumTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanEnumTypeDeclaration(pkg, filepath, name, sCtx)
		} else if sm := re.FunctionDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			impure := false
			if sm[2] > 0 && string(sCtx.line[sm[2]:sm[3]]) == "impure" {
				impure = true
			}
			name := string(sCtx.line[sm[4]:sm[5]])
			syms, err = scanFunctionDeclaration(pkg, filepath, impure, name, sCtx)
		} else if sm := re.ProcedureDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanProcedureDeclaration(pkg, filepath, name, sCtx)
		} else if sm := re.ProtectedTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanProtectedTypeDeclaration(pkg, filepath, name, sCtx)
		} else if sm := re.RecordTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanRecordTypeDeclaration(pkg, filepath, name, sCtx)
		} else if sm := re.SubtypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanSubtypeDeclaration(pkg, filepath, name, sCtx)
		} else if sm := re.SomeTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanSomeTypeDeclaration(parent, filepath, name, sCtx)
		} else if (len(re.End.FindIndex(sCtx.line)) > 0 && bytes.Contains(bytes.ToLower(sCtx.line), []byte(strings.ToLower(name)))) ||
			(len(re.EndPackage.FindIndex(sCtx.line)) > 0) ||
			(len(re.EndWithSemicolon.FindIndex(sCtx.line)) > 0) {
			pkg.codeEnd = sCtx.endIdx
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
		if syms != nil {
			sCtx.symbolAdded()
		}
	}

	return pkg, fmt.Errorf("'%s' package declaration end line not found", name)
}

func scanPackageInstantiation(parent sym.Symbol, filepath string, name string, sCtx *scanContext) (sym.Symbol, error) {
	pi := PackageInstantiation{
		symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		pi.docStart = sCtx.docStart
		pi.docEnd = sCtx.docEnd
	}

	for sCtx.proceed() {
		if len(re.EndsWithSemicolon.FindIndex(sCtx.line)) > 0 {
			pi.codeEnd = sCtx.endIdx
			return pi, nil
		}
	}

	return pi, fmt.Errorf("'%s' package instantiation line with ';' not found", name)
}

func scanEnumTypeDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	t := Enum{
		symbol: symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			lineNum:   sCtx.lineNum,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		t.docStart = sCtx.docStart
		t.docEnd = sCtx.docEnd
	}

	for {
		if len(re.EndsWithSemicolon.FindIndex(sCtx.line)) > 0 {
			t.codeEnd = sCtx.endIdx
			return []sym.Symbol{t}, nil
		}

		if !sCtx.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' enum declaration line with ';' not found", name)
}

func scanArrayTypeDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	t := Array{
		symbol: symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			lineNum:   sCtx.lineNum,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		t.docStart = sCtx.docStart
		t.docEnd = sCtx.docEnd
	}

	for {
		if len(re.EndsWithSemicolon.FindIndex(sCtx.line)) > 0 {
			t.codeEnd = sCtx.endIdx
			return []sym.Symbol{t}, nil
		}

		if !sCtx.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' array declaration end line not found", name)
}

func scanProtectedTypeDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	prot := Protected{
		parent:    parent,
		filepath:  filepath,
		key:       strings.ToLower(name),
		name:      name,
		lineNum:   sCtx.lineNum,
		codeStart: sCtx.startIdx,
		Funcs:     map[sym.ID]sym.Symbol{},
		Procs:     map[sym.ID]sym.Symbol{},
	}

	if sCtx.docPresent {
		prot.docStart = sCtx.docStart
		prot.docEnd = sCtx.docEnd
	}

	for {
		var err error
		var syms []sym.Symbol

		if sm := re.FunctionDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			impure := false
			if sm[2] > 0 && string(sCtx.line[sm[2]:sm[3]]) == "impure" {
				impure = true
			}
			name := string(sCtx.line[sm[4]:sm[5]])
			syms, err = scanFunctionDeclaration(prot, filepath, impure, name, sCtx)
		} else if sm := re.ProcedureDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			syms, err = scanProcedureDeclaration(prot, filepath, name, sCtx)
		} else if (len(re.End.FindIndex(sCtx.line)) > 0 && bytes.Contains(bytes.ToLower(sCtx.line), []byte(strings.ToLower(name)))) ||
			(len(re.EndProtected.FindIndex(sCtx.line)) > 0) ||
			(len(re.EndWithSemicolon.FindIndex(sCtx.line)) > 0) {
			prot.codeEnd = sCtx.endIdx
			return []sym.Symbol{prot}, nil
		}

		if err != nil {
			return nil, fmt.Errorf("protected '%s': %v", name, err)
		}
		for _, s := range syms {
			err = prot.AddSymbol(s)
			if err != nil {
				return nil, fmt.Errorf("protected '%s': %v", name, err)
			}
		}
		if syms != nil {
			sCtx.symbolAdded()
		}

		if !sCtx.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' protected declaration line with ';' not found", name)
}

func scanRecordTypeDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	t := Record{
		symbol: symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			lineNum:   sCtx.lineNum,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		t.docStart = sCtx.docStart
		t.docEnd = sCtx.docEnd
	}

	for {
		if (len(re.End.FindIndex(sCtx.line)) > 0 && bytes.Contains(bytes.ToLower(sCtx.line), []byte(strings.ToLower(name)))) ||
			(len(re.EndRecord.FindIndex(sCtx.line)) > 0) ||
			(len(re.EndWithSemicolon.FindIndex(sCtx.line)) > 0) {
			t.codeEnd = sCtx.endIdx
			return []sym.Symbol{t}, nil
		}
		if !sCtx.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' record declaration line with ';' not found", name)
}

func scanSomeTypeDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	if !sCtx.lookahead() {
		return nil, fmt.Errorf("some type declaration line with type kind not found")
	}

	if len(re.StartsWithProtected.FindIndex(sCtx.lookaheadLine)) > 0 {
		return scanProtectedTypeDeclaration(parent, filepath, name, sCtx)
	} else if len(re.StartsWithRecord.FindIndex(sCtx.lookaheadLine)) > 0 {
		return scanRecordTypeDeclaration(parent, filepath, name, sCtx)
	} else if len(re.StartsWithRoundBracket.FindIndex(sCtx.lookaheadLine)) > 0 {
		return scanEnumTypeDeclaration(parent, filepath, name, sCtx)
	} else if len(re.StartsWithArray.FindIndex(sCtx.lookaheadLine)) > 0 {
		return scanArrayTypeDeclaration(parent, filepath, name, sCtx)
	}

	return nil, nil
}

func scanConstantDeclaration(parent sym.Symbol, filepath string, names []string, sCtx *scanContext) ([]sym.Symbol, error) {
	c := Constant{
		symbol{
			parent:    parent,
			filepath:  filepath,
			lineNum:   sCtx.lineNum,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		c.docStart = sCtx.docStart
		c.docEnd = sCtx.docEnd
	}

	syms := []sym.Symbol{}

	for {
		if len(re.EndsWithSemicolon.FindIndex(sCtx.line)) > 0 {
			c.codeEnd = sCtx.endIdx
			for _, n := range names {
				c.key = strings.ToLower(n)
				c.name = n
				syms = append(syms, c)
			}
			return syms, nil
		}
		if !sCtx.proceed() {
			break
		}
	}

	return syms, fmt.Errorf("constant declaration line with ';' not found")
}

func scanFunctionDeclaration(parent sym.Symbol, filepath string, impure bool, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	f := Function{
		impure: impure,
		symbol: symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			lineNum:   sCtx.lineNum,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		f.docStart = sCtx.docStart
		f.docEnd = sCtx.docEnd
	}

	for {
		if len(re.EndsWithReturn.FindIndex(sCtx.line)) > 0 {
			f.codeEnd = sCtx.endIdx
			return []sym.Symbol{f}, nil
		}
		if !sCtx.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' function declaration line with return not found", name)
}

func scanProcedureDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	p := Procedure{
		symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			lineNum:   sCtx.lineNum,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		p.docStart = sCtx.docStart
		p.docEnd = sCtx.docEnd
	}

	hasParams := false
	if bytes.Contains(sCtx.line, []byte("(")) {
		hasParams = true
	}

	for {
		if hasParams {
			if len(re.EndsWithRoundBracketAndSemicolon.FindIndex(sCtx.line)) > 0 {
				p.codeEnd = sCtx.endIdx
				return []sym.Symbol{p}, nil
			}
		} else {
			if len(re.EndsWithSemicolon.FindIndex(sCtx.line)) > 0 {
				p.codeEnd = sCtx.endIdx
				return []sym.Symbol{p}, nil
			}
		}
		if !sCtx.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("'%s' procedure declaration line with return not found", name)
}

func scanSubtypeDeclaration(parent sym.Symbol, filepath string, name string, sCtx *scanContext) ([]sym.Symbol, error) {
	t := Subtype{
		symbol{
			parent:    parent,
			filepath:  filepath,
			key:       strings.ToLower(name),
			name:      name,
			lineNum:   sCtx.lineNum,
			codeStart: sCtx.startIdx,
		},
	}

	if sCtx.docPresent {
		t.docStart = sCtx.docStart
		t.docEnd = sCtx.docEnd
	}

	for {
		if len(re.EndsWithSemicolon.FindIndex(sCtx.line)) > 0 {
			t.codeEnd = sCtx.endIdx
			return []sym.Symbol{t}, nil
		}

		if !sCtx.proceed() {
			break
		}
	}

	return nil, fmt.Errorf("subtype declaration line with ';' not found")
}
