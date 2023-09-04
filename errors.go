// Package errors provides errors with a recorded stack trace and optional
// structured details.
//
// The traditional error handling idiom in Go is roughly akin to
//
//	if err != nil {
//	        return err
//	}
//
// which when applied recursively up the call stack results in error reports
// without a stack trace or context.
// The errors package provides error handling primitives to annotate
// errors along the failure path in a way that does not destroy the
// original error.
//
// # Adding a stack trace to an error
//
// When interacting with code which returns errors without a stack trace,
// you can upgrade that error to one with a stack trace using errors.WithStack.
// For example:
//
//	func readAll(r io.Reader) ([]byte, errors.E) {
//	        data, err := ioutil.ReadAll(r)
//	        if err != nil {
//	                return nil, errors.WithStack(err)
//	        }
//	        return data, nil
//	}
//
// errors.WithStack records the stack trace at the point where it was called, so
// use it as close to where the error originated as you can get so that the
// recorded stack trace is more precise.
//
// The example above uses errors.E for the returned error type instead of the
// standard error type. This is not required, but it tells Go that you expect
// that the function returns only errors with a stack trace and Go type system
// then helps you find any cases where this is not so.
//
// errors.WithStack does not record the stack trace if it is already present in
// the error so it is safe to call it if you are unsure if the error contains
// a stack trace.
//
// Errors with a stack trace implement the following interface, returning program
// counters of function invocations:
//
//	type stackTracer interface {
//	        StackTrace() []uintptr
//	}
//
// You can use standard runtime.CallersFrames to obtain stack trace frame
// information (e.g., function name, source code file and line).
// You can also use errors.StackFormatter to format the stack trace.
//
// Although the stackTracer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// # Adding context to an error
//
// Sometimes an error occurs in a low-level function and the error messages
// returned from it are too low-level, too.
// You can use errors.Wrap to construct a new higher-level error while
// recording the original error as a cause.
//
//	image, err := readAll(imageFile)
//	if err != nil {
//	        return nil, errors.Wrap(err, "reading image failed")
//	}
//
// In the example above we returned a new error with a new message,
// hiding the low-level details. The returned error implements the
// following interface:
//
//	type causer interface {
//	        Cause() error
//	}
//
// which enables access to the underlying low-level error. You can also
// use errors.Cause to obtain the cause.
//
// Although the causer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// Sometimes you do not want to hide the error message but just add to it.
// You can use errors.WithMessage, which adds a prefix to the existing message,
// or errors.Errorf, which gives you more control over the new message.
//
//	errors.WithMessage(err, "reading image failed")
//	errors.Errorf("reading image failed (%w)", err)
//
// Example new messages could then be, respectively:
//
//	"reading image failed: connection error"
//	"reading image failed (connection error)"
//
// # Adding details to an error
//
// Errors returned by this package implement the detailer interface:
//
//	type detailer interface {
//	        Details() map[string]interface{}
//	}
//
// which enables access to a map with optional additional details about
// the error. Returned map can be modified in-place. You can also use
// errors.Details and errors.AllDetails to access details:
//
//	errors.Details(err)["url"] = "http://example.com"
//
// You can also use errors.WithDetails as an alternative to errors.WithStack
// if you also want to add details while recording the stack trace:
//
//	func readAll(r io.Reader, filename string) ([]byte, errors.E) {
//	        data, err := ioutil.ReadAll(r)
//	        if err != nil {
//	                return nil, errors.WithDetails(err, "filename", filename)
//	        }
//	        return data, nil
//	}
//
// # Working with the tree of errors
//
// Errors which implement the following standard unwrapper interfaces:
//
//	type unwrapper interface {
//	        Unwrap() error
//	}
//
//	type unwrapper interface {
//	        Unwrap() error[]
//	}
//
// form a tree of errors where a wrapping error points its parent,
// wrapped, error(s). Errors returned from this package implement this
// interface to return the original error or errors, when they exist.
// This enables us to have constant base errors which we annotate
// with a stack trace before we return them:
//
//	var ErrAuthentication = errors.Base("authentication error")
//	var ErrMissingPassphrase = errors.BaseWrap(ErrAuthentication, "missing passphrase")
//	var ErrInvalidPassphrase = errors.BaseWrap(ErrAuthentication, "invalid passphrase")
//
//	func authenticate(passphrase string) errors.E {
//	        if passphrase == "" {
//	                return errors.WithStack(ErrMissingPassphrase)
//	        } else if passphrase != "open sesame" {
//	                return errors.WithStack(ErrInvalidPassphrase)
//	        }
//	        return nil
//	}
//
// Or with details:
//
//	func authenticate(username, passphrase string) errors.E {
//	        if passphrase == "" {
//	                return errors.WithDetails(ErrMissingPassphrase, "username", username)
//	        } else if passphrase != "open sesame" {
//	                return errors.WithDetails(ErrInvalidPassphrase, "username", username)
//	        }
//	        return nil
//	}
//
// We can use errors.Is to determine which error has been returned:
//
//	if errors.Is(err, ErrMissingPassphrase) {
//	        fmt.Println("Please provide a passphrase to unlock the doors.")
//	}
//
// Works across the tree, too:
//
//	if errors.Is(err, ErrAuthentication) {
//	        fmt.Println("Failed to unlock the doors.")
//	}
//
// To access details, use:
//
//	errors.AllDetails(err)["username"]
//
// You can join multiple errors into one error by calling errors.Join.
// Join also records the stack trace at the point it was called.
//
// # Formatting errors
//
// All errors with a stack trace returned from this package implement fmt.Formatter
// interface and can be formatted by the fmt package. The following verbs are supported:
//
//	%s    the error message
//	%v    same as %s
//	%+v   together with the error message include also the stack trace,
//	      ends with a newline
package errors

