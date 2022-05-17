package vhdl

import (
	"github.com/m-kru/go-thdl/internal/utils"
	"log"
	"os"
	"sync"
)

func Gen(filepaths []string, wg *sync.WaitGroup) {
	var filesWg sync.WaitGroup

	for _, fp := range filepaths {
		filesWg.Add(1)
		go genFile(fp, &filesWg)
	}

	filesWg.Wait()
	wg.Done()
}

// genFile regenerates file only if there is anything to generate.
func genFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if utils.IsIgnoredVHDLFile(filepath) {
		return
	}

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("reading %s: %v", filepath, err)
	}

	gens, err := scanFile(fileContent)
	if err != nil {
		log.Fatalf("scanning file %s: %v", filepath, err)
	}

	if len(gens) == 0 {
		return
	}
}
