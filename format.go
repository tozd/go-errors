package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Copied from fmt/print.go.
const (
	percentBangString = "%!"
	nilAngleString    = "<nil>"
	badPrecString     = "%!(BADPREC)"
)

const (
	stackTraceHelp     = "stack trace (most recent call first):\n"
	multipleErrorsHelp = "the above error joins multiple errors:\n"
	causeHelp          = "the above error was caused by the following error:\n"
)

// Similar to one in fmt/print.go.
func badVerb(s fmt.State, verb rune, arg interface{}) {
	_, _ = io.WriteString(s, percentBangString)
	_, _ = io.WriteString(s, string([]rune{verb}))
	_, _ = io.WriteString(s, "(")
	if arg != nil {
		_, _ = io.WriteString(s, reflect.TypeOf(arg).String())
		_, _ = io.WriteString(s, "=")
		fmt.Fprintf(s, "%v", arg)
	} else {
		_, _ = io.WriteString(s, nilAngleString)
	}
	_, _ = io.WriteString(s, ")")
}

// Copied from zerolog/console.go.
func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

func writeLinesPrefixed(st fmt.State, linePrefix, s string) {
	lines := strings.Split(s, "\n")
	// Trim empty lines at start.
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	// Trim empty lines at end.
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for _, line := range lines {
		_, _ = io.WriteString(st, linePrefix)
		_, _ = io.WriteString(st, line)
		_, _ = io.WriteString(st, "\n")
	}
}

func useFormatter(err error) bool {
	switch err.(type) { //nolint:errorlint
	case stackTracer, pkgStackTracer, goErrorsStackTracer, detailer:
		return false
	}

	_, ok := err.(fmt.Formatter) //nolint:errorlint
	return ok
}

func isForeignFormatter(err error) bool {
	// Our errors implement fmt.Formatter but we want to return false for them because
	// they just call into our Formatter which would lead to infinite recursion.
	switch err.(type) { //nolint:errorlint
	case *fundamental, *msgWithStack, *msgWithoutStack, *msgJoined, *withStack, *withoutStack, *cause:
		return false
	}

	_, ok := err.(fmt.Formatter) //nolint:errorlint
	return ok
}

func formatError(s fmt.State, indent int, err error) {
	linePrefix := ""
	if indent > 0 {
		width, ok := s.Width()
		if ok {
			linePrefix = strings.Repeat(strings.Repeat(" ", width), indent)
		} else {
			linePrefix = strings.Repeat("\t", indent)
		}
	}

	var cause error
	var errs []error
	precision, ok := s.Precision()
	if !ok {
		// We explicitly set it to 0.
		// See: https://github.com/golang/go/issues/61913
		precision = 0
	}

	if precision >= 2 && isForeignFormatter(err) || err == nil {
		writeLinesPrefixed(s, linePrefix, fmt.Sprintf(fmt.FormatString(s, 'v'), err))
		// Here we return because we assume formatting does recurse itself or at least
		// the user requested us to not recuse if the error implements fmt.Formatter.
		return
	}

	if useFormatter(err) {
		writeLinesPrefixed(s, linePrefix, fmt.Sprintf(fmt.FormatString(s, 'v'), err))
		// Here we still recurse ourselves because we assume formatting just formats the error and
		// does not recurse if it does not implement those interfaces which we checked in useFormatter.
		if precision == 1 || precision == 3 {
			cause, errs = causeOrJoined(err)
		}
	} else {
		formatMsg(s, linePrefix, err)
		var details map[string]interface{}
		if s.Flag('#') || precision == 1 || precision == 3 {
			details, cause, errs = allDetailsUntilCauseOrJoined(err)
		}
		if s.Flag('#') {
			formatDetails(s, linePrefix, details)
		}
		if s.Flag('+') {
			formatStack(s, linePrefix, err)
		}
	}

	if precision == 1 || precision == 3 { //nolint:nestif
		// It is possible that both cause and errs is set. A bit strange, but we allow it.
		// In that case we first recurse into errs and then into the cause, so that it is
		// clear which "error above" joins the errors (not the cause). Because cause is
		// not indented it is hopefully clearer that "error above" does not mean the last
		// error among joined but the one higher up before indentation.
		if len(errs) > 0 {
			if s.Flag('-') {
				if s.Flag(' ') {
					_, _ = io.WriteString(s, "\n")
				}
				writeLinesPrefixed(s, linePrefix, multipleErrorsHelp)
			}
			for _, e := range errs {
				// e should never be nil, but we still check.
				if e != nil {
					if s.Flag(' ') {
						_, _ = io.WriteString(s, "\n")
					}
					formatError(s, indent+1, e)
				}
			}
		}
		if cause != nil {
			if s.Flag('-') {
				if s.Flag(' ') {
					_, _ = io.WriteString(s, "\n")
				}
				writeLinesPrefixed(s, linePrefix, causeHelp)
			}
			if s.Flag(' ') {
				_, _ = io.WriteString(s, "\n")
			}
			formatError(s, indent, cause)
		}
	}
}

