//go:build go1.20

// We want this file to be after errors_test.go so that extra test cases are appended to others.
package errors_test

import "gitlab.com/tozd/go/errors"

func init() {
	currentStackSize := len(callers()) + 1

	tests = append(tests, []testStruct{
		// errors.Errorf with multiple %w without stack
		{errors.Errorf("%w, %w", errors.Base("foo1"), errors.Base("foo2")), "foo1, foo2", currentStackSize, 5},

		// errors.Errorf with multiple %w with stack
		{errors.Errorf("%w, %w", errors.New("foo1"), errors.New("foo2")), "foo1, foo2", 3 * currentStackSize, 7},
	}...)
}
