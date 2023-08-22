package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

type stackTracer interface {
	StackTrace() []uintptr
}

func w() string {
	return "%w"
}

type errorWithFormat struct {
	vMsg string
}

func (*errorWithFormat) Error() string {
	return "foobar"
}

func (e *errorWithFormat) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.vMsg)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, "foobar1")
	case 'q':
		fmt.Fprintf(s, "%q", "foobar2")
	}
}

type errorWithFormatAndStack struct {
	vMsg string
}

func (*errorWithFormatAndStack) Error() string {
	return "foobar"
}

func (e *errorWithFormatAndStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.vMsg)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, "foobar1")
	case 'q':
		fmt.Fprintf(s, "%q", "foobar2")
	}
}

func (e *errorWithFormatAndStack) StackTrace() []uintptr {
	return nil
}

type errorWithCauseAndWrap struct {
	msg   string
	cause error
	wrap  error
}

func (e *errorWithCauseAndWrap) Error() string {
	return e.msg
}

func (e *errorWithCauseAndWrap) Cause() error {
	return e.cause
}

func (e *errorWithCauseAndWrap) Unwrap() error {
	return e.wrap
}

type testStruct struct {
	err    error
	want   string
	format string
	stack  int
	extra  int
}

var tests []testStruct

func init() {
	// We make errors inside a function so that the stack trace is
	// different from errors made through errors.* call.
	parentErr, noMsgErr := func() (errors.E, errors.E) {
		return errors.New("parent"), errors.New("")
	}()
	parentWithFormat1Err, parentWithFormat2Err := func() (errors.E, errors.E) {
		return errors.WithStack(&errorWithFormat{"foobar\nmore data"}), errors.WithStack(&errorWithFormat{"foobar\nmore data\n"})
	}()
	parentPkgError := func() error {
		return pkgerrors.New("parent")
	}()

	parentNoStackErr := stderrors.New("parent")
	noMsgNoStackErr := stderrors.New("")

	// Current stack plus call to errors.*.
	currentStackSize := len(callers()) + 1
	parentErrStackSize := len(parentErr.StackTrace())

	tests = append(tests, []testStruct{
		// errors.New
		{errors.New(""), "", "% +-.1v", currentStackSize, 0},
		{errors.New("foo"), "foo", "% +-.1v", currentStackSize, 1},
		{errors.New("string with format specifiers: %v"), "string with format specifiers: %v", "% +-.1v", currentStackSize, 1},
		{errors.New("foo with newline\n"), "foo with newline\n", "% +-.1v", currentStackSize, 1},

		// errors.Errorf without %w
		{errors.Errorf(""), "", "% +-.1v", currentStackSize, 0},
		{errors.Errorf("read error without format specifiers"), "read error without format specifiers", "% +-.1v", currentStackSize, 1},
		{errors.Errorf("read error with %d format specifier", 1), "read error with 1 format specifier", "% +-.1v", currentStackSize, 1},
		{errors.Errorf("read error with newline\n"), "read error with newline\n", "% +-.1v", currentStackSize, 1},

		// errors.Errorf with %w and parent with stack
		{errors.Errorf("read error with parent: %w", parentErr), "read error with parent: parent", "% +-.1v", parentErrStackSize, 1},
		{errors.Errorf(`read error (parent "%w")`, parentErr), `read error (parent "parent")`, "% +-.1v", parentErrStackSize, 1},
		{errors.Errorf("read error (%w) and newline\n", parentErr), "read error (parent) and newline\n", "% +-.1v", parentErrStackSize, 1},
		{errors.Errorf("%w", noMsgErr), "", "% +-.1v", parentErrStackSize, 0},

		// errors.Errorf with %w and parent without stack
		{errors.Errorf("read error with parent: %w", parentNoStackErr), "read error with parent: parent", "% +-.1v", currentStackSize, 1},
		{errors.Errorf(`read error (parent "%w")`, parentNoStackErr), `read error (parent "parent")`, "% +-.1v", currentStackSize, 1},
		{errors.Errorf("read error (%w) and newline\n", parentNoStackErr), "read error (parent) and newline\n", "% +-.1v", currentStackSize, 1},
		{errors.Errorf("%w", noMsgNoStackErr), "", "% +-.1v", currentStackSize, 0},

		// errors.WithStack and parent without stack
		{errors.WithStack(io.EOF), "EOF", "% +-.1v", currentStackSize, 1},
		{errors.WithStack(errors.Base("EOF")), "EOF", "% +-.1v", currentStackSize, 1},
		{errors.WithStack(errors.Base("")), "", "% +-.1v", currentStackSize, 0},
		{errors.WithStack(errors.Base("foobar\n")), "foobar\n", "% +-.1v", currentStackSize, 1},

		// errors.WithStack and parent with stack
		{errors.WithStack(parentErr), "parent", "% +-.1v", parentErrStackSize, 1},
		{errors.WithStack(noMsgErr), "", "% +-.1v", parentErrStackSize, 0},

		// errors.WithStack and parent with custom %+v which is ignored
		{errors.WithStack(&errorWithFormat{"foobar\nmore data"}), "foobar", "% +-.1v", currentStackSize, 1},
		{errors.WithStack(&errorWithFormat{"foobar\nmore data\n"}), "foobar", "% +-.1v", currentStackSize, 1},

		// errors.WithDetails and parent without stack
		{errors.WithDetails(io.EOF), "EOF", "% +-.1v", currentStackSize, 1},
		{errors.WithDetails(errors.Base("EOF")), "EOF", "% +-.1v", currentStackSize, 1},
		{errors.WithDetails(errors.Base("")), "", "% +-.1v", currentStackSize, 0},
		{errors.WithDetails(errors.Base("foobar\n")), "foobar\n", "% +-.1v", currentStackSize, 1},

		// errors.WithDetails and parent with stack
		{errors.WithDetails(parentErr), "parent", "% +-.1v", parentErrStackSize, 1},
		{errors.WithDetails(noMsgErr), "", "% +-.1v", parentErrStackSize, 0},

		// errors.WithDetails and parent with custom %+v which is ignored
		{errors.WithDetails(&errorWithFormat{"foobar\nmore data"}), "foobar", "% +-.1v", currentStackSize, 1},
		{errors.WithDetails(&errorWithFormat{"foobar\nmore data\n"}), "foobar", "% +-.1v", currentStackSize, 1},

		// errors.WithMessage and parent without stack
		{errors.WithMessage(parentNoStackErr, "read error"), "read error: parent", "% +-.1v", currentStackSize, 1},
		{errors.WithMessage(parentNoStackErr, ""), "parent", "% +-.1v", currentStackSize, 1},
		{errors.WithMessage(parentNoStackErr, "read error\n"), "read error\nparent", "% +-.1v", currentStackSize, 2},
		{errors.WithMessage(noMsgNoStackErr, "read error"), "read error", "% +-.1v", currentStackSize, 1},
		{errors.WithMessage(noMsgNoStackErr, ""), "", "% +-.1v", currentStackSize, 0},
		{errors.WithMessage(io.EOF, "read error"), "read error: EOF", "% +-.1v", currentStackSize, 1},

		// errors.WithMessage twice
		{errors.WithMessage(errors.WithMessage(io.EOF, "read error"), "client error"), "client error: read error: EOF", "% +-.1v", currentStackSize, 1},

		// errors.WithMessage and parent with stack
		{errors.WithMessage(parentErr, "read error"), "read error: parent", "% +-.1v", parentErrStackSize, 1},
		{errors.WithMessage(parentErr, ""), "parent", "% +-.1v", parentErrStackSize, 1},
		{errors.WithMessage(parentErr, "read error\n"), "read error\nparent", "% +-.1v", parentErrStackSize, 2},
		{errors.WithMessage(noMsgErr, ""), "", "% +-.1v", parentErrStackSize, 0},
		{errors.WithMessage(noMsgErr, "read error"), "read error", "% +-.1v", parentErrStackSize, 1},

		// errors.WithMessage and parent with custom %+v which is ignored and no stack
		{errors.WithMessage(&errorWithFormat{"foobar\nmore data"}, "read error"), "read error: foobar", "% +-.1v", currentStackSize, 1},
		{errors.WithMessage(&errorWithFormat{"foobar\nmore data\n"}, "read error"), "read error: foobar", "% +-.1v", currentStackSize, 1},

		// errors.WithMessage and parent with custom %+v which is ignored and stack
		{errors.WithMessage(parentWithFormat1Err, "read error"), "read error: foobar", "% +-.1v", parentErrStackSize, 1},
		{errors.WithMessage(parentWithFormat2Err, "read error"), "read error: foobar", "% +-.1v", parentErrStackSize, 1},

		// errors.WithMessagef
		{errors.WithMessagef(parentNoStackErr, "read error %d", 1), "read error 1: parent", "% +-.1v", currentStackSize, 1},
		// We use w() to prevent static analysis.
		{errors.WithMessagef(parentNoStackErr, "read error ("+w()+")", noMsgNoStackErr), "read error (%!w(*errors.errorString=&{})): parent", "% +-.1v", currentStackSize, 1},

		// errors.Wrap and parent without stack, there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		{errors.Wrap(parentNoStackErr, "read error"), "read error", "% +-.1v", currentStackSize, 3 + 2},
		{errors.Wrap(parentNoStackErr, ""), "", "% +-.1v", currentStackSize, 3 + 1},
		{errors.Wrap(parentNoStackErr, "read error\n"), "read error\n", "% +-.1v", currentStackSize, 3 + 2},
		{errors.Wrap(io.EOF, "read error"), "read error", "% +-.1v", currentStackSize, 3 + 2},
		// There is no "the above error was caused by the following error" message.
		{errors.Wrap(noMsgNoStackErr, "read error"), "read error", "% +-.1v", currentStackSize, 3 + 1},
		{errors.Wrap(noMsgNoStackErr, ""), "", "% +-.1v", currentStackSize, 3 + 0},

		// errors.Wrap and parent with stack, there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		// + 1 for additional stack trace (most recent call first)" line
		{errors.Wrap(parentErr, "read error"), "read error", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 2 + 1},
		{errors.Wrap(parentErr, ""), "", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 1 + 1},
		{errors.Wrap(parentErr, "read error\n"), "read error\n", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 2 + 1},
		{errors.Wrap(noMsgErr, "read error"), "read error", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 1 + 1},
		{errors.Wrap(noMsgErr, ""), "", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 0 + 1},

		// errors.Wrap and parent with custom %+v and no stack, there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		{errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"), "read error", "% +-.1v", currentStackSize, 3 + 3},
		{errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"), "read error", "% +-.1v", currentStackSize, 3 + 3},

		// errors.Wrap and parent with custom %+v and no stack, there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		{errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"), "read error", "% +-.3v", currentStackSize, 3 + 3},
		{errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"), "read error", "% +-.3v", currentStackSize, 3 + 3},

		// errors.Wrap and parent with custom %+v (which is ignored) and stack there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		// + 1 for additional stack trace (most recent call first)" line
		{errors.Wrap(parentWithFormat1Err, "read error"), "read error", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 2 + 1},
		{errors.Wrap(parentWithFormat2Err, "read error"), "read error", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 2 + 1},

		// errors.Wrap and parent with custom %+v (which is ignored) and stack there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		// + 1 for additional stack trace (most recent call first)" line
		{errors.Wrap(parentWithFormat1Err, "read error"), "read error", "% +-.3v", currentStackSize + parentErrStackSize, 3 + 2 + 1},
		{errors.Wrap(parentWithFormat2Err, "read error"), "read error", "% +-.3v", currentStackSize + parentErrStackSize, 3 + 2 + 1},

		// errors.Wrapf
		{errors.Wrapf(parentNoStackErr, "read error %d", 1), "read error 1", "% +-.1v", currentStackSize, 3 + 2},
		// We use w() to prevent static analysis.
		{errors.Wrapf(parentNoStackErr, "read error ("+w()+")", noMsgNoStackErr), "read error (%!w(*errors.errorString=&{}))", "% +-.1v", currentStackSize, 3 + 2},

		// errors.Errorf with %w and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.Errorf("read error with parent: %w", parentPkgError), "read error with parent: parent", "% +-.1v", parentErrStackSize, 1},

		// errors.WithStack and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.WithStack(parentPkgError), "parent", "% +-.1v", parentErrStackSize, 1},

		// errors.WithStack and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.WithStack(parentPkgError), "parent", "% +-.1v", parentErrStackSize, 1},

		// errors.WithDetails and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.WithDetails(parentPkgError), "parent", "% +-.1v", parentErrStackSize, 1},

		// errors.WithDetails and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.WithDetails(parentPkgError), "parent", "% +-.1v", parentErrStackSize, 1},

		// errors.Wrap and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.Wrap(parentPkgError, "read error"), "read error", "% +-.1v", currentStackSize + parentErrStackSize, 3 + 3},

		// errors.Wrap and github.com/pkg/errors parent,
		// formatting of the cause is fully done by parentPkgError in this case,
		// there are still three lines extra for "The above error was caused by the
		// following error" + lines for error messages, but
		// there is no second "stack trace (most recent call first)" line,
		// a final newline is still added
		{errors.Wrap(parentPkgError, "read error"), "read error", "% +-.3v", currentStackSize + parentErrStackSize, 3 + 2},

		// errors.WithMessage and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.WithMessage(parentPkgError, "read error"), "read error: parent", "% +-.1v", parentErrStackSize, 1},

		// errors.WithMessage and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.WithMessage(parentPkgError, "read error"), "read error: parent", "% +-.3v", parentErrStackSize, 1},

		// Wrap behaves like New and Errorf if provided error is nil.
		{errors.Wrap(nil, "foo"), "foo", "% +-.1v", currentStackSize, 1},
		{errors.Wrap(nil, "read error without format specifiers"), "read error without format specifiers", "% +-.1v", currentStackSize, 1},

		// errors.Join.
		{errors.Join(errors.Base("foo1"), errors.Base("foo2")), "foo1\nfoo2", "% +-.1v", currentStackSize, 2 + 3 + 1 + 1 + 1},
		{errors.Join(errors.New("foo1"), errors.New("foo2")), "foo1\nfoo2", "% +-.1v", 3 * currentStackSize, 2 + 3 + 2 + 3},
	}...)
}

func TestErrors(t *testing.T) {
	t.Parallel()

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.EqualError(t, tt.err, tt.want)
			assert.Implements(t, (*stackTracer)(nil), tt.err)
			assert.Equal(t, tt.want, fmt.Sprintf("%s", tt.err))
			assert.Equal(t, tt.want, fmt.Sprintf("%v", tt.err))
			assert.Equal(t, fmt.Sprintf("%q", tt.want), fmt.Sprintf("%q", tt.err))
			stackTrace := fmt.Sprintf(tt.format, tt.err)
			// Expected stack size (2 lines per frame), plus "stack trace
			// (most recent call first)" line, plus extra lines.
			assert.Equal(t, tt.stack*2+1+tt.extra, strings.Count(stackTrace, "\n"), stackTrace)
		})
	}
}

func TestWithStackNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, errors.WithStack(nil), nil)
}

func TestWithDetailsNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, errors.WithDetails(nil), nil)
}

func TestWithMessageNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, errors.WithMessage(nil, "no error"), nil)
}

func TestWithMessagefNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, errors.WithMessagef(nil, "no error"), nil)
}

func TestJoinNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, errors.Join(nil))
	assert.Nil(t, errors.Join(nil, nil))
}

// stderrors.New, etc values are not expected to be compared by value.
// Assert that various kinds of errors have a functional equality operator,
// even if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	t.Parallel()

	vals := []error{
		nil,
		io.EOF,
		stderrors.New("EOF"),
		errors.New("EOF"),
		errors.Errorf("EOF"),
		errors.Wrap(io.EOF, "EOF"),
		errors.Wrapf(io.EOF, "EOF%d", 2),
		errors.WithMessage(nil, "whoops"),
		errors.WithMessage(io.EOF, "whoops"),
		errors.WithStack(io.EOF),
		errors.WithStack(nil),
	}

	for i := range vals {
		for j := range vals {
			// Must not panic.
			_ = vals[i] == vals[j] //nolint:errorlint
		}
	}
}

func TestBases(t *testing.T) {
	t.Parallel()

	// Current stack plus call to errors.*.
	currentStackSize := len(callers()) + 1

	grandparent := errors.Base("grandparent")
	parent := errors.BaseWrap(grandparent, "parent")
	err := errors.WithStack(parent)
	assert.EqualError(t, err, "parent")
	assert.Implements(t, (*stackTracer)(nil), err)
	stackTrace := fmt.Sprintf("% +-v", err)
	// Expected stack size (2 lines per frame), plus "stack trace
	// (most recent call first)" line, plus extra lines.
	assert.Equal(t, currentStackSize*2+1+1, strings.Count(stackTrace, "\n"), stackTrace)
	assert.ErrorIs(t, err, parent)
	assert.ErrorIs(t, err, grandparent)
}

