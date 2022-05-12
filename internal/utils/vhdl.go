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
	"package": true, "port": true, "positive": true, "procedure": true, "protected": true, "pure": true,
	"range": true, "record": true, "report": true, "return": true,
	"severity": true, "signal": true, "signed": true, "std_logic": true, "std_logic_vector": true, "string": true, "subtype": true,
	"time": true, "to": true, "true": true, "type": true,
	"unsigned": true,
}

func vhdlBold(s string, prefix string, suffix string) string {
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
					aux := prefix + s[startIdx:endIdx] + suffix
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

func VHDLTerminalBold(s string) string {
	return vhdlBold(s, "\033[1m", "\033[0m")
}

func VHDLHTMLBold(s string) string {
	return vhdlBold(s, "<b>", "</b>")
}

// VHDLDecomment assumes that string has already passed Deindent() process.
func VHDLDecomment(s string) string {
	b := strings.Builder{}

	lineStart := true
	firstSpace := false
	potentialComment := false

	for _, r := range s {
		if r == '\n' || r == '\r' {
			lineStart = true
			firstSpace = false
			potentialComment = false
			b.WriteRune(r)
		} else if lineStart {
			if r == '-' {
				potentialComment = true
			} else {
				_, _ = b.WriteRune(r)
			}
			lineStart = false
		} else if potentialComment {
			if r != '-' {
				b.WriteRune('-')
				b.WriteRune(r)
				firstSpace = false
			} else {
				firstSpace = true
			}
			potentialComment = false
		} else if firstSpace {
			if r != ' ' {
				b.WriteRune(r)
			}
			firstSpace = false
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func VHDLDeindentDecomment(s string) string {
	// NOTE: This is not optimal implementation, as it iterates over
	// the string twice. However, currently it is not performance
	// bottleneck, so there is no need to optimize it.
	return VHDLDecomment(Deindent(s))
}
