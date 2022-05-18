package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/m-kru/go-thdl/internal/gen/gen"
	"github.com/m-kru/go-thdl/internal/vhdl/re"
	_ "log"
	_ "os"
)

// scanFile returns a list of units containing Generables within single file.
// If there is nothing to be generated within the unit, then
// the unit is not included in the list.
func scanFile(fileContent []byte) ([]unit, error) {
	units := []unit{}
	unit := unit{}

	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	sCtx := scanContext{scanner: scanner}

	appendUnit := func() {
		if unit.name != "" {
			units = append(units, unit)
		}
	}

	for sCtx.proceed() {
		if sm := re.ArchitectureDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			appendUnit()
			unit.name = string(sCtx.line[sm[2]:sm[3]])
			unit.lineNum = sCtx.lineNum
			unit.typ = "architecture"
			unit.gens = map[string]gen.Generable{}
		} else if sm := re.PackageDeclaration.FindIndex(sCtx.line); len(sm) > 0 {
			appendUnit()
			unit.name = string(sCtx.line[sm[2]:sm[3]])
			unit.lineNum = sCtx.lineNum
			unit.typ = "package"
			unit.gens = map[string]gen.Generable{}
		} else if len(thdlGenLine.FindIndex(sCtx.line)) > 0 {
			gen, err := scanGenerable(&sCtx)
			if err != nil {
				return nil, err
			}
			if gen != nil {
				unit.gens[gen.Name()] = gen
			}
		}
	}

	appendUnit()

	return units, nil
}

func scanGenerable(sCtx *scanContext) (gen.Generable, error) {
	if !sCtx.proceed() {
		return nil, fmt.Errorf("cannot scan generable, EOF")
	}

	if len(re.EmptyLine.FindIndex(sCtx.line)) > 0 || len(re.CommentLine.FindIndex(sCtx.line)) > 0 {
		return nil, nil
	} else if sm := re.EnumTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
		name := string(sCtx.line[sm[2]:sm[3]])
		return scanEnumTypeDeclaration(sCtx, name)
	}

	panic("should never happen")
}

// scanEnumTypeDeclaration assumes that current line in the scanContext contains the '(' character.
func scanEnumTypeDeclaration(sCtx *scanContext, name string) (enum, error) {
	enum := enum{name: name, values: []string{}}

	sCtx.line = bytes.Split((sCtx.line), []byte("("))[1]
	for {
		sCtx.decomment()
		vals := bytes.Split(sCtx.line, []byte(","))
		for _, v := range vals {
			if len(v) == 0 {
				continue
			}
			v = bytes.Trim(v, " \t")
			if v[len(v)-1] == ')' || v[len(v)-1] == ';' {
				v = v[:len(v)-1]
				v = bytes.Trim(v, " \t")
			}
			// Check one more time, needed in case of ");"
			if v[len(v)-1] == ')' {
				v = v[:len(v)-1]
			}
			if len(v) == 0 {
				continue
			}
			enum.values = append(enum.values, string(bytes.Trim(v, " \t")))
		}

		if bytes.Contains(sCtx.line, []byte(")")) {
			break
		}

		sCtx.proceed()
	}

	return enum, nil
}
