package vhdl

import (
	"regexp"
)

var clockPortMapWithFrequenciesRegexp *regexp.Regexp = regexp.MustCompile(`cl(oc)?k_?(\d+).*=>.*cl(oc)?k_?(\d+)`)

func checkClockPortMapping(line string) (string, bool) {
	matches := clockPortMapWithFrequenciesRegexp.FindStringSubmatch(line)

	if len(matches) > 0 {
		if matches[2] != matches[4] {
			return "clock frequency mismatch", false
		}
	}

	return "", true
}
