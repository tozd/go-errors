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

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    the source file
//	%d    the source line
//	%n    the function name
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the formatting of some verbs, as follows:
//
//	%+s   the full function name and full compile-time path of the source file,
//	      separated by \n\t (<funcname>\n\t<path>)
//	%+v   equivalent to %+s:%d
func (f frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
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

// stack represents a stack of program counters.
type stack []uintptr

// Format formats the stack of frames according to the fmt.Formatter interface.
// For each frame in the stack, ended by \n:
//
//	%s	  lists the source file
//	%d    lists the source line
//	%n    lists the function name
//	%v	  lists the source file and source line
//
// Format accepts flags that alter the formatting of some verbs, as follows:
//
//	%+s   lists the full function name and full compile-time path of the source file,
//	      separated by \n\t (<funcname>\n\t<path>)
//	%+v   lists the full function name and full compile-time path of the source file
//	      with the source line, separated by \n\t
//	      (<funcname>\n\t<path>:<line>)
func (s stack) Format(st fmt.State, verb rune) {
	if len(s) == 0 {
		return
	}
	frames := runtime.CallersFrames(s)
	for {
		f, more := frames.Next()
		frame(f).Format(st, verb)
		io.WriteString(st, "\n")
		if !more {
			break
		}
	}
}

func (s stack) MarshalJSON() ([]byte, error) {
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

func (s stack) StackTrace() []uintptr {
	return s
}

func callers() stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) //nolint:gomnd
	var st stack = pcs[0:n]
	return st
}

// funcname removes the path prefix component of a function's name.
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

// StackFormat formats the provided stack trace as text.
//
// The stack trace can come from errors in this package, from
// runtime.Callers, or from somewhere else.
//
// For each frame in the stack, ended by \n:
//
//	%s	  lists the source file
//	%d    lists the source line
//	%n    lists the function name
//	%v	  lists the source file and source line
//
// Format accepts flags that alter the formatting of some verbs, as follows:
//
//	%+s   lists the full function name and full compile-time path of the source file,
//	      separated by \n\t (<funcname>\n\t<path>)
//	%+v   lists the full function name and full compile-time path of the source file
//	      with the source line, separated by \n\t
//	      (<funcname>\n\t<path>:<line>)
func StackFormat(format string, s []uintptr) string {
	return fmt.Sprintf(format, stack(s))
}

// StackMarshalJSON marshals the provided stack trace as JSON.
//
// The stack trace can come from errors in this package, from
// runtime.Callers, or from somewhere else.
//
// JSON consists of an array of frame objects, each with
// (function) name, file (name), and line fields.
func StackMarshalJSON(s []uintptr) ([]byte, E) {
	b, err := stack(s).MarshalJSON()
	return b, WithStack(err)
}
