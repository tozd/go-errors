package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.New("error"),
		"%s",
		"^error$",
	}, {
		errors.New("error"),
		"%v",
		"^error$",
	}, {
		errors.New("error"),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:28\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%q",
		`^"error"$`,
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatErrorf(t *testing.T) {
	parentErr := errors.New("error")
	parentNoStackErr := stderrors.New("error")

	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Errorf("%s", "error"),
		"%s",
		"^error$",
	}, {
		errors.Errorf("%s", "error"),
		"%v",
		"^error$",
	}, {
		errors.Errorf("%s", "error"),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:65\n",
	}, {
		errors.Errorf("%w", parentErr),
		"%s",
		"^error$",
	}, {
		errors.Errorf("%w", parentErr),
		"%v",
		"^error$",
	}, {
		errors.Errorf("%w", parentErr),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:49\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentNoStackErr),
		"%s",
		"^error$",
	}, {
		errors.Errorf("%w", parentNoStackErr),
		"%v",
		"^error$",
	}, {
		errors.Errorf("%w", parentNoStackErr),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:96\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWithStack(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.WithStack(io.EOF),
		"%s",
		"^EOF$",
	}, {
		errors.WithStack(io.EOF),
		"%v",
		"^EOF$",
	}, {
		errors.WithStack(io.EOF),
		"%+v",
		"^EOF\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:126\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.New("error"),
		),
		"%s",
		"^error$",
	}, {
		errors.WithStack(
			errors.New("error"),
		),
		"%v",
		"^error$",
	}, {
		errors.WithStack(
			errors.New("error"),
		),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:147\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.Base("error"),
		),
		"%s",
		"^error$",
	}, {
		errors.WithStack(
			errors.Base("error"),
		),
		"%v",
		"^error$",
	}, {
		errors.WithStack(
			errors.Base("error"),
		),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:168\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.WithStack(io.EOF),
		),
		"%+v",
		"^EOF\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:179\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.Errorf("error%d", 1),
		),
		"%+v",
		"error1\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:189\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%s",
		"^error2$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%v",
		"^error2$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%+v",
		"^error2\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:226\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:227\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%s",
		"error",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%v",
		"error",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:251\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.Wrap(
			errors.Wrap(io.EOF, "error1"),
			"error2",
		),
		"%+v",
		"^error2\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:261\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"error1\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:262\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.Wrap(
			errors.New("error with space"),
			"context",
		),
		"%q",
		`"context"`,
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Wrapf(io.EOF, "error%d", 2),
		"%s",
		"^error2$",
	}, {
		errors.Wrapf(io.EOF, "error%d", 2),
		"%v",
		"^error2$",
	}, {
		errors.Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"^error2\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:309\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {}, {
		errors.Wrapf(
			errors.New("error"),
			"error%d", 2,
		),
		"%s",
		"^error2$",
	}, {
		errors.Wrapf(
			errors.New("error"),
			"error%d", 2,
		),
		"%v",
		"^error2$",
	}, {
		errors.Wrapf(
			errors.New("error"),
			"error%d", 2,
		),
		"%+v",
		"^error2\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:333\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:334\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWithMessage(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.WithMessage(
			errors.New("error"), "error2",
		),
		"%s",
		"^error2: error$",
	}, {
		errors.WithMessage(
			errors.New("error"), "error2",
		),
		"%v",
		"^error2: error$",
	}, {
		errors.WithMessage(
			errors.New("error"), "error2",
		),
		"%+v",
		"^error2: error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:377\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(io.EOF, "addition1"),
		"%s",
		"^addition1: EOF$",
	}, {
		errors.WithMessage(io.EOF, "addition1"),
		"%v",
		"^addition1: EOF$",
	}, {
		errors.WithMessage(io.EOF, "addition1"),
		"%+v",
		"^addition1: EOF\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:394\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.WithMessage(io.EOF, "addition1"),
			"addition2",
		),
		"%v",
		"^addition2: addition1: EOF$",
	}, {
		errors.WithMessage(
			errors.WithMessage(io.EOF, "addition1"),
			"addition2",
		),
		"%+v",
		"^addition2: addition1: EOF\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:410\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.WithMessage(io.EOF, "error1"),
			"error2",
		),
		"%+v",
		"^error2\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:420\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"error1: EOF\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:421\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.Errorf("error%d", 1),
			"error2",
		),
		"%+v",
		"^error2: error1\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:438\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.WithStack(io.EOF),
			"error",
		),
		"%+v",
		"^error: EOF\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:449\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.Wrap(
				errors.WithStack(io.EOF),
				"inside-error",
			),
			"outside-error",
		),
		"%+v",
		"^outside-error: inside-error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:460\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nThe above error was caused by the following error:\n\n" +
			"EOF\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:461\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

// This is mid-stack inlined in go 1.12+.
func wrappedNew(message string) error {
	return errors.New(message)
}

func TestFormatWrappedNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		wrappedNew("error"),
		"%+v",
		"^error\n" +
			"Stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.wrappedNew\n" +
			"\t.+/format_test.go:489\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrappedNew\n" +
			"\t.+/format_test.go:498\n",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}
