package vhdl

import (
	"regexp"
	"strings"
)

var startsWithWhenRegexp *regexp.Regexp = regexp.MustCompile(`^\s*when\b`)

var positiveResetPortMapRegexp *regexp.Regexp = regexp.MustCompile(`re?se?t((p)|(p_i)|(_p)|(_i)|(_p_i)|(_i_p))?\s*=>\s*(.+)`)
var positiveResetRegexp *regexp.Regexp = regexp.MustCompile(`re?se?t((p)|(p_i)|(_p)|(_i)|(_p_i)|(_i_p))?\b`)

var negativeResetPortMapRegexp *regexp.Regexp = regexp.MustCompile(`re?se?t((n)|(n_i)|(_n)|(_n_i)|(_i_n))\s*=>\s*(.+)`)
var negativeResetRegexp *regexp.Regexp = regexp.MustCompile(`re?se?t((n)|(n_i)|(_n)|(_n_i)|(_i_n))\b`)

var startsWithNotRegexp *regexp.Regexp = regexp.MustCompile(`^not((\s+)|(\s*\())`)

func checkResetPortMapping(line string) (string, bool) {
	if len(startsWithWhenRegexp.FindStringIndex(line)) > 0 {
		return "", true
	}

	if matches := positiveResetPortMapRegexp.FindStringSubmatch(line); len(matches) > 0 {
		if msg, ok := checkPositiveResetPortMapping(matches); !ok {
			return msg, ok
		}
	} else if matches := negativeResetPortMapRegexp.FindStringSubmatch(line); len(matches) > 0 {
		if msg, ok := checkNegativeResetPortMapping(matches); !ok {
			return msg, ok
		}
	}

	return "", true
}

func checkPositiveResetPortMapping(matches []string) (string, bool) {
	assignee := matches[len(matches)-1]

	if strings.HasPrefix(assignee, "'1'") {
		return "positive reset stuck to '1'", false
	}

	negated := false
	if len(startsWithNotRegexp.FindStringIndex(assignee)) > 0 {
		negated = true
	}

	reset := ""

	if len(negativeResetRegexp.FindStringIndex(assignee)) > 0 {
		reset = "negative"
	} else if len(positiveResetRegexp.FindStringIndex(assignee)) > 0 {
		reset = "positive"
	}

	if reset == "negative" && !negated {
		return "positive reset mapped to negative reset", false
	} else if reset == "positive" && negated {
		return "positive reset mapped to negated positive reset", false
	}

	return "", true
}

func checkNegativeResetPortMapping(matches []string) (string, bool) {
	assignee := matches[len(matches)-1]

	if strings.HasPrefix(assignee, "'0'") {
		return "negative reset stuck to '0'", false
	}

	negated := false
	if len(startsWithNotRegexp.FindStringIndex(assignee)) > 0 {
		negated = true
	}

	reset := ""

	if len(negativeResetRegexp.FindStringIndex(assignee)) > 0 {
		reset = "negative"
	} else if len(positiveResetRegexp.FindStringIndex(assignee)) > 0 {
		reset = "positive"
	}

	if reset == "positive" && !negated {
		return "negative reset mapped to positive reset", false
	} else if reset == "negative" && negated {
		return "negative reset mapped to negated negative reset", false
	}

	return "", true
}
