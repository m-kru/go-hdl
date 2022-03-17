package vhdl

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/m-kru/go-thdl/internal/check/rprt"
)

var ignoreNextLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*--thdl:ignore`)
var ignoreThisLineRegExp *regexp.Regexp = regexp.MustCompile(`--thdl:ignore\s*$`)
var commentLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*--`)

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
	defer f.Close()

	pCtx := processContext{sensitivityList: []string{}}

	ioScanner := bufio.NewScanner(f)
	lineNum := uint(0)
	ignoreNextLine := false
	for ioScanner.Scan() {
		lineNum += 1
		line := ioScanner.Text()

		if len(ignoreNextLineRegExp.FindStringIndex(line)) > 0 {
			ignoreNextLine = true
			continue
		} else if ignoreNextLine {
			ignoreNextLine = false
			continue
		} else if len(commentLineRegExp.FindStringIndex(line)) > 0 {
			continue
		} else if len(ignoreThisLineRegExp.FindStringIndex(line)) > 0 {
			continue
		}

		lineLower := strings.ToLower(line)

		if msg, ok := checkClockPortMapping(lineLower); !ok {
			rprt.Report(filepath, msg, lineNum, line)
		}

		if msg, ok := checkResetPortMapping(lineLower); !ok {
			rprt.Report(filepath, msg, lineNum, line)
		}

		if msg, ok := checkResetIfCondition(lineLower); !ok {
			rprt.Report(filepath, msg, lineNum, line)
		}

		if msg, ok := checkProcessSensitivityList(lineLower, lineNum, &pCtx); !ok {
			rprt.Report(filepath, msg, lineNum, line)
		}
	}
}
