package vhdl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/m-kru/go-thdl/internal/gen/gen"
	"github.com/m-kru/go-thdl/internal/utils"
	"github.com/m-kru/go-thdl/internal/vhdl/re"
	"log"
	"os"
	"strings"
	"sync"
)

func Gen(filepaths []string, wg *sync.WaitGroup) {
	var filesWg sync.WaitGroup

	for _, fp := range filepaths {
		filesWg.Add(1)
		go processFile(fp, &filesWg)
	}

	filesWg.Wait()
	wg.Done()
}

// processFile regenerates file only if there is anything to generate.
func processFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if utils.IsIgnoredVHDLFile(filepath) {
		return
	}

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("reading %s: %v", filepath, err)
	}

	units, err := scanFile(fileContent)
	if err != nil {
		log.Fatalf("%s: %v", filepath, err)
	}

	if len(units) == 0 {
		return
	}

	newContent, err := genNewFileContent(fileContent, units)
	if err != nil {
		log.Fatalf("%s: %v", filepath, err)
	}

	// We can assume that the file already exists so the perm is discarded anyway.
	err = os.WriteFile(filepath, newContent, 0)
	if err != nil {
		log.Fatalf("writing file %s: %v", filepath, err)
	}
}

func genNewFileContent(fileContent []byte, units []unit) ([]byte, error) {
	sCtx := scanContext{scanner: bufio.NewScanner(bytes.NewReader(fileContent))}
	b := strings.Builder{}

	write := func(line []byte) {
		b.Write(line)
		b.WriteRune('\n')
	}

	for _, u := range units {
		err := genDesignUnit(u, &sCtx, &b)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %v", u.typ, u.name, err)
		}
	}

	for sCtx.scan() {
		write(sCtx.line)
	}

	return []byte(b.String()), nil
}

func genDesignUnit(u unit, sCtx *scanContext, b *strings.Builder) error {
	inUnit := false
	gotoThdlEnd := false
	for {
		if !sCtx.scan() {
			if gotoThdlEnd {
				return fmt.Errorf("'--thdl:end' line not found")
			}
			break
		}

		if gotoThdlEnd {
			if len(thdlEndLine.FindIndex(sCtx.line)) > 0 {
				break
			}
			continue
		}

		if sCtx.lineNum == u.lineNum {
			inUnit = true
		}

		if inUnit {
			if u.typ == "architecture" {
			} else if u.typ == "package" {
				if len(thdlStartLine.FindIndex(sCtx.line)) > 0 {
					genPackage(u.gens, false, false, b)
					gotoThdlEnd = true
					continue
				} else if len(re.EndPackage.FindIndex(sCtx.line)) > 0 ||
					(len(re.End.FindIndex(sCtx.line)) > 0 && bytes.Contains(bytes.ToLower(sCtx.line), []byte(strings.ToLower(u.name)))) {
					genPackage(u.gens, false, true, b)
					b.Write(sCtx.line)
					b.WriteRune('\n')
					break
				}
			} else if u.typ == "package body" {
				if len(thdlStartLine.FindIndex(sCtx.line)) > 0 {
					genPackage(u.gens, true, false, b)
					gotoThdlEnd = true
					continue
				} else if len(re.EndPackageBody.FindIndex(sCtx.line)) > 0 ||
					(len(re.End.FindIndex(sCtx.line)) > 0 && bytes.Contains(bytes.ToLower(sCtx.line), []byte(strings.ToLower(u.name)))) {
					genPackage(u.gens, true, true, b)
					b.Write(sCtx.line)
					b.WriteRune('\n')
					break
				}
			} else {
				panic("should never happen")
			}
		}

		b.Write(sCtx.line)
		b.WriteRune('\n')
	}

	return nil
}

// body is false for package and true for package body.
func genPackage(gens map[string]gen.Generable, body bool, extraEmptyLines bool, b *strings.Builder) {
	if extraEmptyLines {
		b.WriteRune('\n')
	}

	b.WriteString(startCommentMsg)
	for _, g := range gens {
		var s string
		if body {
			s = g.GenDefinition([]string{})
		} else {
			s = g.GenDeclaration([]string{})
		}
		b.WriteString(s)
	}
	b.WriteString(endCommentMsg)

	if extraEmptyLines {
		b.WriteRune('\n')
	}
}
