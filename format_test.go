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
			"\t.+/format_test.go:29\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:36\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"  .+/format_test.go:44\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:52\n" +
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
		errors.New("error\n"), //nolint:revive
		"% v",
		"^error\n$",
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
			"\t.+/format_test.go:98\n",
	}, {
		errors.Errorf("%s", "error"),
		"%-+v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:104\n",
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
			"\t.+/format_test.go:81\n" +
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
			"\t.+/format_test.go:135\n" +
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
			"\t.+/format_test.go:83\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:83\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"  .+/format_test.go:83\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:83\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWithStack(t *testing.T) { //nolint: dupl
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
			"\t.+/format_test.go:204\n" +
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
			"\t.+/format_test.go:225\n" +
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
			"\t.+/format_test.go:246\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.WithStack(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:257\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:267\n" +
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
			"\t.+/format_test.go:289\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			pkgerrors.New("error"),
		),
		"% +-#2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"  .+/format_test.go:299\n" +
			"(.+\n  .+:\\d+\n)+$",
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
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:336\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:337\n" +
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
			"\t.+/format_test.go:361\n" +
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
			"\t.+/format_test.go:371\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:372\n" +
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
			"\t.+/format_test.go:411\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:412\n" +
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
			"\t.+/format_test.go:428\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:429\n" +
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
			"  .+/format_test.go:444\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:445\n" +
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
			"\t.+/format_test.go:460\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error was caused by the following error:\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:461\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%+.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:477\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:478\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"% +.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:491\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:492\n" +
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
			"\t.+/format_test.go:506\n" +
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
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:565\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
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
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:589\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:590\n" +
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
		"% +-.1v",
		"^error2: error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:633\n" +
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
			"\t.+/format_test.go:650\n" +
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
			"\t.+/format_test.go:666\n" +
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
			"\t.+/format_test.go:676\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:677\n" +
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
			"\t.+/format_test.go:694\n" +
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
			"\t.+/format_test.go:705\n" +
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
			"\t.+/format_test.go:716\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:717\n" +
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
			"\t.+/format_test.go:748\n" +
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
			"\t.+/format_test.go:758\n" +
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
			"\t.+/format_test.go:772\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"EOF\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:773\n" +
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
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.wrappedNew\n" +
			"\t.+/format_test.go:821\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrappedNew\n" +
			"\t.+/format_test.go:830\n",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatWithDetails(t *testing.T) { //nolint: dupl
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
			"\t.+/format_test.go:861\n" +
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
			"\t.+/format_test.go:882\n" +
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
			"\t.+/format_test.go:903\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.WithDetails(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:914\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:924\n" +
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
			"\t.+/format_test.go:946\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
		),
		"%+#v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:956\n" +
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
			"\t.+/format_test.go:965\n" +
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
			"\t.+/format_test.go:1004\n" +
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
			"\t.+/format_test.go:1020\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.error))
		})
	}
}

func TestFormatter(t *testing.T) {
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
			"\t.+/format_test.go:1065\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1075\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1085\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1095\n" +
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
			"\t.+/format_test.go:1177\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1187\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1197\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:1207\n" +
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
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, errors.Formatter{tt.error}))
		})
	}
}