import (
	"fmt"
	"strings"
	"unsafe"

	pkgerrors "github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() []uintptr
}

type pkgStackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

type goErrorsStackTracer interface {
	Callers() []uintptr
}

type causer interface {
	Cause() error
}

type unwrapper interface {
	Unwrap() error
}

type unwrapperJoined interface {
	Unwrap() []error
}

type detailer interface {
	Details() map[string]interface{}
}

func getExistingStackTrace(err error) []uintptr {
	for err != nil {
		switch e := err.(type) { //nolint:errorlint
		case stackTracer:
			return e.StackTrace()
		case pkgStackTracer:
			st := e.StackTrace()
			return *(*[]uintptr)(unsafe.Pointer(&st))
		case goErrorsStackTracer:
			return e.Callers()
		}
		c, ok := err.(causer) //nolint:errorlint
		if ok && c.Cause() != nil {
			return nil
		}
		e, ok := err.(unwrapperJoined) //nolint:errorlint
		if ok && len(e.Unwrap()) > 0 {
			return nil
		}
		err = Unwrap(err)
	}
	return nil
}

// prefixMessage eagerly builds a new message with the provided prefix.
// This is a trade-off which consumes more memory but allows one to cheaply
// call Error multiple times.
func prefixMessage(msg, prefix string) string {
	message := strings.Builder{}
	if len(prefix) > 0 {
		message.WriteString(prefix)
		if prefix[len(prefix)-1] != '\n' && len(msg) > 0 {
			message.WriteString(": ")
		}
	}
	if len(msg) > 0 {
		message.WriteString(msg)
	}
	return message.String()
}

