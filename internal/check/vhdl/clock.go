package vhdl

import (
	"bytes"
	"regexp"
)

var clockPortMapWithFrequenciesRegexp *regexp.Regexp = regexp.MustCompile(`cl(oc)?k_?(\d+).*=>.*cl(oc)?k_?(\d+)`)

func checkClockPortMapping(line []byte) (string, bool) {
	matches := clockPortMapWithFrequenciesRegexp.FindSubmatch(line)

	if len(matches) > 0 {
		if !bytes.Equal(matches[2], matches[4]) {
			return "clock frequency mismatch", false
		}
	}

	return "", true
}
