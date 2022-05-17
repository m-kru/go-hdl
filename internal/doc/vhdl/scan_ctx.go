package vhdl

import (
	"bufio"
)

type scanContext struct {
	scanner *bufio.Scanner
	line    []byte // Lowercase actual line.
	lineNum uint32

	startIdx uint32 // Line start index.
	endIdx   uint32 // Line end index.

	docPresent bool
	docStart   uint32
	docEnd     uint32

	lookaheadLine []byte
}

// proceed returns false on EOF.
func (sc *scanContext) proceed() bool {
GETLINE:
	if sc.lookaheadLine != nil {
		sc.line = sc.lookaheadLine
		sc.lookaheadLine = nil
	} else if ok := sc.scanner.Scan(); !ok {
		return false
	} else {
		sc.line = sc.scanner.Bytes()
	}

	sc.lineNum += 1

	sc.startIdx = sc.endIdx
	sc.endIdx += uint32(len(sc.line)) + 1

	if len(emptyLine.FindIndex(sc.line)) > 0 {
		sc.docPresent = false
		goto GETLINE
	} else if len(commentLine.FindIndex(sc.line)) > 0 {
		sc.docEnd = sc.endIdx
		if !sc.docPresent {
			sc.docStart = sc.startIdx
			sc.docPresent = true
		}
	}

	return true
}

func (sc *scanContext) lookahead() bool {
	if sc.lookaheadLine != nil {
		panic("cannot lookahead more than one line")
	}

	if ok := sc.scanner.Scan(); !ok {
		return false
	}

	sc.lookaheadLine = sc.scanner.Bytes()

	return true
}

// symbolAdded function must be called whenever any symbol is added
// to any symbol container.
func (sc *scanContext) symbolAdded() {
	sc.docPresent = false
}
