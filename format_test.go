package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

func TestFormatNew(t *testing.T) {
	t.Parallel()

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
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:31\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:38\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"  .+/format_test.go:46\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:54\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%q",
		`^"error"$`,
	}, {
		errors.New("error"),
		"% v",
		"^error\n$",
	}, {
		errors.New("error\n"),
		"% v",
		"^error\n$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatErrorf(t *testing.T) {
	t.Parallel()

	parentErr := errors.New("error")
	parentNoStackErr := stderrors.New("error")
	parentPkgErr := pkgerrors.New("error")

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
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:106\n",
	}, {
		errors.Errorf("%s", "error"),
		"%-+v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:112\n",
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
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:89\n" +
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
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:143\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%s",
		"^error$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%v",
		"^error$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:91\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:91\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"  .+/format_test.go:91\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:91\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWithStack(t *testing.T) {
	t.Parallel()

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
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:218\n" +
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
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:239\n" +
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
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:260\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.WithStack(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:271\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:281\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			pkgerrors.New("error"),
		),
		"%s",
		"^error$",
	}, {
		errors.WithStack(
			pkgerrors.New("error"),
		),
		"%v",
		"^error$",
	}, {
		errors.WithStack(
			pkgerrors.New("error"),
		),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:303\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			pkgerrors.New("error"),
		),
		"% +-#2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"  .+/format_test.go:313\n" +
			"(.+\n  .+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWrap(t *testing.T) {
	t.Parallel()

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
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:356\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:357\n" +
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
		"% +-.1v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:381\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.Wrap(
			errors.Wrap(io.EOF, "error1"),
			"error2",
		),
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:391\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:392\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.Wrap(
			errors.New("error with space"),
			"context",
		),
		"%q",
		`"context"`,
	}, {
		errors.Wrap(
			pkgerrors.New("error"),
			"error2",
		),
		"%s",
		"^error2$",
	}, {
		errors.Wrap(
			pkgerrors.New("error"),
			"error2",
		),
		"%v",
		"^error2$",
	}, {
		errors.Wrap(
			pkgerrors.New("error"),
			"error2",
		),
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:431\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:432\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			pkgerrors.New("error"),
			"error2",
		),
		"% +-.3v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:448\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:449\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			pkgerrors.New("error"),
			"error2",
		),
		"% +-2.3v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"  .+/format_test.go:464\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:465\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%+-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:480\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error was caused by the following error:\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:481\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%+.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:497\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:498\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"% +.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:511\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:512\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"% +-#v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:526\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"% -.1v",
		"^error2\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%-.1v",
		"^error2\n" +
			"the above error was caused by the following error:\n" +
			"error\n$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%.1v",
		"^error2\n" +
			"error\n$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWrapf(t *testing.T) {
	t.Parallel()

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
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:591\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
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
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:615\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:616\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWithMessage(t *testing.T) {
	t.Parallel()

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
		"% +-.1v",
		"^error2: error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:665\n" +
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
		"% +-.1v",
		"^addition1: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:682\n" +
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
		"% +-.1v",
		"^addition2: addition1: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:698\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.WithMessage(io.EOF, "error1"),
			"error2",
		),
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:708\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:709\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.Errorf("error%d", 1),
			"error2",
		),
		"% +-.1v",
		"^error2: error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:726\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.WithStack(io.EOF),
			"error",
		),
		"% +-.1v",
		"^error: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:737\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.Wrap(
				errors.WithStack(io.EOF),
				"inside-error",
			),
			"outside-error",
		),
		"% +-.1v",
		"^outside-error: inside-error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:748\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:749\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			pkgerrors.New("error"), "error2",
		),
		"%s",
		"^error2: error$",
	}, {
		errors.WithMessage(
			pkgerrors.New("error"), "error2",
		),
		"%v",
		"^error2: error$",
	}, {
		errors.WithMessage(
			pkgerrors.New("error"), "error2",
		),
		"% +-.1v",
		"^error2: error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:780\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.Wrap(
				errors.WithStack(io.EOF),
				"inside-error",
			),
			"outside-error",
		),
		"% +-.v",
		"^outside-error: inside-error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:790\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.Wrap(
				errors.WithStack(io.EOF),
				"inside-error",
			),
			"outside-error",
		),
		"%+.1v",
		"^outside-error: inside-error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:804\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"EOF\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:805\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(
			errors.Wrap(
				errors.WithStack(io.EOF),
				"inside-error",
			),
			"outside-error",
		),
		"% -.1v",
		"^outside-error: inside-error\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.WithMessage(
			errors.Wrap(
				errors.WithStack(io.EOF),
				"inside-error",
			),
			"outside-error",
		),
		"%.1v",
		"^outside-error: inside-error\n" +
			"EOF\n$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

