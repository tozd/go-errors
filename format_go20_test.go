//go:build go1.20
// +build go1.20

package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

func TestJoinErrorf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		error

		format string
		want   string
	}{{
		errors.Errorf("multiple: %w %w", errors.New("error1"), errors.New("error2")),
		"% #-+.1v",
		"^multiple: error1 error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoinErrorf\n" +
			"\t.+/format_go20_test.go:23\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoinErrorf\n" +
			"\t\t.+/format_go20_test.go:23\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoinErrorf\n" +
			"\t\t.+/format_go20_test.go:23\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, errors.Formatter{Error: tt.error}))
		})
	}
}
