package vhdl

import (
	"bufio"
	"bytes"
)

type scanContext struct {
	scanner *bufio.Scanner
	lineNum uint32
	line    []byte
}

// proceed returns false on EOF.
func (sc *scanContext) proceed() bool {
	sc.lineNum += 1

	if !sc.scanner.Scan() {
		return false
	}

	sc.line = sc.scanner.Bytes()

	return true
}

// decomment removes the comment at the end of the line.
func (sc *scanContext) decomment() {
	sc.line = bytes.Split(sc.line, []byte("--"))[0]
}
