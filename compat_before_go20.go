//go:build !go1.20
// +build !go1.20

package errors

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

const (
	rune1Max     = 1<<7 - 1
	rune2Max     = 1<<11 - 1
	rune3Max     = 1<<16 - 1
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
	maskx        = 0b00111111
	tx           = 0b10000000
	t2           = 0b11000000
	t3           = 0b11100000
	t4           = 0b11110000
)

// Copied from unicode/utf8/utf8.go available from Go 1.18 on.
func appendRune(p []byte, r rune) []byte {
	// This function is inlineable for fast handling of ASCII.
	if uint32(r) <= rune1Max {
		return append(p, byte(r))
	}
	return appendRuneNonASCII(p, r)
}

func appendRuneNonASCII(p []byte, r rune) []byte {
	// Negative values are erroneous. Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune2Max:
		return append(p, t2|byte(r>>6), tx|byte(r)&maskx)
	case i > utf8.MaxRune, surrogateMin <= i && i <= surrogateMax:
		r = utf8.RuneError
		fallthrough
	case i <= rune3Max:
		return append(p, t3|byte(r>>12), tx|byte(r>>6)&maskx, tx|byte(r)&maskx)
	default:
		return append(p, t4|byte(r>>18), tx|byte(r>>12)&maskx, tx|byte(r>>6)&maskx, tx|byte(r)&maskx)
	}
}

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
	b = appendRune(b, verb)
	return string(b)
}