// This is a trade-off which consumes more memory but allows one to cheaply
// call Error multiple times.
func joinMessages(errs []error) string {
	// Same implementation as standard library's joinError's Error.
	var b []byte
	for i, err := range errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

// E interface can be used in as a return type instead of the standard error
// interface to annotate which functions return an error with a stack trace
// and details.
// This is useful so that you know when you should use WithStack or WithDetails
// (for functions which do not return E) and when not (for functions which do
// return E).
//
// If you call WithStack on an error with a stack trace nothing bad happens
// (same error is simply returned), it just pollutes the code. So this
// interface is defined to help. (Calling WithDetails on an error with details
// adds an additional and independent layer of details on
// top of any existing details.)
type E interface {
	error
	stackTracer
	detailer
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) E {
	return &fundamentalError{
		msg:     message,
		stack:   callers(),
		details: nil,
	}
}

// Errorf return an error with the supplied message
// formatted according to a format specifier.
// It supports %w format verb to wrap an existing error.
// Errorf also records the stack trace at the point it was called,
// unless wrapped error already have a stack trace.
// If %w is provided multiple times, then a stack trace is always recorded.
func Errorf(format string, args ...interface{}) E {
	err := fmt.Errorf(format, args...) //nolint:goerr113
	var errs []error
	// Errorf itself maybe wrapped an error or errors so we can use a type switch here
	// and do not need to (and should not) use As to determine if that happened.
	switch u := err.(type) { //nolint:errorlint
	case unwrapperJoined:
		errs = u.Unwrap()
	case unwrapper:
		errs = []error{u.Unwrap()}
	}
	if len(errs) > 1 {
		return &msgJoinedError{
			errs:    errs,
			msg:     err.Error(),
			stack:   callers(),
			details: nil,
		}
	} else if len(errs) == 1 {
		unwrap := errs[0]
		st := getExistingStackTrace(unwrap)
		if len(st) == 0 {
			st = callers()
		}

		return &msgError{
			err:     unwrap,
			msg:     err.Error(),
			stack:   st,
			details: nil,
		}
	}

	return &fundamentalError{
		msg:     err.Error(),
		stack:   callers(),
		details: nil,
	}
}

// fundamentalError is an error that has a message and a stack,
// but does not wrap another error.
type fundamentalError struct {
	msg     string
	stack   []uintptr
	details map[string]interface{}
}

func (e *fundamentalError) Error() string {
	return e.msg
}

func (e *fundamentalError) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, formatString(s, verb), Formatter{e})
}

func (e fundamentalError) MarshalJSON() ([]byte, error) {
	return marshalJSONError(&e)
}

func (e *fundamentalError) StackTrace() []uintptr {
	return e.stack
}

func (e *fundamentalError) Details() map[string]interface{} {
	if e.details == nil {
		e.details = make(map[string]interface{})
	}
	return e.details
}

// msgError wraps another error and has its own stack and msg.
type msgError struct {
	err     error
	msg     string
	stack   []uintptr
	details map[string]interface{}
}

func (e *msgError) Error() string {
	return e.msg
}

func (e *msgError) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, formatString(s, verb), Formatter{e})
}

func (e msgError) MarshalJSON() ([]byte, error) {
	return marshalJSONError(&e)
}

func (e *msgError) Unwrap() error {
	return e.err
}

func (e *msgError) StackTrace() []uintptr {
	return e.stack
}

func (e *msgError) Details() map[string]interface{} {
	if e.details == nil {
		e.details = make(map[string]interface{})
	}
	return e.details
}

// msgJoinedError wraps multiple errors
// and has its own stack and msg.
type msgJoinedError struct {
	errs    []error
	msg     string
	stack   []uintptr
	details map[string]interface{}
}

func (e *msgJoinedError) Error() string {
	return e.msg
}

func (e *msgJoinedError) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, formatString(s, verb), Formatter{e})
}

func (e msgJoinedError) MarshalJSON() ([]byte, error) {
	return marshalJSONError(&e)
}

func (e *msgJoinedError) Unwrap() []error {
	return e.errs
}

func (e *msgJoinedError) StackTrace() []uintptr {
	return e.stack
}

func (e *msgJoinedError) Details() map[string]interface{} {
	if e.details == nil {
		e.details = make(map[string]interface{})
	}
	return e.details
}

