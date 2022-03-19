// Package rprt must be used to report violations to user.
package rprt

import (
	"fmt"
	"sync/atomic"
)

var violationCounter uint32 = 0

func ViolationCount() uint32 {
	return violationCounter
}

var fmtStr = "%s: %s\n%d:%s\n\n"

func Report(filepath string, msg string, lineNum uint, line []byte) {
	atomic.AddUint32(&violationCounter, 1)

	fmt.Printf(fmtStr, filepath, msg, lineNum, string(line))
}
