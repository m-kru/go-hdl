package vhdl

import (
	"bytes"
	"fmt"
	"regexp"
	_ "strings"
)

var processRegexp *regexp.Regexp = regexp.MustCompile(`\bprocess\b`)
var endProcessRegexp *regexp.Regexp = regexp.MustCompile(`^\s*end\s+process\b`)
var startsWithBegin *regexp.Regexp = regexp.MustCompile(`^\s*begin\b`)
var startsWithWait *regexp.Regexp = regexp.MustCompile(`^\s*wait\b`)
var processWithSensitivityListRegexp *regexp.Regexp = regexp.MustCompile(`\bprocess\b\s*\((.*)\)`)
var ingEdgeRegexp *regexp.Regexp = regexp.MustCompile(`\(?\s*(ris|fall)ing_edge\s*\(\s*((\w*)|(\w*\s*\(\w\)))\s*\)\s*\)?`)

type processContext struct {
	sensitivityListLineNum uint
	sensitivityListLine    string
	sensitivityList        []string
}

// inSensitivityList return true if signal s is present in the sensitivity list.
func (pc processContext) inSensitivityList(s string) bool {
	for i := range pc.sensitivityList {
		if pc.sensitivityList[i] == s {
			return true
		}
	}
	return false
}

func checkProcessSensitivityList(line []byte, lineNum uint, pc *processContext) (string, bool) {
	if matches := processWithSensitivityListRegexp.FindSubmatch(line); len(matches) > 0 {
		pc.sensitivityListLineNum = lineNum
		pc.sensitivityListLine = string(line)
		pc.sensitivityList = parseSensitivityList(matches[1])
	} else if len(endProcessRegexp.FindIndex(line)) > 0 {
		pc.sensitivityListLineNum = 0
		pc.sensitivityListLine = ""
		pc.sensitivityList = []string{}
		return "", true
	} else if len(processRegexp.FindIndex(line)) > 0 {
		if aux := startsWithBegin.FindIndex(line); len(aux) > 0 {
			return "", true
		}
		pc.sensitivityListLineNum = lineNum
		pc.sensitivityListLine = string(line)
		pc.sensitivityList = []string{}
	}

	if matches := ingEdgeRegexp.FindSubmatch(line); len(matches) > 0 {
		// Ignore typical test bench use cases.
		if aux := startsWithWait.FindIndex(line); len(aux) > 0 {
			return "", true
		}
		// Ignore some rare, but synthesizable constructs.
		if bytes.Contains(line, []byte("<=")) && bytes.Contains(line, []byte("when")) {
			return "", true
		}

		signal := matches[2]

		if len(pc.sensitivityList) == 0 {
			return fmt.Sprintf(
					"'%s' found in the edge function, but sensitivity list is missing\n%d:%s",
					signal, pc.sensitivityListLineNum, pc.sensitivityListLine,
				),
				false
		}

		if !pc.inSensitivityList(string(signal)) {
			return fmt.Sprintf(
					"'%s' not found in the sensitivity list\n%d:%s",
					signal, pc.sensitivityListLineNum, pc.sensitivityListLine,
				),
				false
		}
	}

	return "", true
}

func parseSensitivityList(s []byte) []string {
	list := []string{}
	elems := bytes.Split(s, []byte(","))
	for _, e := range elems {
		if !bytes.Equal(e, []byte("")) {
			list = append(list, string(bytes.Trim(e, " \t")))
		}
	}
	return list
}