func TestCause(t *testing.T) {
	t.Parallel()

	assert.Nil(t, errors.Cause(errors.Base("foo")))
	assert.Nil(t, errors.Cause(errors.New("foo")))
	assert.Nil(t, errors.Cause(errors.WithMessage(errors.Base("foo"), "bar")))

	err := errors.Base("foo")
	assert.Equal(t, err, errors.Cause(errors.Wrap(err, "bar")))
	assert.Equal(t, err, errors.Cause(errors.WithMessage(errors.Wrap(err, "bar"), "zar")))

	wrap := &errorWithCauseAndWrap{"test", nil, nil}
	assert.Nil(t, errors.Cause(wrap))

	wrap.wrap = err
	assert.Nil(t, errors.Cause(wrap))

	wrap.wrap = errors.Wrap(err, "bar")
	assert.Equal(t, err, errors.Cause(wrap))
}

func TestDetails(t *testing.T) {
	t.Parallel()

	err := errors.New("test")
	errors.Details(err)["zoo"] = "base"
	errors.Details(err)["foo"] = "bar"
	assert.Equal(t, map[string]interface{}{"zoo": "base", "foo": "bar"}, errors.Details(err))
	assert.Equal(t, map[string]interface{}{"zoo": "base", "foo": "bar"}, errors.AllDetails(err))
	err2 := errors.WithDetails(err)
	errors.Details(err2)["foo"] = "baz"
	errors.Details(err2)["foo2"] = "bar2"
	assert.Equal(t, map[string]interface{}{"zoo": "base", "foo": "bar"}, errors.Details(err))
	assert.Equal(t, map[string]interface{}{"zoo": "base", "foo": "bar"}, errors.AllDetails(err))
	assert.Equal(t, map[string]interface{}{"foo2": "bar2", "foo": "baz"}, errors.Details(err2))
	assert.Equal(t, map[string]interface{}{"foo2": "bar2", "foo": "baz", "zoo": "base"}, errors.AllDetails(err2))
	err3 := errors.WithDetails(err, "foo", "baz", "foo2", "bar2")
	assert.Equal(t, map[string]interface{}{"foo2": "bar2", "foo": "baz", "zoo": "base"}, errors.AllDetails(err3))
}
