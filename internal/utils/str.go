package utils

import (
	"strings"
)

func IsSingleLine(s string) bool {
	if strings.Count(s, "\n") > 1 {
		return false
	}

	return true
}

// FirstLine returns first line from the string s without '\n' rune.
func FirstLine(s string) string {
	return strings.Split(s, "\n")[0]
}

func Dewhitespace(s string) string {
	b := strings.Builder{}

	inIndent := true
	inWhitespace := false

	for _, r := range s {
		if inIndent {
			if r != ' ' && r != '\t' {
				b.WriteRune(r)
				inIndent = false
			}
		} else {
			if r == ' ' || r == '\t' {
				inWhitespace = true
			} else {
				if inWhitespace {
					inWhitespace = false
					b.WriteRune(' ')
				}
				b.WriteRune(r)
			}
		}
	}

	return b.String()
}
