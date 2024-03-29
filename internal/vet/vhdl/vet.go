package vhdl

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/m-kru/go-hdl/internal/utils"
	"github.com/m-kru/go-hdl/internal/vet/rprt"
)

var ignoreNextLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*--hdl:ignore`)
var ignoreThisLineRegExp *regexp.Regexp = regexp.MustCompile(`--hdl:ignore\s*$`)
var commentLineRegExp *regexp.Regexp = regexp.MustCompile(`^\s*--`)

func Vet(filepaths []string, wg *sync.WaitGroup) {
	var filesWg sync.WaitGroup

	for _, fp := range filepaths {
		filesWg.Add(1)
		go vetFile(fp, &filesWg)
	}

	filesWg.Wait()
	wg.Done()
}

func vetFile(filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if utils.IsIgnoredVHDLFile(filepath) {
		return
	}

	pCtx := processContext{sensitivityList: []string{}}

	f, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("error reading file %s: %v", filepath, err)
	}

	ioScanner := bufio.NewScanner(bytes.NewReader(f))
	lineNum := uint(0)
	ignoreNextLine := false
	for ioScanner.Scan() {
		lineNum += 1
		line := ioScanner.Bytes()

		if len(ignoreNextLineRegExp.FindIndex(line)) > 0 {
			ignoreNextLine = true
			continue
		} else if ignoreNextLine {
			ignoreNextLine = false
			continue
		} else if len(commentLineRegExp.FindIndex(line)) > 0 {
			continue
		} else if len(ignoreThisLineRegExp.FindIndex(line)) > 0 {
			continue
		}

		lineLower := bytes.ToLower(line)

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