// WithStack annotates err with a stack trace at the point WithStack was called,
// if err does not already have a stack trace.
// If err is nil, WithStack returns nil.
//
// Use this instead of Wrap when you just want to convert an existing error
// into one with a stack trace. Use it as close to where the error originated
// as you can get.
func WithStack(err error) E {
	if err == nil {
		return nil
	}

	e, ok := err.(E) //nolint:errorlint
	if ok {
		return e
	}

	st := getExistingStackTrace(err)
	if len(st) == 0 {
		st = callers()
	}

	return &noMsgError{
		err:     err,
		stack:   st,
		details: nil,
	}
}

// noMsgError wraps another error and has its
// own stack and but does not have its own msg.
type noMsgError struct {
	err     error
	stack   []uintptr
	details map[string]interface{}
}

func (e *noMsgError) Error() string {
	return e.err.Error()
}

func (e *noMsgError) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, formatString(s, verb), Formatter{e})
}

func (e noMsgError) MarshalJSON() ([]byte, error) {
	return marshalJSONError(&e)
}

func (e *noMsgError) Unwrap() error {
	return e.err
}

func (e *noMsgError) StackTrace() []uintptr {
	return e.stack
}

func (e *noMsgError) Details() map[string]interface{} {
	if e.details == nil {
		e.details = make(map[string]interface{})
	}
	return e.details
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// Wrapping is done even if err already has a stack trace.
// It records the original error as a cause.
// If err is nil, Wrap behaves like New.
//
// Use this when you want to make a new error,
// preserving the cause of the new error.
func Wrap(err error, message string) E {
	if err == nil {
		return &fundamentalError{
			msg:     message,
			stack:   callers(),
			details: nil,
		}
	}
	return &causeError{
		err:     err,
		msg:     message,
		stack:   callers(),
		details: nil,
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the supplied message
// formatted according to a format specifier.
// Wrapping is done even if err already has a stack trace.
// It records the original error as a cause.
// It does not support %w format verb (use %s instead if you
// need to incorporate cause's error message).
// If err is nil, Wrapf behaves like Errorf, but without
// support for %w format verb.
//
// Use this when you want to make a new error,
// preserving the cause of the new error.
func Wrapf(err error, format string, args ...interface{}) E {
	if err == nil {
		return &fundamentalError{
			msg:     fmt.Sprintf(format, args...),
			stack:   callers(),
			details: nil,
		}
	}
	return &causeError{
		err:     err,
		msg:     fmt.Sprintf(format, args...),
		stack:   callers(),
		details: nil,
	}
}

// causeError records another error as a causeError
// and has its own stack and msg.
type causeError struct {
	err     error
	msg     string
	stack   []uintptr
	details map[string]interface{}
}

func (e *causeError) Error() string {
	return e.msg
}

func (e *causeError) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, formatString(s, verb), Formatter{e})
}

func (e causeError) MarshalJSON() ([]byte, error) {
	return marshalJSONError(&e)
}

func (e *causeError) Unwrap() error {
	return e.err
}

func (e *causeError) Cause() error {
	return e.err
}

func (e *causeError) StackTrace() []uintptr {
	return e.stack
}

func (e *causeError) Details() map[string]interface{} {
	if e.details == nil {
		e.details = make(map[string]interface{})
	}
	return e.details
}

// WithMessage annotates err with a prefix message.
// If err does not have a stack trace, stack strace is recorded as well.
//
// It does not support controlling the delimiter. Use Errorf if you need that.
//
// If err is nil, WithMessage returns nil.
func WithMessage(err error, prefix string) E {
	if err == nil {
		return nil
	}

	st := getExistingStackTrace(err)
	if len(st) == 0 {
		st = callers()
	}

	return &msgError{
		err:     err,
		msg:     prefixMessage(err.Error(), prefix),
		stack:   st,
		details: nil,
	}
}

// WithMessagef annotates err with a prefix message
// formatted according to a format specifier.
// If err does not have a stack trace, stack strace is recorded as well.
//
// It does not support %w format verb or controlling the delimiter.
// Use Errorf if you need that.
//
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) E {
	if err == nil {
		return nil
	}

	st := getExistingStackTrace(err)
	if len(st) == 0 {
		st = callers()
	}

	return &msgError{
		err:     err,
		msg:     prefixMessage(err.Error(), fmt.Sprintf(format, args...)),
		stack:   st,
		details: nil,
	}
}