func formatMsg(s fmt.State, linePrefix string, err error) {
	writeLinesPrefixed(s, linePrefix, err.Error())
}

// Similar to writeFields in zerolog/console.go.
func formatDetails(s fmt.State, linePrefix string, details map[string]interface{}) {
	fields := make([]string, len(details))
	i := 0
	for field := range details {
		fields[i] = field
		i++
	}
	sort.Strings(fields)
	for _, field := range fields {
		value := details[field]
		var v string
		switch tValue := value.(type) {
		case string:
			if needsQuote(tValue) {
				v = strconv.Quote(tValue)
			} else {
				v = tValue
			}
		case json.Number:
			v = string(tValue)
		default:
			b, err := marshalWithoutEscapeHTML(tValue)
			if err != nil {
				v = fmt.Sprintf("[error: %v]", err)
			} else {
				v = string(b)
			}
		}
		writeLinesPrefixed(s, linePrefix, fmt.Sprintf("%s=%s\n", field, v))
	}
}

func formatStack(s fmt.State, linePrefix string, err error) {
	st := getExistingStackTrace(err)
	if len(st) == 0 {
		return
	}

	if s.Flag('-') {
		writeLinesPrefixed(s, linePrefix, stackTraceHelp)
	}
	var result string
	width, ok := s.Width()
	if ok {
		result = fmt.Sprintf("%+*v", width, StackFormatter(st))
	} else {
		result = fmt.Sprintf("%+v", StackFormatter(st))
	}
	writeLinesPrefixed(s, linePrefix, result)
}

// Formatter formats an error as text using the fmt.Formatter interface
// and marshals the error as JSON.
type Formatter struct {
	Error error
}

// Format formats the error as text according to the fmt.Formatter interface.
//
// The error does not have to necessary come from this package and it will be formatted
// in the same way if it implements interfaces used by this package (e.g., stackTracer
// or detailer interfaces. By default, only if those interfaces are not implemented,
// but fmt.Formatter interface is, formatting will be delegated to the error itself.
// You can change this default through format precision.
//
// Errors which do come from this package can be directly formatted by the fmt package
// in the same way as this function does as they implement fmt.Formatter interface.
// If you are not sure about the source of the error, it is safe to call this function
// on them as well.
//
// The following verbs are supported:
//
//	%s    the error message
//	%q    the quoted error message
//	%v    by default the same as %s
//
// You can control how is %v formatted through the width and precision arguments and
// flags. The width argument controls the width of the indent step in spaces. The default
// (no width argument) indents with a tab step.
// Width is passed through to the stack trace formatting.
//
// The following flags for %v are supported:
//
//	'#'   list details as key=value lines after the error message, when available
//	'+'   follow with the %+v formatted stack trace, if available
//	'-'   add human friendly messages to delimit parts of the text
//	' '   add extra newlines to separate parts of the text better
//
// Precision is specified by a period followed by a decimal number and enable
// modes of operation. The following modes are supported:
//
//	.0    do not change default behavior, this is the default
//	.1    recurse into error causes and joined errors
//	.2    prefer error's fmt.Formatter interface implementation if error implements it
//	.3    recurse into error causes and joined errors, but prefer fmt.Formatter
//	      interface implementation if any error implements it; this means that
//	      recursion stops if error's formatter does not recurse
//
// When any flag or non-zero precision mode is used, it is assured that the text
// ends with a newline, if it does not already do so.
func (f Formatter) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		precision, ok := s.Precision()
		if !ok {
			// We explicitly set it to 0.
			// See: https://github.com/golang/go/issues/61913
			precision = 0
		}
		if precision < 0 || precision > 3 {
			_, _ = io.WriteString(s, badPrecString)
			break
		}
		if s.Flag('#') || s.Flag('+') || s.Flag('-') || s.Flag(' ') || precision > 0 {
			formatError(s, 0, f.Error)
			break
		}
		fallthrough
	case 's':
		if f.Error != nil {
			_, _ = io.WriteString(s, f.Error.Error())
		} else {
			fmt.Fprintf(s, "%s", f.Error)
		}
	case 'q':
		if f.Error != nil {
			fmt.Fprintf(s, "%q", f.Error.Error())
		} else {
			fmt.Fprintf(s, "%q", f.Error)
		}
	default:
		badVerb(s, verb, f.Error)
	}
}
