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

func TestErrors(t *testing.T) {
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

	tests := []struct {
		err   error
		want  string
		stack int
		extra int
	}{
		// errors.New
		{errors.New(""), "", currentStackSize, 0},
		{errors.New("foo"), "foo", currentStackSize, 1},
		{errors.New("string with format specifiers: %v"), "string with format specifiers: %v", currentStackSize, 1},
		{errors.New("foo with newline\n"), "foo with newline\n", currentStackSize, 1},

		// errors.Errorf without %w
		{errors.Errorf(""), "", currentStackSize, 0},
		{errors.Errorf("read error without format specifiers"), "read error without format specifiers", currentStackSize, 1},
		{errors.Errorf("read error with %d format specifier", 1), "read error with 1 format specifier", currentStackSize, 1},
		{errors.Errorf("read error with newline\n"), "read error with newline\n", currentStackSize, 1},

		// errors.Errorf with %w and parent with stack
		{errors.Errorf("read error with parent: %w", parentErr), "read error with parent: parent", parentErrStackSize, 1},
		{errors.Errorf(`read error (parent "%w")`, parentErr), `read error (parent "parent")`, parentErrStackSize, 1},
		{errors.Errorf("read error (%w) and newline\n", parentErr), "read error (parent) and newline\n", parentErrStackSize, 1},
		{errors.Errorf("%w", noMsgErr), "", parentErrStackSize, 0},

		// errors.Errorf with %w and parent without stack
		{errors.Errorf("read error with parent: %w", parentNoStackErr), "read error with parent: parent", currentStackSize, 1},
		{errors.Errorf(`read error (parent "%w")`, parentNoStackErr), `read error (parent "parent")`, currentStackSize, 1},
		{errors.Errorf("read error (%w) and newline\n", parentNoStackErr), "read error (parent) and newline\n", currentStackSize, 1},
		{errors.Errorf("%w", noMsgNoStackErr), "", currentStackSize, 0},

		// errors.WithStack and parent without stack
		{errors.WithStack(io.EOF), "EOF", currentStackSize, 1},
		{errors.WithStack(errors.Base("EOF")), "EOF", currentStackSize, 1},
		{errors.WithStack(errors.Base("")), "", currentStackSize, 0},
		{errors.WithStack(errors.Base("foobar\n")), "foobar\n", currentStackSize, 1},

		// errors.WithStack and parent with stack
		{errors.WithStack(parentErr), "parent", parentErrStackSize, 1},
		{errors.WithStack(noMsgErr), "", parentErrStackSize, 0},

		// errors.WithStack and parent with custom %+v
		{errors.WithStack(&errorWithFormat{"foobar\nmore data"}), "foobar", currentStackSize, 2},
		{errors.WithStack(&errorWithFormat{"foobar\nmore data\n"}), "foobar", currentStackSize, 2},

		// errors.WithDetails and parent without stack
		{errors.WithDetails(io.EOF), "EOF", currentStackSize, 1},
		{errors.WithDetails(errors.Base("EOF")), "EOF", currentStackSize, 1},
		{errors.WithDetails(errors.Base("")), "", currentStackSize, 0},
		{errors.WithDetails(errors.Base("foobar\n")), "foobar\n", currentStackSize, 1},

		// errors.WithDetails and parent with stack
		{errors.WithDetails(parentErr), "parent", parentErrStackSize, 1},
		{errors.WithDetails(noMsgErr), "", parentErrStackSize, 0},

		// errors.WithDetails and parent with custom %+v
		{errors.WithDetails(&errorWithFormat{"foobar\nmore data"}), "foobar", currentStackSize, 2},
		{errors.WithDetails(&errorWithFormat{"foobar\nmore data\n"}), "foobar", currentStackSize, 2},

		// errors.WithMessage and parent without stack
		{errors.WithMessage(parentNoStackErr, "read error"), "read error: parent", currentStackSize, 1},
		{errors.WithMessage(parentNoStackErr, ""), "parent", currentStackSize, 1},
		{errors.WithMessage(parentNoStackErr, "read error\n"), "read error\nparent", currentStackSize, 2},
		{errors.WithMessage(noMsgNoStackErr, "read error"), "read error", currentStackSize, 1},
		{errors.WithMessage(noMsgNoStackErr, ""), "", currentStackSize, 0},
		{errors.WithMessage(io.EOF, "read error"), "read error: EOF", currentStackSize, 1},

		// errors.WithMessage twice
		{errors.WithMessage(errors.WithMessage(io.EOF, "read error"), "client error"), "client error: read error: EOF", currentStackSize, 1},

		// errors.WithMessage and parent with stack
		{errors.WithMessage(parentErr, "read error"), "read error: parent", parentErrStackSize, 1},
		{errors.WithMessage(parentErr, ""), "parent", parentErrStackSize, 1},
		{errors.WithMessage(parentErr, "read error\n"), "read error\nparent", parentErrStackSize, 2},
		{errors.WithMessage(noMsgErr, ""), "", parentErrStackSize, 0},
		// "read error" is prefixed to the parent's stack trace output, so there is no extra line for the error message
		{errors.WithMessage(noMsgErr, "read error"), "read error", parentErrStackSize, 0},

		// errors.WithMessage and parent with custom %+v and no stack
		{errors.WithMessage(&errorWithFormat{"foobar\nmore data"}, "read error"), "read error: foobar", currentStackSize, 2},
		{errors.WithMessage(&errorWithFormat{"foobar\nmore data\n"}, "read error"), "read error: foobar", currentStackSize, 2},

		// errors.WithMessage and parent with custom %+v and stack
		{errors.WithMessage(parentWithFormat1Err, "read error"), "read error: foobar", parentErrStackSize, 2},
		{errors.WithMessage(parentWithFormat2Err, "read error"), "read error: foobar", parentErrStackSize, 2},

		// errors.WithMessagef
		{errors.WithMessagef(parentNoStackErr, "read error %d", 1), "read error 1: parent", currentStackSize, 1},
		// We use w() to prevent static analysis.
		{errors.WithMessagef(parentNoStackErr, "read error ("+w()+")", noMsgNoStackErr), "read error (%!w(*errors.errorString=&{})): parent", currentStackSize, 1},

		// errors.Wrap and parent without stack, there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		{errors.Wrap(parentNoStackErr, "read error"), "read error", currentStackSize, 3 + 2},
		{errors.Wrap(parentNoStackErr, ""), "", currentStackSize, 3 + 1},
		{errors.Wrap(parentNoStackErr, "read error\n"), "read error\n", currentStackSize, 3 + 2},
		{errors.Wrap(io.EOF, "read error"), "read error", currentStackSize, 3 + 2},
		// There is no "the above error was caused by the following error" message.
		{errors.Wrap(noMsgNoStackErr, "read error"), "read error", currentStackSize, 1},
		{errors.Wrap(noMsgNoStackErr, ""), "", currentStackSize, 0},

		// errors.Wrap and parent with stack, there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		// + 1 for additional stack trace (most recent call first)" line
		{errors.Wrap(parentErr, "read error"), "read error", currentStackSize + parentErrStackSize, 3 + 2 + 1},
		{errors.Wrap(parentErr, ""), "", currentStackSize + parentErrStackSize, 3 + 1 + 1},
		{errors.Wrap(parentErr, "read error\n"), "read error\n", currentStackSize + parentErrStackSize, 3 + 2 + 1},
		{errors.Wrap(noMsgErr, "read error"), "read error", currentStackSize + parentErrStackSize, 3 + 1 + 1},
		{errors.Wrap(noMsgErr, ""), "", currentStackSize + parentErrStackSize, 3 + 0 + 1},

		// errors.Wrap and parent with custom %+v and no stack, there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		{errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"), "read error", currentStackSize, 3 + 3},
		{errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"), "read error", currentStackSize, 3 + 3},

		// errors.Wrap and parent with custom %+v and stack there are three lines extra for
		// "the above error was caused by the following error" + lines for error messages
		// + 1 for additional stack trace (most recent call first)" line
		{errors.Wrap(parentWithFormat1Err, "read error"), "read error", currentStackSize + parentErrStackSize, 3 + 3 + 1},
		{errors.Wrap(parentWithFormat2Err, "read error"), "read error", currentStackSize + parentErrStackSize, 3 + 3 + 1},

		// errors.Wrapf
		{errors.Wrapf(parentNoStackErr, "read error %d", 1), "read error 1", currentStackSize, 3 + 2},
		// We use w() to prevent static analysis.
		{errors.Wrapf(parentNoStackErr, "read error ("+w()+")", noMsgNoStackErr), "read error (%!w(*errors.errorString=&{}))", currentStackSize, 3 + 2},

		// errors.Errorf with %w and github.com/pkg/errors parent,
		// we format the stack trace in this case
		{errors.Errorf("read error with parent: %w", parentPkgError), "read error with parent: parent", parentErrStackSize, 1},

		// errors.WithStack and github.com/pkg/errors parent,
		// formatting is fully done by parentPkgError in this case,
		// there is still one line for the error message, but
		// there is no "stack trace (most recent call first)" line,
		// and no final newline
		{errors.WithStack(parentPkgError), "parent", parentErrStackSize, 1 - 1 - 1},

		// errors.WithDetails and github.com/pkg/errors parent,
		// formatting is fully done by parentPkgError in this case,
		// there is still one line for the error message, but
		// there is no "stack trace (most recent call first)" line,
		// and no final newline
		{errors.WithDetails(parentPkgError), "parent", parentErrStackSize, 1 - 1 - 1},

		// errors.Wrap and github.com/pkg/errors parent,
		// formatting of the cause is fully done by parentPkgError in this case,
		// there are still three lines extra for "The above error was caused by the
		// following error" + lines for error messages, but
		// there is no second "stack trace (most recent call first)" line,
		// a final newline is still added
		{errors.Wrap(parentPkgError, "read error"), "read error", currentStackSize + parentErrStackSize, 3 + 2},

		// errors.WithMessage and github.com/pkg/errors parent,
		// formatting of the cause is fully done by parentPkgError in this case,
		// additional message is just prefixed, there is still one line for the
		// error message, but there is no "stack trace (most recent call first)" line,
		// and no final newline
		{errors.WithMessage(parentPkgError, "read error"), "read error: parent", parentErrStackSize, 1 - 1 - 1},

		// Wrap behaves like New and Errorf if provided error is nil.
		{errors.Wrap(nil, "foo"), "foo", currentStackSize, 1},
		{errors.Wrap(nil, "read error without format specifiers"), "read error without format specifiers", currentStackSize, 1},

		// errors.Join.
		{errors.Join(errors.Base("foo1"), errors.Base("foo2")), "foo1\nfoo2", currentStackSize, 4},
		{errors.Join(errors.New("foo1"), errors.New("foo2")), "foo1\nfoo2", 3 * currentStackSize, 6},
	}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.EqualError(t, tt.err, tt.want)
			assert.Implements(t, (*stackTracer)(nil), tt.err)
			assert.Equal(t, tt.want, fmt.Sprintf("%s", tt.err))
			assert.Equal(t, tt.want, fmt.Sprintf("%v", tt.err))
			assert.Equal(t, fmt.Sprintf("%q", tt.want), fmt.Sprintf("%q", tt.err))
			stackTrace := fmt.Sprintf("%+v", tt.err)
			// Expected stack size (2 lines per frame), plus "Stack trace
			// (most recent call first)" line, plus extra lines.
			assert.Equal(t, tt.stack*2+1+tt.extra, strings.Count(stackTrace, "\n"), stackTrace)
		})
	}
}