// Cause returns the result of calling the Cause method on err, if err's
// type contains a Cause method returning error.
// Otherwise, the err is unwrapped and the process is repeated.
// If unwrapping is not possible, Cause returns nil.
// Unwrapping stops if it encounters an error with
// Unwrap() method returning multiple errors.
func Cause(err error) error {
	for err != nil {
		c, ok := err.(causer) //nolint:errorlint
		if ok {
			cause := c.Cause()
			if cause != nil {
				return cause //nolint:wrapcheck
			}
		}
		e, ok := err.(unwrapperJoined) //nolint:errorlint
		if ok && len(e.Unwrap()) > 0 {
			return nil
		}
		err = Unwrap(err)
	}
	return err
}

// Unjoin returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning multiple errors.
// Otherwise, the err is unwrapped and the process is repeated.
// If unwrapping is not possible, Unjoin returns nil.
// Unwrapping stops if it encounters an error with the Cause
// method returning error.
func Unjoin(err error) []error {
	for err != nil {
		e, ok := err.(unwrapperJoined) //nolint:errorlint
		if ok {
			errs := e.Unwrap()
			if len(errs) > 0 {
				return errs
			}
		}
		c, ok := err.(causer) //nolint:errorlint
		if ok && c.Cause() != nil {
			return nil
		}
		err = Unwrap(err)
	}
	return nil
}

// Details returns the result of calling the Details method on err,
// if err's type contains a Details method returning initialized map.
// Otherwise, the err is unwrapped and the process is repeated.
// If unwrapping is not possible, Details returns nil.
// Unwrapping stops if it encounters an error with the Cause
// method returning error, or Unwrap() method returning
// multiple errors.
//
// You can modify returned map to modify err's details.
func Details(err error) map[string]interface{} {
	for err != nil {
		dd := detailsOf(err)
		if dd != nil {
			return dd
		}
		c, ok := err.(causer) //nolint:errorlint
		if ok && c.Cause() != nil {
			return nil
		}
		e, ok := err.(unwrapperJoined) //nolint:errorlint
		if ok && len(e.Unwrap()) > 0 {
			return nil
		}
		err = Unwrap(err)
	}
	return nil
}

// Returns details of the err if it implements detailer interface.
// It does not unwrap and recurse.
func detailsOf(err error) map[string]interface{} {
	if err == nil {
		return nil
	}
	d, ok := err.(detailer) //nolint:errorlint
	if ok {
		return d.Details()
	}
	return nil
}

// AllDetails returns a map build from calling the Details method on err
// and populating the map with key/value pairs which are not yet
// present. Afterwards, the err is unwrapped and the process is repeated.
// Unwrapping stops if it encounters an error with the Cause
// method returning error, or Unwrap() method returning
// multiple errors.
func AllDetails(err error) map[string]interface{} {
	res := make(map[string]interface{})
	for err != nil {
		for key, value := range detailsOf(err) {
			if _, ok := res[key]; !ok {
				res[key] = value
			}
		}
		c, ok := err.(causer) //nolint:errorlint
		if ok && c.Cause() != nil {
			return res
		}
		e, ok := err.(unwrapperJoined) //nolint:errorlint
		if ok && len(e.Unwrap()) > 0 {
			return res
		}
		err = Unwrap(err)
	}
	return res
}

