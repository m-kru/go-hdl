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
	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	newContent := strings.Builder{}

	write := func(line []byte) {
		newContent.Write(line)
		newContent.WriteRune('\n')
	}

	var inUnit bool
	var gotoThdlEnd bool
	lineNum := uint(0)
	for _, unit := range units {
		inUnit = false
		gotoThdlEnd = false
		for {
			if !scanner.Scan() {
				if gotoThdlEnd {
					return nil, fmt.Errorf(
						"%s %s, '--thdl:end' line not found", unit.typ, unit.name,
					)
				}
				break
			}

			lineNum += 1
			line := scanner.Bytes()

			if gotoThdlEnd {
				if len(thdlEndLine.FindIndex(line)) > 0 {
					break
				}
				continue
			}

			if lineNum == unit.lineNum {
				inUnit = true
			}

			if inUnit {
				if unit.typ == "architecture" {
				} else if unit.typ == "package" {
					if len(thdlStartLine.FindIndex(line)) > 0 {
						genPackage(unit.gens, false, &newContent)
						gotoThdlEnd = true
						continue
					} else if len(re.EndPackage.FindIndex(line)) > 0 {
						genPackage(unit.gens, true, &newContent)
						write(line)
						break
					}
				} else {
					panic("should never happen")
				}
			}

			write(line)
		}
	}

	for scanner.Scan() {
		write(scanner.Bytes())
	}

	return []byte(newContent.String()), nil
}

func genPackage(gens map[string]gen.Generable, extraEmptyLines bool, b *strings.Builder) {
	if extraEmptyLines {
		b.WriteRune('\n')
	}

	b.WriteString(startCommentMsg)
	for _, g := range gens {
		b.WriteString(g.GenDeclaration([]string{}))
	}
	b.WriteString(endCommentMsg)

	if extraEmptyLines {
		b.WriteRune('\n')
	}
}
