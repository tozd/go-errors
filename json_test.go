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
func jsonEqual(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) bool {
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

func TestJSON(t *testing.T) {
	t.Parallel()

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
		`{"error":"error","stack":[],"cause":{"error":"bar"}}`,
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
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			jsonError, err := json.Marshal(tt.error)
			require.NoError(t, err)
			jsonEqual(t, tt.want, string(jsonError))
		})
	}
}
