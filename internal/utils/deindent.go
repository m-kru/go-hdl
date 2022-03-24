package utils

import (
	"strings"
)

func Deindent(s string) string {
	b := strings.Builder{}

	checkingIndent := true
	indent := 0
	indentCnt := 0

	for _, r := range s {
		if checkingIndent {
			if r == ' ' || r == '\t' {
				indent += 1
			} else {
				checkingIndent = false
				indentCnt = indent
				_, _ = b.WriteRune(r)
			}
		} else {
			if r == '\n' || r == '\r' {
				indentCnt = 0
				_, _ = b.WriteRune(r)
			} else if indentCnt < indent {
				indentCnt += 1
			} else {
				_, _ = b.WriteRune(r)
			}
		}
	}

	return b.String()
}
