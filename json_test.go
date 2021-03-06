package errors_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/tozd/go/errors"
)

// See: https://github.com/stretchr/testify/issues/535
func equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	// We do not compare JSON arrays (which are used to represent stack traces).
	opt := cmp.Transformer("stack", func(_ []interface{}) []interface{} {
		return []interface{}{}
	})

	if !cmp.Equal(expected, actual, opt) {
		diff := cmp.Diff(expected, actual, opt)
		return assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s%s", expected, actual, diff), msgAndArgs...)
	}
	return true
}

// jsonEqual ignores and does not compare JSON arrays (which are used to represent stack traces).
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
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			jsonError, err := json.Marshal(tt.error)
			require.NoError(t, err)
			jsonEqual(t, tt.want, string(jsonError))
		})
	}
}
