package vhdl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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

	var f *os.File
	var err error

	for {
		f, err = os.Open(filepath)
		if err != nil {
			if strings.HasSuffix(err.Error(), "too many open files") {
				time.Sleep(1e6)
			} else {
				log.Fatalf("check file %s: %v", filepath, err)
			}
		} else {
			break
		}
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

		if msg, ok := checkResetIfCondition(lineLower); !ok {
			fmt.Printf("%s:%d: %s\n%s\n\n", filepath, lineNum, msg, line)
		}

	}
}
