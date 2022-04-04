package utils

import (
	"strings"
)

var VHDLKeywords map[string]bool = map[string]bool{
	"array": true, "assert": true,
	"begin": true, "boolean": true, "buffer": true,
	"constant": true,
	"downto":   true,
	"end":      true, "entity": true,
	"failure": true, "false": true, "function": true,
	"generic": true,
	"impure":  true, "in": true, "inout": true, "integer": true, "is": true,
	"natural": true,
	"of":      true, "others": true, "out": true,
	"package": true, "port": true, "positive": true, "procedure": true, "pure": true,
	"range": true, "record": true, "report": true, "return": true,
	"severity": true, "signal": true, "signed": true, "std_logic": true, "std_logic_vector": true, "string": true, "subtype": true,
	"time": true, "to": true, "true": true, "type": true,
	"unsigned": true,
}

func VHDLTerminalBold(s string) string {
	b := strings.Builder{}

	inWord := false
	startIdx := 0
	endIdx := 0

	var prevR rune
	inString := false
	inComment := false

	for i, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' ||
			r == ':' || r == ';' || r == ',' || r == '(' || r == ')' {
			if r == '\n' || r == '\r' {
				inComment = false
			}
			if inWord {
				_, ok := VHDLKeywords[strings.ToLower(s[startIdx:endIdx])]
				if ok && !inComment && !inString {
					aux := "\033[1m" + s[startIdx:endIdx] + "\033[0m"
					_, _ = b.WriteString(aux)
				} else {
					_, _ = b.WriteString(s[startIdx:endIdx])
				}
			}
			inWord = false
			_, _ = b.WriteString(s[i : i+1])
		} else {
			if !inWord {
				startIdx = endIdx
				inWord = true
			}
		}
		endIdx += 1
		if r == '-' && prevR == '-' && !inString {
			inComment = true
		}
		if r == '"' {
			if prevR != '\\' {
				if inString {
					inString = false
				} else {
					inString = true
				}
			}
		}
		prevR = r
	}

	return b.String()
}