func TestWithStackNil(t *testing.T) {
	assert.Nil(t, errors.WithStack(nil), nil)
}

func TestWithDetailsNil(t *testing.T) {
	assert.Nil(t, errors.WithDetails(nil), nil)
}

func TestWithMessageNil(t *testing.T) {
	assert.Nil(t, errors.WithMessage(nil, "no error"), nil)
}

func TestWithMessagefNil(t *testing.T) {
	assert.Nil(t, errors.WithMessagef(nil, "no error"), nil)
}

func TestJoinNil(t *testing.T) {
	assert.Nil(t, errors.Join(nil))
	assert.Nil(t, errors.Join(nil, nil))
}

// stderrors.New, etc values are not expected to be compared by value.
// Assert that various kinds of errors have a functional equality operator,
// even if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
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
			_ = vals[i] == vals[j]
		}
	}
}

func TestBases(t *testing.T) {
	// Current stack plus call to errors.*.
	currentStackSize := len(callers()) + 1

	grandparent := errors.Base("grandparent")
	parent := errors.BaseWrap(grandparent, "parent")
	err := errors.WithStack(parent)
	assert.EqualError(t, err, "parent")
	assert.Implements(t, (*stackTracer)(nil), err)
	stackTrace := fmt.Sprintf("%+v", err)
	// Expected stack size (2 lines per frame), plus "Stack trace
	// (most recent call first)" line, plus extra lines.
	assert.Equal(t, currentStackSize*2+1+1, strings.Count(stackTrace, "\n"), stackTrace)
	assert.ErrorIs(t, err, parent)
	assert.ErrorIs(t, err, grandparent)
}

func TestCause(t *testing.T) {
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
