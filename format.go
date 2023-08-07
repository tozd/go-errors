package errors

import (
	"fmt"
	"io"
	"reflect"
)

// Copied from fmt/print.go.
const (
	percentBangString = "%!"
	nilAngleString    = "<nil>"
)

func badVerb(s fmt.State, verb rune, arg interface{}) {
	io.WriteString(s, percentBangString)
	io.WriteString(s, string([]rune{verb}))
	io.WriteString(s, "(")
	if arg != nil {
		io.WriteString(s, reflect.TypeOf(arg).String())
		io.WriteString(s, "=")
		fmt.Fprintf(s, "%v", arg)
	} else {
		io.WriteString(s, nilAngleString)
	}
	io.WriteString(s, ")")
}
