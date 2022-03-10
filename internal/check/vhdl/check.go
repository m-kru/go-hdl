package vhdl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func Check(filepaths []string, wg *sync.WaitGroup) {
	var filesWg sync.WaitGroup

	for _, fp := range filepaths {
		filesWg.Add(1)
		go checkFile(fp, &filesWg)
	}

	filesWg.Wait()
	wg.Done()
}

func checkFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("check file %s: %v", filepath, err)
	}

	ioScanner := bufio.NewScanner(f)
	lineNum := 0
	for ioScanner.Scan() {
		lineNum += 1
		line := ioScanner.Text()
		lineLower := strings.ToLower(line)

		if msg, ok := checkClockPortMapping(lineLower); !ok {
			fmt.Printf("%s:%d: %s\n%s\n\n", filepath, lineNum, msg, line)
		}

		if msg, ok := checkResetPortMapping(lineLower); !ok {
			fmt.Printf("%s:%d: %s\n%s\n\n", filepath, lineNum, msg, line)
		}

	}
}
