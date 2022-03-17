package vhdl

import (
	"fmt"
	"regexp"
	"strings"
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
	for i, _ := range pc.sensitivityList {
		if pc.sensitivityList[i] == s {
			return true
		}
	}
	return false
}

func checkProcessSensitivityList(line string, lineNum uint, pc *processContext) (string, bool) {
	if matches := processWithSensitivityListRegexp.FindStringSubmatch(line); len(matches) > 0 {
		pc.sensitivityListLineNum = lineNum
		pc.sensitivityListLine = line
		pc.sensitivityList = parseSensitivityList(matches[1])
	} else if len(endProcessRegexp.FindStringIndex(line)) > 0 {
		pc.sensitivityListLineNum = 0
		pc.sensitivityListLine = ""
		pc.sensitivityList = []string{}
		return "", true
	} else if len(processRegexp.FindStringIndex(line)) > 0 {
		if aux := startsWithBegin.FindStringIndex(line); len(aux) > 0 {
			return "", true
		}
		pc.sensitivityListLineNum = lineNum
		pc.sensitivityListLine = line
		pc.sensitivityList = []string{}
	}

	if matches := ingEdgeRegexp.FindStringSubmatch(line); len(matches) > 0 {
		// Ignore typical test bench use cases.
		if aux := startsWithWait.FindStringIndex(line); len(aux) > 0 {
			return "", true
		}
		// Ignore some rare, but synthesizable constructs.
		if strings.Contains(line, "<=") && strings.Contains(line, "when") {
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

		if !pc.inSensitivityList(signal) {
			return fmt.Sprintf(
					"'%s' not found in the sensitivity list\n%d:%s",
					signal, pc.sensitivityListLineNum, pc.sensitivityListLine,
				),
				false
		}
	}

	return "", true
}

func parseSensitivityList(s string) []string {
	list := []string{}
	elems := strings.Split(s, ",")
	for _, e := range elems {
		if e != "" {
			list = append(list, strings.Trim(e, " \t"))
		}
	}
	return list
}