// This is mid-stack inlined in go 1.12+.
func wrappedNew(message string) error {
	return errors.New(message)
}

func TestFormatWrappedNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		error
		format string
		want   string
	}{{
		wrappedNew("error"),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.wrappedNew\n" +
			"\t.+/format_test.go:857\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrappedNew\n" +
			"\t.+/format_test.go:868\n",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWithDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.WithDetails(io.EOF),
		"%s",
		"^EOF$",
	}, {
		errors.WithDetails(io.EOF),
		"%v",
		"^EOF$",
	}, {
		errors.WithDetails(io.EOF),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:905\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
		),
		"%s",
		"^error$",
	}, {
		errors.WithDetails(
			errors.New("error"),
		),
		"%v",
		"^error$",
	}, {
		errors.WithDetails(
			errors.New("error"),
		),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:926\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.Base("error"),
		),
		"%s",
		"^error$",
	}, {
		errors.WithDetails(
			errors.Base("error"),
		),
		"%v",
		"^error$",
	}, {
		errors.WithDetails(
			errors.Base("error"),
		),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:947\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.WithDetails(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:958\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:968\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			pkgerrors.New("error"),
		),
		"%s",
		"^error$",
	}, {
		errors.WithDetails(
			pkgerrors.New("error"),
		),
		"%v",
		"^error$",
	}, {
		errors.WithDetails(
			pkgerrors.New("error"),
		),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:990\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
		),
		"%+#v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:1000\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
			"foo", 1,
			"bar", "baz",
			"quote", "one\ntwo",
		),
		"%+#v",
		"^error\n" +
			"bar=baz\n" +
			"foo=1\n" +
			"quote=\"one\\\\ntwo\"\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:1009\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
			"foo", 1,
			"bar", "baz",
			"quote", "one\ntwo",
		),
		"%#v",
		"^error\n" +
			"bar=baz\n" +
			"foo=1\n" +
			"quote=\"one\\\\ntwo\"\n$",
	}, {
		errors.WithDetails(
			errors.New("error"),
			"foo", 1,
			"bar", "baz",
			"quote", "one\ntwo",
		),
		"% #-v",
		"^error\n" +
			"bar=baz\n" +
			"foo=1\n" +
			"quote=\"one\\\\ntwo\"\n$",
	}, {
		errors.WithDetails(
			errors.New("error"),
			"foo", 1,
			"bar", "baz",
			"quote", "one\ntwo",
		),
		"% #-+v",
		"^error\n" +
			"bar=baz\n" +
			"foo=1\n" +
			"quote=\"one\\\\ntwo\"\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:1048\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
			"foo", 1,
			"bar", "baz",
			"quote", "one\ntwo",
		),
		"%#-+v",
		"^error\n" +
			"bar=baz\n" +
			"foo=1\n" +
			"quote=\"one\\\\ntwo\"\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:1064\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		error
		format string
		want   string
	}{{
		&errorWithFormat{"foobar\nmore data"},
		"% #-+.3v",
		"^foobar\nmore data\n$",
	}, {
		&errorWithFormat{"foobar\nmore data"},
		"% #-+.1v",
		"^foobar\nmore data\n$",
	}, {
		&errorWithFormat{"foobar\nmore data\n"},
		"% #-+.3v",
		"^foobar\nmore data\n$",
	}, {
		&errorWithFormat{"foobar\nmore data\n"},
		"% #-+.1v",
		"^foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1115\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1125\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1135\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1145\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% #-+.1v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% #-+.3v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"%v",
		"^test$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"%.2v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% .2v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, &errorWithFormat{"foobar\nmore data"}},
		"% #-+.1v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, &errorWithFormat{"foobar\nmore data"}},
		"% #-+.3v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormat{"foobar\nmore data"}, nil},
		"% #-+.1v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormat{"foobar\nmore data"}, nil},
		"% #-+.3v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormat{"foobar\nmore data\n"}, nil},
		"% #-+.1v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormat{"foobar\nmore data\n"}, nil},
		"% #-+.3v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		&errorWithFormatAndStack{"foobar\nmore data"},
		"% #-+.3v",
		"^foobar\nmore data\n$",
	}, {
		&errorWithFormatAndStack{"foobar\nmore data"},
		"% #-+.1v",
		"^foobar\n$",
	}, {
		&errorWithFormatAndStack{"foobar\nmore data\n"},
		"% #-+.3v",
		"^foobar\nmore data\n$",
	}, {
		&errorWithFormatAndStack{"foobar\nmore data\n"},
		"% #-+.1v",
		"^foobar\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1227\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1237\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1247\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1257\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% #-+.1v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% #-+.3v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"%v",
		"^test$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"%.2v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, nil},
		"% .2v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, &errorWithFormatAndStack{"foobar\nmore data"}},
		"% #-+.1v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", nil, &errorWithFormatAndStack{"foobar\nmore data"}},
		"% #-+.3v",
		"^test\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormatAndStack{"foobar\nmore data"}, nil},
		"% #-+.1v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormatAndStack{"foobar\nmore data"}, nil},
		"% #-+.3v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormatAndStack{"foobar\nmore data\n"}, nil},
		"% #-+.1v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\n$",
	}, {
		&errorWithCauseAndWrap{"test", &errorWithFormatAndStack{"foobar\nmore data\n"}, nil},
		"% #-+.3v",
		"^test\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, errors.Formatter{tt.error}))
		})
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Join(),
		"% #-+.3v",
		"^<nil>\n$",
	}, {
		errors.Join(),
		"%s",
		"^%!s\\(<nil>\\)$",
	}, {
		errors.Join(),
		"%q",
		"^%!q\\(<nil>\\)$",
	}, {
		errors.Join(errors.New("error")),
		"% #-+.1v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1355\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error"), nil),
		"% #-+.1v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1363\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1371\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1379\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1379\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1379\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+2.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  .+/format_test.go:1399\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"  error1\n" +
			"  stack trace \\(most recent call first\\):\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"    .+/format_test.go:1399\n" +
			"(.+\n    .+:\\d+\n)+" +
			"\n" +
			"  error2\n" +
			"  stack trace \\(most recent call first\\):\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"    .+/format_test.go:1399\n" +
			"(.+\n    .+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1419\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error joins multiple errors:\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1419\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1419\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #+.1v",
		"^error1\nerror2\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1438\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1438\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1438\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#+.1v",
		"^error1\nerror2\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1455\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1455\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1455\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #.1v",
		"^error1\nerror2\n" +
			"\n" +
			"\terror1\n" +
			"\n" +
			"\terror2\n$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#.1v",
		"^error1\nerror2\n" +
			"\terror1\n" +
			"\terror2\n$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-.1v",
		"^error1\nerror2\n" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\n" +
			"\terror2\n$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#-.1v",
		"^error1\nerror2\n" +
			"the above error joins multiple errors:\n" +
			"\terror1\n" +
			"\terror2\n$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1499\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1499\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1499\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+.3v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1519\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1519\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1519\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+2.3v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  .+/format_test.go:1537\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"  error1\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  \t.+/format_test.go:1537\n" +
			"(.+\n  \t.+:\\d+\n)+" +
			"\n" +
			"  error2\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  \t.+/format_test.go:1537\n" +
			"(.+\n  \t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(errors.Join(errors.New("error1"), errors.New("error2")), "message"),
		"% #-+.1v",
		"^message: error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1555\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1555\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1555\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(errors.Join(errors.New("error1"), errors.New("error2")), "foo", "bar"),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"foo=bar\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1575\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1575\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1575\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, errors.Formatter{tt.error}))
		})
	}
}
