package utils

import (
	"log"
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
	"others":  true, "out": true,
	"package": true, "port": true, "positive": true, "procedure": true, "pure": true,
	"range": true, "record": true, "report": true,
	"severity": true, "signed": true, "std_logic": true, "std_logic_vector": true, "string": true, "subtype": true,
	"time": true, "to": true, "true": true, "type": true,
	"unsigned": true,
}

func VHDLTerminalBold(s string) string {
	var err error

	b := strings.Builder{}

	inWord := false
	startIdx := 0
	endIdx := 0

	for i, _ := range s {
		if s[i:i+1] == " " || s[i:i+1] == "\t" || s[i:i+1] == "\n" ||
			s[i:i+1] == ":" || s[i:i+1] == ";" || s[i:i+1] == "," ||
			s[i:i+1] == "(" || s[i:i+1] == ")" {
			if inWord {
				if _, ok := VHDLKeywords[strings.ToLower(s[startIdx:endIdx])]; ok {
					aux := "\033[1m" + s[startIdx:endIdx] + "\033[0m"
					_, err = b.WriteString(aux)
					if err != nil {
						log.Fatalf("VHDLTerminalBold: %v", err)
					}
				} else {
					_, err = b.WriteString(s[startIdx:endIdx])
					if err != nil {
						log.Fatalf("VHDLTerminalBold: %v", err)
					}
				}
			}
			inWord = false
			_, err = b.WriteString(s[i : i+1])
			if err != nil {
				log.Fatalf("VHDLTerminalBold: %v", err)
			}
		} else {
			if !inWord {
				startIdx = endIdx
				inWord = true
			}
		}
		endIdx += 1
	}

	return b.String()
}
