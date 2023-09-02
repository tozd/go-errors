package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"strings"
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
			"\t.+/format_test.go:32\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:39\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"  .+/format_test.go:47\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:55\n" +
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

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, tt.error)
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
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
			"\t.+/format_test.go:112\n",
	}, {
		errors.Errorf("%s", "error"),
		"%-+v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:118\n",
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
			"\t.+/format_test.go:95\n" +
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
			"\t.+/format_test.go:149\n" +
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
			"\t.+/format_test.go:97\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:97\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"  .+/format_test.go:97\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:97\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, tt.error)
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
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
			"\t.+/format_test.go:229\n" +
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
			"\t.+/format_test.go:250\n" +
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
			"\t.+/format_test.go:271\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.WithStack(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:282\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:292\n" +
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
			"\t.+/format_test.go:314\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			pkgerrors.New("error"),
		),
		"% +-#2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"  .+/format_test.go:324\n" +
			"(.+\n  .+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, tt.error)
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
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
			"\t.+/format_test.go:372\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:373\n" +
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
			"\t.+/format_test.go:397\n" +
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
			"\t.+/format_test.go:407\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:408\n" +
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
			"\t.+/format_test.go:447\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:448\n" +
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
			"\t.+/format_test.go:464\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:465\n" +
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
			"  .+/format_test.go:480\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:481\n" +
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
			"\t.+/format_test.go:496\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error was caused by the following error:\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:497\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%+.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:513\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:514\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"% +.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:527\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:528\n" +
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
			"\t.+/format_test.go:542\n" +
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

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			if !strings.Contains(tt.format, ".3v") {
				err2 := copyThroughJSON(t, tt.error)
				got2 := fmt.Sprintf(tt.format, err2)
				assert.Equal(t, got, got2)
			}
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
			"\t.+/format_test.go:614\n" +
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
			"\t.+/format_test.go:638\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:639\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, tt.error)
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
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
			"\t.+/format_test.go:693\n" +
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
			"\t.+/format_test.go:710\n" +
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
			"\t.+/format_test.go:726\n" +
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
			"\t.+/format_test.go:736\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:737\n" +
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
			"\t.+/format_test.go:754\n" +
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
			"\t.+/format_test.go:765\n" +
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
			"\t.+/format_test.go:776\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:777\n" +
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
			"\t.+/format_test.go:808\n" +
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
			"\t.+/format_test.go:818\n" +
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
			"\t.+/format_test.go:832\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"EOF\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:833\n" +
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

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, tt.error)
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
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
			"\t.+/format_test.go:890\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrappedNew\n" +
			"\t.+/format_test.go:901\n",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, tt.error)
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
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
			"\t.+/format_test.go:943\n" +
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
			"\t.+/format_test.go:964\n" +
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
			"\t.+/format_test.go:985\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.WithDetails(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:996\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:1006\n" +
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
			"\t.+/format_test.go:1028\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
		),
		"%+#v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:1038\n" +
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
			"\t.+/format_test.go:1047\n" +
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
			"\t.+/format_test.go:1086\n" +
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
			"\t.+/format_test.go:1102\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, tt.error)
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, tt.error)
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
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
			"\t.+/format_test.go:1158\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1168\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1178\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1188\n" +
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
			"\t.+/format_test.go:1270\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1280\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1290\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1300\n" +
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

			got := fmt.Sprintf(tt.format, errors.Formatter{tt.error})
			assert.Regexp(t, tt.want, got)

			if !strings.Contains(got, "more data") {
				err2 := copyThroughJSON(t, errors.Formatter{tt.error})
				got2 := fmt.Sprintf(tt.format, errors.Formatter{err2})
				assert.Equal(t, got, got2)
			}
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
			"\t.+/format_test.go:1405\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error"), nil),
		"% #-+.1v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1413\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1421\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1429\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1429\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1429\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+2.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  .+/format_test.go:1449\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"  error1\n" +
			"  stack trace \\(most recent call first\\):\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"    .+/format_test.go:1449\n" +
			"(.+\n    .+:\\d+\n)+" +
			"\n" +
			"  error2\n" +
			"  stack trace \\(most recent call first\\):\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"    .+/format_test.go:1449\n" +
			"(.+\n    .+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1469\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error joins multiple errors:\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1469\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1469\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #+.1v",
		"^error1\nerror2\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1488\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1488\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1488\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#+.1v",
		"^error1\nerror2\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1505\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1505\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1505\n" +
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
			"\t.+/format_test.go:1549\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1549\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1549\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+.3v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1569\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1569\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1569\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+2.3v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  .+/format_test.go:1587\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"  error1\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  \t.+/format_test.go:1587\n" +
			"(.+\n  \t.+:\\d+\n)+" +
			"\n" +
			"  error2\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  \t.+/format_test.go:1587\n" +
			"(.+\n  \t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(errors.Join(errors.New("error1"), errors.New("error2")), "message"),
		"% #-+.1v",
		"^message: error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1605\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1605\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1605\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(errors.Join(errors.New("error1"), errors.New("error2")), "foo", "bar"),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"foo=bar\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:1625\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins multiple errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1625\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:1625\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, errors.Formatter{tt.error})
			assert.Regexp(t, tt.want, got)

			if !strings.Contains(tt.format, ".3v") {
				err2 := copyThroughJSON(t, errors.Formatter{tt.error})
				got2 := fmt.Sprintf(tt.format, errors.Formatter{err2})
				assert.Equal(t, got, got2)
			}
		})
	}
}
