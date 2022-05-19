// Package enc implements functions for different number encodings.
package enc

import (
	"strings"
)

func OneHot(i int, width int) string {
	b := strings.Builder{}
	for j := width - 1; j >= 0; j-- {
		r := '0'
		if j == i {
			r = '1'
		}
		b.WriteRune(r)
	}
	return b.String()
}
