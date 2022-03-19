package vhdl

import (
	"bytes"
	"regexp"
	_ "strings"
)

var positiveReset string = `re?se?t((p)|(p_i)|(_p)|(_i)|(_p_i)|(_i_p))?\b`
var negativeReset string = `re?se?t((n)|(n_i)|(_n)|(_n_i)|(_i_n))\b`

var startsWithWhenRegexp *regexp.Regexp = regexp.MustCompile(`^\s*when\b`)

var positiveResetPortMapRegexp *regexp.Regexp = regexp.MustCompile(positiveReset + `\s*=>\s*(.+)`)
var positiveResetRegexp *regexp.Regexp = regexp.MustCompile(positiveReset)

var negativeResetPortMapRegexp *regexp.Regexp = regexp.MustCompile(negativeReset + `\s*=>\s*(.+)`)
var negativeResetRegexp *regexp.Regexp = regexp.MustCompile(negativeReset)

var startsWithNotRegexp *regexp.Regexp = regexp.MustCompile(`^not((\s+)|(\s*\())`)

var positiveResetInvalidIfConditionRegexp = regexp.MustCompile(`^\s*if((\s+)|(\s*\(\s*))` + positiveReset + `\s*=\s*'0'((\s*)|(\s*\)\s*))then`)
var positiveResetInvalidIfConditionNoRHSRegexp = regexp.MustCompile(`^\s*if\s+not(\s+|(\s*\(\s*))` + positiveReset + `(\s*|(\s*\)\s*))then`)

var negativeResetInvalidIfConditionRegexp = regexp.MustCompile(`^\s*if((\s+)|(\s*\(\s*))` + negativeReset + `\s*=\s*'1'((\s*)|(\s*\)\s*))then`)
var negativeResetInvalidIfConditionNoRHSRegexp = regexp.MustCompile(`^\s*if((\s+)|(\s*\(\s*))` + negativeReset + `((\s+)|(\s*\)\s*))then`)

func checkResetPortMapping(line []byte) (string, bool) {
	if len(startsWithWhenRegexp.FindIndex(line)) > 0 {
		return "", true
	}

	if matches := positiveResetPortMapRegexp.FindSubmatch(line); len(matches) > 0 {
		if msg, ok := checkPositiveResetPortMapping(matches); !ok {
			return msg, ok
		}
	} else if matches := negativeResetPortMapRegexp.FindSubmatch(line); len(matches) > 0 {
		if msg, ok := checkNegativeResetPortMapping(matches); !ok {
			return msg, ok
		}
	}

	return "", true
}

func checkPositiveResetPortMapping(matches [][]byte) (string, bool) {
	assignee := matches[len(matches)-1]

	if bytes.HasPrefix(assignee, []byte("'1'")) {
		return "positive reset stuck to '1'", false
	}

	negated := false
	if len(startsWithNotRegexp.FindIndex(assignee)) > 0 {
		negated = true
	}

	reset := ""

	if len(negativeResetRegexp.FindIndex(assignee)) > 0 {
		reset = "negative"
	} else if len(positiveResetRegexp.FindIndex(assignee)) > 0 {
		reset = "positive"
	}

	if reset == "negative" && !negated {
		return "positive reset mapped to negative reset", false
	} else if reset == "positive" && negated {
		return "positive reset mapped to negated positive reset", false
	}

	return "", true
}

func checkNegativeResetPortMapping(matches [][]byte) (string, bool) {
	assignee := matches[len(matches)-1]

	if bytes.HasPrefix(assignee, []byte("'0'")) {
		return "negative reset stuck to '0'", false
	}

	negated := false
	if len(startsWithNotRegexp.FindIndex(assignee)) > 0 {
		negated = true
	}

	reset := ""

	if len(negativeResetRegexp.FindIndex(assignee)) > 0 {
		reset = "negative"
	} else if len(positiveResetRegexp.FindIndex(assignee)) > 0 {
		reset = "positive"
	}

	if reset == "positive" && !negated {
		return "negative reset mapped to positive reset", false
	} else if reset == "negative" && negated {
		return "negative reset mapped to negated negative reset", false
	}

	return "", true
}

func checkResetIfCondition(line []byte) (string, bool) {
	if msg, ok := checkPositiveResetIfCondition(line); !ok {
		return msg, ok
	}

	if msg, ok := checkNegativeResetIfCondition(line); !ok {
		return msg, ok
	}

	return "", true
}

func checkPositiveResetIfCondition(line []byte) (string, bool) {
	msg := "invalid positive reset condition"

	if len(positiveResetInvalidIfConditionRegexp.FindIndex(line)) > 0 {
		return msg, false
	}

	if len(positiveResetInvalidIfConditionNoRHSRegexp.FindIndex(line)) > 0 {
		return msg, false
	}

	return "", true
}

func checkNegativeResetIfCondition(line []byte) (string, bool) {
	msg := "invalid negative reset condition"

	if len(negativeResetInvalidIfConditionRegexp.FindIndex(line)) > 0 {
		return msg, false
	}

	if len(negativeResetInvalidIfConditionNoRHSRegexp.FindIndex(line)) > 0 {
		return msg, false
	}

	return "", true
}
