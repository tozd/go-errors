package errors

import (
	"fmt"
	"io"
	"unsafe"

	pkgerrors "github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() []uintptr
}

type pkgStackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

// E interface can be used in as a return type instead of the standard error
// interface to annotate which functions return an error with a stack trace.
// This is useful so that you know when you should use WithStack (for functions
// which do not return E) and when not (for functions which do return E).
// If you call WithStack on an error with a stack trace nothing bad happens
// (same error is simply returned), it just pollutes the code. So this
// interface is defined to help.
type E interface {
	error
	stackTracer
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) E {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf return an error with the supplied message
// formatted according to a format specifier.
// It supports %w format verb to wrap an existing error.
// Errorf also records the stack trace at the point it was called,
// unless wrapped error already have a stack trace.
//
// When formatting the returned error using %+v, formatting
// is not delegated to the wrapped error (when there is one),
// giving you full control of the message and formatted error.
func Errorf(format string, args ...interface{}) E {
	err := fmt.Errorf(format, args...)
	u, ok := err.(interface {
		Unwrap() error
	})
	if ok {
		unwrap := u.Unwrap()
		if _, ok := unwrap.(stackTracer); ok {
			return &errorf{
				unwrap,
				err.Error(),
			}
		} else if _, ok := unwrap.(pkgStackTracer); ok {
			return &errorf{
				unwrap,
				err.Error(),
			}
		}

		return &errorfWithStack{
			unwrap,
			err.Error(),
			callers(),
		}
	}

	return &fundamental{
		msg:   err.Error(),
		stack: callers(),
	}
}

// fundamental is an error that has a message and a stack,
// but does not wrap another error.
type fundamental struct {
	msg string
	stack
}

func (f *fundamental) Error() string {
	return f.msg
}

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(f.msg) > 0 {
				io.WriteString(s, f.msg)
				if f.msg[len(f.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "Stack trace (most recent call first):\n")
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

type errorf struct {
	error
	msg string
}

func (w *errorf) Error() string {
	return w.msg
}

func (w *errorf) Unwrap() error {
	return w.error
}

func (w *errorf) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(w.msg) > 0 {
				io.WriteString(s, w.msg)
				if w.msg[len(w.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "Stack trace (most recent call first):\n")
			stack(w.StackTrace()).Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.msg)
	case 'q':
		fmt.Fprintf(s, "%q", w.msg)
	}
}

func (w *errorf) StackTrace() []uintptr {
	switch e := w.error.(type) {
	case stackTracer:
		return e.StackTrace()
	case pkgStackTracer:
		st := e.StackTrace()
		return *(*[]uintptr)(unsafe.Pointer(&st))
	default:
		panic(New("not possible"))
	}
}

type errorfWithStack struct {
	error
	msg string
	stack
}

func (w *errorfWithStack) Error() string {
	return w.msg
}

func (w *errorfWithStack) Unwrap() error {
	return w.error
}

func (w *errorfWithStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(w.msg) > 0 {
				io.WriteString(s, w.msg)
				if w.msg[len(w.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "Stack trace (most recent call first):\n")
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.msg)
	case 'q':
		fmt.Fprintf(s, "%q", w.msg)
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called,
// if err does not already have a stack trace.
// If err is nil, WithStack returns nil.
//
// When formatting the returned error using %+v, formatting
// is delegated to the wrapped error. The stack trace is added only
// if the wrapped error does not already have it.
//
// Use this instead of Wrap when you just want to convert an existing error
// into one with a stack trace. Use it as close to where the error originated
// as you can get.
func WithStack(err error) E {
	if err == nil {
		return nil
	}
	if e, ok := err.(E); ok {
		return e
	} else if _, ok := err.(pkgStackTracer); ok {
		return &withPkgStack{
			err,
		}
	}

	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	stack
}

func (w *withStack) Unwrap() error {
	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			unwrap := fmt.Sprintf("%+v", w.error)
			if len(unwrap) > 0 {
				io.WriteString(s, unwrap)
				if unwrap[len(unwrap)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "Stack trace (most recent call first):\n")
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

type withPkgStack struct {
	error
}

func (w *withPkgStack) Unwrap() error {
	return w.error
}

func (w *withPkgStack) StackTrace() []uintptr {
	// We know error has pkgStackTracer interface because we construct it only then.
	st := w.error.(pkgStackTracer).StackTrace()
	return *(*[]uintptr)(unsafe.Pointer(&st))
}

func (w *withPkgStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.error)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// Wrapping is done even if err already has a stack trace.
// It records the original error as a cause.
// If err is nil, Wrap returns nil.
//
// When formatting the returned error using %+v, formatting
// of the cause is delegated to the wrapped error.
//
// Use this when you want to make a new error,
// preserving the cause of the new error.
func Wrap(err error, message string) E {
	if err == nil {
		return nil
	}
	return &wrapped{
		err,
		message,
		callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the supplied message
// formatted according to a format specifier.
// Wrapping is done even if err already has a stack trace.
// It records the original error as a cause.
// It does not support %w format verb.
// If err is nil, Wrapf returns nil.
//
// When formatting the returned error using %+v, formatting
// of the cause is delegated to the wrapped error.
//
// Use this when you want to make a new error,
// preserving the cause of the new error.
func Wrapf(err error, format string, args ...interface{}) E {
	if err == nil {
		return nil
	}
	return &wrapped{
		err,
		fmt.Sprintf(format, args...),
		callers(),
	}
}

type wrapped struct {
	error
	msg string
	stack
}

func (w *wrapped) Error() string {
	return w.msg
}

func (w *wrapped) Unwrap() error {
	return w.error
}

func (w *wrapped) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(w.msg) > 0 {
				io.WriteString(s, w.msg)
				if w.msg[len(w.msg)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			fmt.Fprintf(s, "Stack trace (most recent call first):\n")
			w.stack.Format(s, verb)
			unwrap := fmt.Sprintf("%+v", w.error)
			if len(unwrap) > 0 {
				io.WriteString(s, "\nThe above error was caused by the following error:\n\n")
				io.WriteString(s, unwrap)
				if unwrap[len(unwrap)-1] != '\n' {
					io.WriteString(s, "\n")
				}
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.msg)
	case 'q':
		fmt.Fprintf(s, "%q", w.msg)
	}
}

// WithMessage annotates err with a prefix message.
// If err does not have a stack trace, stack strace is recorded as well.
//
// It does not support controlling the delimiter. Use Errorf if you need that.
//
// When formatting the returned error using %+v, formatting
// is delegated to the wrapped error, prefixing it
// with the message. The stack trace is added only
// if the wrapped error does not already have it.
//
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) E {
	if err == nil {
		return nil
	}
	if e, ok := err.(E); ok {
		return &withMessage{
			e,
			message,
		}
	} else if _, ok := err.(pkgStackTracer); ok {
		return &withMessage{
			err,
			message,
		}
	}

	return &withMessageAndStack{
		err,
		message,
		callers(),
	}
}

// WithMessagef annotates err with a prefix message
// formatted according to a format specifier.
// If err does not have a stack trace, stack strace is recorded as well.
//
// It does not support %w format verb or controlling the delimiter.
// Use Errorf if you need that.
//
// When formatting the returned error using %+v, formatting
// is delegated to the wrapped error, prefixing it
// with the message. The stack trace is added only
// if the wrapped error does not already have it.
//
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) E {
	if err == nil {
		return nil
	}
	if _, ok := err.(stackTracer); ok {
		return &withMessage{
			err,
			fmt.Sprintf(format, args...),
		}
	} else if _, ok := err.(pkgStackTracer); ok {
		return &withMessage{
			err,
			fmt.Sprintf(format, args...),
		}
	}

	return &withMessageAndStack{
		err,
		fmt.Sprintf(format, args...),
		callers(),
	}
}

type withMessage struct {
	error
	msg string
}

func (w *withMessage) Error() string {
	message := ""
	unwrap := w.error.Error()
	if len(w.msg) > 0 {
		message += w.msg
		if w.msg[len(w.msg)-1] != '\n' && len(unwrap) > 0 {
			message += ": "
		}
	}
	if len(unwrap) > 0 {
		message += unwrap
	}
	return message
}

func (w *withMessage) Unwrap() error {
	return w.error
}

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			unwrap := fmt.Sprintf("%+v", w.error)
			if len(w.msg) > 0 {
				io.WriteString(s, w.msg)
				if w.msg[len(w.msg)-1] != '\n' && len(unwrap) > 0 {
					io.WriteString(s, ": ")
				}
			}
			if len(unwrap) > 0 {
				io.WriteString(s, unwrap)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func (w *withMessage) StackTrace() []uintptr {
	switch e := w.error.(type) {
	case stackTracer:
		return e.StackTrace()
	case pkgStackTracer:
		st := e.StackTrace()
		return *(*[]uintptr)(unsafe.Pointer(&st))
	default:
		panic(New("not possible"))
	}
}

type withMessageAndStack struct {
	error
	msg string
	stack
}

func (w *withMessageAndStack) Error() string {
	message := ""
	unwrap := w.error.Error()
	if len(w.msg) > 0 {
		message += w.msg
		if w.msg[len(w.msg)-1] != '\n' && len(unwrap) > 0 {
			message += ": "
		}
	}
	if len(unwrap) > 0 {
		message += unwrap
	}
	return message
}

func (w *withMessageAndStack) Unwrap() error {
	return w.error
}

func (w *withMessageAndStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			needsNewline := false
			unwrap := fmt.Sprintf("%+v", w.error)
			if len(w.msg) > 0 {
				io.WriteString(s, w.msg)
				if w.msg[len(w.msg)-1] != '\n' {
					if len(unwrap) > 0 {
						io.WriteString(s, ": ")
					} else {
						needsNewline = true
					}
				}
			}
			if len(unwrap) > 0 {
				io.WriteString(s, unwrap)
				if unwrap[len(unwrap)-1] != '\n' {
					needsNewline = true
				}
			}
			if needsNewline {
				io.WriteString(s, "\n")
			}
			fmt.Fprintf(s, "Stack trace (most recent call first):\n")
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
