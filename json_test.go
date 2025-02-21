package errors_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/tozd/go/errors"
)

// See: https://github.com/stretchr/testify/issues/535
func equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	// We do not compare stack frames.
	opt := cmpopts.IgnoreSliceElements(func(el interface{}) bool {
		v, ok := el.(map[string]interface{})
		if !ok {
			return false
		}
		_, name := v["name"]
		_, file := v["file"]
		_, line := v["line"]
		return name && file && line
	})

	if !cmp.Equal(expected, actual, opt) {
		diff := cmp.Diff(expected, actual, opt)
		return assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s%s", expected, actual, diff), msgAndArgs...)
	}
	return true
}

// jsonEqual ignores and does not compare frames.
func jsonEqual(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) bool { //nolint:unparam
	t.Helper()
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		return assert.Fail(t, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()), msgAndArgs...)
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		return assert.Fail(t, fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()), msgAndArgs...)
	}

	return equal(t, expectedJSONAsInterface, actualJSONAsInterface, msgAndArgs...)
}

type testValueReceiverError struct{}

func (e testValueReceiverError) Error() string {
	return "error"
}

func TestJSON(t *testing.T) {
	t.Parallel()

	testErr := &testStructJoined{msg: "test2"}

	tests := []struct {
		error
		want string
	}{{
		errors.New("error"),
		`{"error":"error","stack":[]}`,
	}, {
		errors.Errorf("error: %w", errors.Base("foobar")),
		`{"error":"error: foobar","stack":[]}`,
	}, {
		errors.Errorf("error: %w", errors.New("foobar")),
		`{"error":"error: foobar","stack":[]}`,
	}, {
		errors.Errorf("error: %w", pkgerrors.New("foobar")),
		`{"error":"error: foobar","stack":[]}`,
	}, {
		errors.WithStack(errors.Base("error")),
		`{"error":"error","stack":[]}`,
	}, {
		errors.WithStack(pkgerrors.New("error")),
		`{"error":"error","stack":[]}`,
	}, {
		errors.WithMessage(errors.Base("foobar"), "error"),
		`{"error":"error: foobar","stack":[]}`,
	}, {
		errors.WithMessage(errors.New("foobar"), "error"),
		`{"error":"error: foobar","stack":[]}`,
	}, {
		errors.WithMessage(pkgerrors.New("foobar"), "error"),
		`{"error":"error: foobar","stack":[]}`,
	}, {
		errors.Wrap(errors.Base("foobar"), "error"),
		`{"error":"error","stack":[],"cause":{"error":"foobar"}}`,
	}, {
		errors.Wrap(errors.New("foobar"), "error"),
		`{"error":"error","stack":[],"cause":{"error":"foobar","stack":[]}}`,
	}, {
		errors.Wrap(errors.BaseWrap(errors.Base("foo"), "bar"), "error"),
		`{"error":"error","stack":[],"cause":{"error":"bar"}}`,
	}, {
		errors.Wrap(errors.BaseWrap(errors.New("foo"), "bar"), "error"),
		`{"error":"error","stack":[],"cause":{"error":"bar","stack":[]}}`,
	}, {
		errors.Wrap(pkgerrors.New("foobar"), "error"),
		`{"error":"error","stack":[],"cause":{"error":"foobar","stack":[]}}`,
	}, {
		errors.Wrap(pkgerrors.Wrap(pkgerrors.New("foo"), "bar"), "error"),
		`{"error":"error","stack":[],"cause":{"error":"bar: foo","stack":[],"cause":{"error":"bar: foo","cause":{"error":"foo","stack":[]}}}}`,
	}, {
		errors.Wrap(pkgerrors.WithMessage(errors.Base("foo"), "bar"), "error"),
		`{"error":"error","stack":[],"cause":{"error":"bar: foo","cause":{"error":"foo"}}}`,
	}, {
		errors.Join(errors.Base("foobar1"), errors.Base("foobar2")),
		`{"error":"foobar1\nfoobar2","errors":[{"error":"foobar1"},{"error":"foobar2"}],"stack":[]}`,
	}, {
		errors.Join(errors.New("foobar1"), errors.New("foobar2")),
		`{"error":"foobar1\nfoobar2","errors":[{"error":"foobar1","stack":[]},{"error":"foobar2","stack":[]}],"stack":[]}`,
	}, {
		errors.WithDetails(errors.Base("error"), "foo", "bar"),
		`{"error":"error","foo":"bar","stack":[]}`,
	}, {
		errors.WithDetails(errors.Base("error"), "stack", "foobar"),
		`{"error":"error","stack":[]}`,
	}, {
		errors.WithDetails(errors.Join(errors.WithDetails(errors.New("foobar1"), "foo", 1), errors.WithDetails(errors.New("foobar2"), "foo", 2)), "foo", "bar"),
		`{"error":"foobar1\nfoobar2","errors":[{"error":"foobar1","foo":1,"stack":[]},{"error":"foobar2","foo":2,"stack":[]}],"foo":"bar","stack":[]}`,
	}, {
		errors.WithDetails(errors.Wrap(errors.WithDetails(errors.New("foobar"), "foo", 1), "error"), "foo", 2),
		`{"error":"error","foo":2,"stack":[],"cause":{"error":"foobar","foo":1,"stack":[]}}`,
	}, {
		errors.WrapWith(errors.Base("foobar"), errors.Base("error")),
		`{"error":"error","stack":[],"cause":{"error":"foobar"}}`,
	}, {
		errors.WrapWith(errors.Base("foobar"), errors.WrapWith(errors.Base("error1"), errors.Base("error2"))),
		`{"cause":{"error":"foobar"},"error":"error2","errors":[{"cause":{"error":"error1"},"error":"error2","stack":[]}],"stack":[]}`,
	}, {
		errors.WrapWith(errors.Base("foobar"), errors.WithDetails(errors.Base("error1"), "x", "foo")),
		`{"cause":{"error":"foobar"},"error":"error1","errors":[{"error":"error1","stack":[],"x":"foo"}],"stack":[]}`,
	}, {
		errors.WrapWith(errors.Base("foobar"), errors.WithDetails(errors.Base("error"))),
		`{"error":"error","stack":[],"cause":{"error":"foobar"}}`,
	}, {
		errors.WrapWith(errors.Base("foobar"), errors.WithStack(errors.Base("error"))),
		`{"error":"error","stack":[],"cause":{"error":"foobar"}}`,
	}, {
		&testStructJoined{},
		`{}`,
	}, {
		&testStructJoined{msg: "test"},
		`{"error":"test"}`,
	}, {
		&testStructJoined{msg: "test1", cause: testErr},
		`{"cause":{"error":"test2"},"error":"test1"}`,
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{testErr}},
		`{"cause":{"error":"test2"},"error":"test1"}`,
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{testErr, &testStructJoined{msg: "test3"}}},
		`{"cause":{"error":"test2"},"error":"test1","errors":[{"error":"test3"}]}`,
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{testErr, &testStructJoined{msg: "test3"}, &testStructJoined{msg: "test4"}}},
		`{"cause":{"error":"test2"},"error":"test1","errors":[{"error":"test3"},{"error":"test4"}]}`,
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{&testStructJoined{msg: "test3"}, &testStructJoined{msg: "test4"}}},
		`{"cause":{"error":"test2"},"error":"test1","errors":[{"error":"test3"},{"error":"test4"}]}`,
	}, {
		errors.Prefix(errors.New("error"), errors.Base("error2")),
		`{"error":"error2: error","errors":[{"error":"error2"},{"error":"error","stack":[]}],"stack":[]}`,
	}, {
		errors.Prefix(errors.Base("parent"), errors.Base("")),
		`{"error":"parent","stack":[]}`,
	}, {
		testValueReceiverError{},
		`{"error":"error"}`,
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			jsonError, err := json.Marshal(errors.Formatter{Error: tt.error})
			require.NoError(t, err)
			jsonEqual(t, tt.want, string(jsonError))

			err2, errE := errors.UnmarshalJSON(jsonError)
			require.NoError(t, errE, "% -+#.1v", errE)
			jsonError2, err := json.Marshal(err2)
			require.NoError(t, err)
			assert.Equal(t, string(jsonError), string(jsonError2)) //nolint:testifylint
		})
	}
}