// allDetailsUntilCauseOrJoined builds a map with details unwrapping errors
// until it hits a cause or joined errors, also returning it or them.
// This also means that it does not traverse errors returned by Join.
func allDetailsUntilCauseOrJoined(err error) (res map[string]interface{}, cause error, errs []error) { //nolint:revive,stylecheck,nonamedreturns
	res = make(map[string]interface{})
	cause = nil
	errs = nil

	for err != nil {
		for key, value := range detailsOf(err) {
			if _, ok := res[key]; !ok {
				res[key] = value
			}
		}
		c, ok := err.(causer) //nolint:errorlint
		if ok {
			cause = c.Cause()
		}
		e, ok := err.(unwrapperJoined) //nolint:errorlint
		if ok {
			errs = e.Unwrap()
		}
		if cause != nil || len(errs) > 0 {
			// It is possible that both cause and errs is set. A bit strange, but we allow it.
			return
		}
		err = Unwrap(err)
	}

	return
}

// causeOrJoined unwraps err repeatedly until it hits a cause or joined errors,
// returning it or them.
// This also means that it does not traverse errors returned by Join.
func causeOrJoined(err error) (cause error, errs []error) { //nolint:revive,stylecheck,nonamedreturns
	cause = nil
	errs = nil

	for err != nil {
		c, ok := err.(causer) //nolint:errorlint
		if ok {
			cause = c.Cause()
		}
		e, ok := err.(unwrapperJoined) //nolint:errorlint
		if ok {
			errs = e.Unwrap()
		}
		if cause != nil || len(errs) > 0 {
			// It is possible that both cause and errs is set. A bit strange, but we allow it.
			return
		}
		err = Unwrap(err)
	}

	return
}

func initializeDetails(err error) {
	for err != nil {
		detailsOf(err)
		err = Unwrap(err)
	}
}

// WithDetails wraps err implementing the detailer interface to access
// a map with optional additional details about the error.
//
// If err does not have a stack trace, then this call is equivalent
// to calling WithStack, annotating err with a stack trace as well.
//
// Use this when you have an err which implements stackTracer interface
// but does not implement detailer interface as well.
//
// It is also useful when err does implement detailer interface, but you want
// to reuse same err multiple times (e.g., pass same err to multiple
// goroutines), adding different details each time. Calling WithDetails
// wraps err and adds an additional and independent layer of details on
// top of any existing details.
//
// You can provide initial details by providing pairs of keys (strings)
// and values (interface{}).
func WithDetails(err error, kv ...interface{}) E {
	if err == nil {
		return nil
	}

	if len(kv)%2 != 0 {
		panic(New("odd number of arguments for initial details"))
	}

	// We always initialize map because details were explicitly asked for.
	initMap := make(map[string]interface{})
	for i := 0; i < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			panic(Errorf(`key "%v" must be a string, not %T`, kv[i], kv[i]))
		}
		initMap[key] = kv[i+1]
	}

	// Details were explicitly asked for, so we initialize them across
	// the whole stack of errors. It is useful to do this here so that
	// there are no race conditions if AllDetails on the same error is
	// called from multiple goroutines, racing which call will initialize
	// nil maps first.
	initializeDetails(err)

	// Even if err is of type E, we still wrap it into another withStack error to
	// have another layer of details. This is where it is different from WithStack.
	// We do not have to check for type E explicitly because E implements stackTracer
	// so getExistingStackTrace returns its stack trace.
	st := getExistingStackTrace(err)
	if len(st) == 0 {
		st = callers()
	}

	return &noMsgError{
		err:     err,
		stack:   st,
		details: initMap,
	}
}

// Join returns an error that wraps the given errors.
// Join also records the stack trace at the point it was called.
// Any nil error values are discarded.
// Join returns nil if errs contains no non-nil values.
// If there is only one non-nil value, Join behaves
// like WithStack on the non-nil value.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
func Join(errs ...error) E {
	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	} else if len(nonNilErrs) == 1 {
		err := nonNilErrs[0]

		e, ok := err.(E) //nolint:errorlint
		if ok {
			return e
		}

		st := getExistingStackTrace(err)
		if len(st) == 0 {
			st = callers()
		}

		return &noMsgError{
			err:     err,
			stack:   st,
			details: nil,
		}
	}

	return &msgJoinedError{
		errs:    nonNilErrs,
		msg:     joinMessages(nonNilErrs),
		stack:   callers(),
		details: nil,
	}
}
