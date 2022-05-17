package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/m-kru/go-thdl/internal/gen/gen"
	_ "log"
	_ "os"
)

func scanFile(fileContent []byte) (map[string]gen.Generable, error) {
	gens := map[string]gen.Generable{}

	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	sCtx := scanContext{scanner: scanner}

	for sCtx.proceed() {
		if len(thdlGenLine.FindIndex(sCtx.line)) > 0 {
			gen, err := scanGenerable(&sCtx)
			if err != nil {
				return nil, err
			}
			gens[gen.Name()] = gen
		}
	}

	return gens, nil
}

func scanGenerable(sCtx *scanContext) (gen.Generable, error) {
	if !sCtx.proceed() {
		return nil, fmt.Errorf("cannot scan generable, EOF")
	}

	if sm := enumTypeDeclaration.FindSubmatchIndex(sCtx.line); len(sm) > 0 {
		name := string(sCtx.line[sm[2]:sm[3]])
		return scanEnumTypeDeclaration(sCtx, name)
	}

	panic("should never happen")
}

// scanEnumTypeDeclaration assumes that current line in the scanContext contains the '(' character.
func scanEnumTypeDeclaration(sCtx *scanContext, name string) (Enum, error) {
	enum := Enum{name: name, values: []string{}}

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
