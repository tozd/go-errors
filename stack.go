package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// frame represents a program counter inside a stack frame.
type frame runtime.Frame

// file returns the full path to the file that contains the
// function for this frame's pc.
func (f frame) file() string {
	if f.File == "" {
		return "unknown"
	}
	return f.File
}

// line returns the line number of source code of the
// function for this frame's pc.
func (f frame) line() int {
	return f.Line
}

// name returns the name of this function, if known.
func (f frame) name() string {
	if f.Function == "" {
		return "unknown"
	}
	return f.Function
}

// Format formats the frame as text according to the fmt.Formatter interface.
func (f frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, f.name())
			width, ok := s.Width()
			if ok {
				_, _ = io.WriteString(s, "\n")
				_, _ = io.WriteString(s, strings.Repeat(" ", width))
			} else {
				_, _ = io.WriteString(s, "\n\t")
			}
			_, _ = io.WriteString(s, f.file())
		default:
			_, _ = io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		_, _ = io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		_, _ = io.WriteString(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		_, _ = io.WriteString(s, ":")
		f.Format(s, 'd')
	default:
		badVerb(s, verb, f)
	}
}

func (f frame) MarshalJSON() ([]byte, error) {
	if f.Function == "" {
		return []byte("{}"), nil
	}

	return marshalWithoutEscapeHTML(&struct {
		Name string `json:"name,omitempty"`
		File string `json:"file,omitempty"`
		Line int    `json:"line,omitempty"`
	}{
		Name: f.name(),
		File: f.file(),
		Line: f.line(),
	})
}

// StackFormatter formats a provided stack trace as text using
// the fmt.Formatter interface and marshals the provided stack
// trace as JSON.
//
// Examples:
//
//	fmt.Sprintf("%+v", errors.StackFormatter{stack})
//	json.Marshal(errors.StackFormatter{stack})
type StackFormatter []uintptr

// Format formats the stack of frames as text according to the fmt.Formatter interface.
//
// The stack trace can come from errors in this package, from
// runtime.Callers, or from somewhere else.
//
// Each frame in the stack is formatted according to the format and is ended by a newline.
//
// The following verbs are supported:
//
//	%s	  lists the source file basename
//	%d    lists the source line number
//	%n    lists the short function name
//	%v	  equivalent to %s:%d
//
// StackFormat accepts flags that alter the formatting of some verbs, as follows:
//
//	%+s   lists the full function name and full compile-time path of the source file,
//	      separated by \n\t (<funcname>\n\t<path>)
//	%+v   lists the full function name and full compile-time path of the source file
//	      with the source line number, separated by \n\t
//	      (<funcname>\n\t<path>:<line>)
//
// StackFormat also accepts the width argument which controls the width of the indent
// step in spaces. The default (no width argument) indents with a tab step.
func (s StackFormatter) Format(st fmt.State, verb rune) {
	if len(s) == 0 {
		return
	}
	frames := runtime.CallersFrames(s)
	for {
		f, more := frames.Next()
		frame(f).Format(st, verb)
		_, _ = io.WriteString(st, "\n")
		if !more {
			break
		}
	}
}

// MarshalJSON marshals the stack of frames as JSON.
//
// JSON consists of an array of frame objects, each with
// (function) name, file (name), and line fields.
func (s StackFormatter) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}

	output := []byte{'['}
	frames := runtime.CallersFrames(s)
	first := true
	for {
		f, more := frames.Next()
		b, err := frame(f).MarshalJSON()
		if err != nil {
			return nil, WithStack(err)
		}
		if !first {
			output = append(output, ',')
		}
		first = false
		output = append(output, b...)
		if !more {
			break
		}
	}
	output = append(output, ']')
	return output, nil
}

func callers() StackFormatter {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) //nolint:gomnd
	var st StackFormatter = pcs[0:n]
	return st
}

// funcname removes the path prefix component of a function's name.
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
