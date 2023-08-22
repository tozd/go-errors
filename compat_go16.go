//go:build !go1.20
// +build !go1.20

package errors

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Copied from fmt/print.go available from Go 1.20 on.
func formatString(state fmt.State, verb rune) string {
	var tmp [16]byte // Use a local buffer.
	b := append(tmp[:0], '%')
	for _, c := range " +-#0" { // All known flags
		if state.Flag(int(c)) { // The argument is an int for historical reasons.
			b = append(b, byte(c))
		}
	}
	if w, ok := state.Width(); ok {
		b = strconv.AppendInt(b, int64(w), 10)
	}
	if p, ok := state.Precision(); ok {
		b = append(b, '.')
		b = strconv.AppendInt(b, int64(p), 10)
	}
	b = utf8.AppendRune(b, verb)
	return string(b)
}
