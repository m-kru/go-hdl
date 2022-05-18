package vhdl

import (
	"bufio"
	"bytes"
	"github.com/m-kru/go-thdl/internal/utils"
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
		log.Fatalf("scanning file %s: %v", filepath, err)
	}

	if len(units) == 0 {
		return
	}

	newContent := genNewFileContent(fileContent, units)

	// We can assume that the file already exists so the perm is discarded anyway.
	err = os.WriteFile(filepath, newContent, 0)
	if err != nil {
		log.Fatalf("writing file %s: %v", filepath, err)
	}
}

func genNewFileContent(fileContent []byte, units []unit) []byte {
	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	newContent := strings.Builder{}

	lineNum := uint(0)
	for _, _ = range units {
		for scanner.Scan() {
			lineNum += 1
			line := scanner.Bytes()
			newContent.Write(line)
		}
	}

	for scanner.Scan() {
		newContent.Write(scanner.Bytes())
	}

	return []byte(newContent.String())
}
