package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/m-kru/go-thdl/internal/gen/gen"
	"github.com/m-kru/go-thdl/internal/utils"
	"github.com/m-kru/go-thdl/internal/vhdl/re"
	"strings"
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
		if unit.name != "" && len(unit.gens) > 0 {
			units = append(units, unit)
		}
	}

	for sCtx.scan() {
		if sm := re.ArchitectureDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			appendUnit()
			unit.name = string(sCtx.line[sm[2]:sm[3]])
			unit.lineNum = sCtx.lineNum
			unit.typ = "architecture"
			unit.gens = map[string]gen.Generable{}
		} else if sm := re.PackageDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			appendUnit()
			unit.name = string(sCtx.line[sm[2]:sm[3]])
			unit.lineNum = sCtx.lineNum
			unit.typ = "package"
			unit.gens = map[string]gen.Generable{}
		} else if sm := re.PackageBodyDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
			name := string(sCtx.line[sm[2]:sm[3]])
			if strings.ToLower(name) == strings.ToLower(unit.name) {
				if len(unit.gens) > 0 {
					appendUnit()
					unit.name = name
					unit.lineNum = sCtx.lineNum
					unit.typ = "package body"
				}
			}
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
	args := utils.ThdlGenArgs(sCtx.line)

	if !sCtx.scan() {
		return nil, fmt.Errorf("cannot scan generable, EOF")
	}

	if len(re.EmptyLine.FindIndex(sCtx.line)) > 0 || len(re.CommentLine.FindIndex(sCtx.line)) > 0 {
		return nil, nil
	} else if sm := re.EnumTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
		name := string(sCtx.line[sm[2]:sm[3]])
		return scanEnumTypeDeclaration(sCtx, name, args)
	}

	return nil, fmt.Errorf("line %d: cannot process line\n%s", sCtx.lineNum, sCtx.line)
}

// scanEnumTypeDeclaration assumes that current line in the scanContext contains the '(' character.
func scanEnumTypeDeclaration(sCtx *scanContext, name string, args []string) (*enum, error) {
	enum := enum{name: name, values: []string{}}

	err := enum.ParseArgs(args)
	if err != nil {
		return nil, fmt.Errorf("line %d: enum '%s': %v", sCtx.lineNum, name, err)
	}

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

		sCtx.scan()
	}

	return &enum, nil
}
