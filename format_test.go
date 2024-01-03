package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

func TestFormatNew(t *testing.T) {
	t.Parallel()

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+15) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+22) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"  .+/format_test.go:" + strconv.Itoa(offset+30) + "\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.New("error"),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatNew\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+38) + "\n" +
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+19) + "\n",
	}, {
		errors.Errorf("%s", "error"),
		"%-+v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+25) + "\n",
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+2) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+56) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+4) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"% +-#v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+4) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+-2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"  .+/format_test.go:" + strconv.Itoa(offset+4) + "\n" +
			"(.+\n  .+:\\d+\n)+$",
	}, {
		errors.Errorf("%w", parentPkgErr),
		"%+v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatErrorf\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+4) + "\n" +
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+15) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+36) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+57) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.WithStack(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+68) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"^error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+78) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+100) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithStack(
			pkgerrors.New("error"),
		),
		"% +-#2v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithStack\n" +
			"  .+/format_test.go:" + strconv.Itoa(offset+110) + "\n" +
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+21) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+22) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%s",
		"^error$",
	}, {
		errors.Wrap(io.EOF, "error"),
		"%v",
		"^error$",
	}, {
		errors.Wrap(io.EOF, "error"),
		"% +-.1v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+46) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+56) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+57) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.Wrap(
			errors.New("error with space"),
			"context",
		),
		"%q",
		`^"context"$`,
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+96) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+97) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+113) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+114) + "\n" +
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
			"  .+/format_test.go:" + strconv.Itoa(offset+129) + "\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+130) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+145) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error was caused by the following error:\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+146) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"%+.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+162) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+163) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Wrap(
			errors.New("error"),
			"error2",
		),
		"% +.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+176) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrap\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+177) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+191) + "\n" +
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+15) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+39) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapf\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+40) + "\n" +
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+20) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+37) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+53) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+63) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+64) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+81) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+92) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+103) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+104) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+135) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+145) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+159) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"EOF\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithMessage\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+160) + "\n" +
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset-6) + "\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrappedNew\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+7) + "\n",
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+15) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+36) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+57) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.WithDetails(io.EOF),
		),
		"%+-v",
		"^EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+68) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.Errorf("error%d", 1),
		),
		"%+-v",
		"^error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+78) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+100) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(
			errors.New("error"),
		),
		"%+#v",
		"^error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWithDetails\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+110) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+119) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+158) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+174) + "\n" +
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

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+23) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+33) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+43) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormat{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+53) + "\n" +
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
			"\t.+/format_test.go:" + strconv.Itoa(offset+135) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+145) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.3v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+155) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"foobar\nmore data\n$",
	}, {
		errors.Wrap(&errorWithFormatAndStack{"foobar\nmore data\n"}, "read error"),
		"% #-+.1v",
		"^read error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatter\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+165) + "\n" +
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

			got := fmt.Sprintf(tt.format, errors.Formatter{Error: tt.error})
			assert.Regexp(t, tt.want, got)

			if !strings.Contains(got, "more data") {
				err2 := copyThroughJSON(t, errors.Formatter{Error: tt.error})
				got2 := fmt.Sprintf(tt.format, err2)
				assert.Equal(t, got, got2)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	offset := stackOffset(t)

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
			"\t.+/format_test.go:" + strconv.Itoa(offset+19) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error"), nil),
		"% #-+.1v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+27) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+35) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+43) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+43) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+43) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #-+2.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  .+/format_test.go:" + strconv.Itoa(offset+63) + "\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"  error1\n" +
			"  stack trace \\(most recent call first\\):\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"    .+/format_test.go:" + strconv.Itoa(offset+63) + "\n" +
			"(.+\n    .+:\\d+\n)+" +
			"\n" +
			"  error2\n" +
			"  stack trace \\(most recent call first\\):\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"    .+/format_test.go:" + strconv.Itoa(offset+63) + "\n" +
			"(.+\n    .+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+83) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error joins errors:\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+83) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+83) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"% #+.1v",
		"^error1\nerror2\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+102) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+102) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+102) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#+.1v",
		"^error1\nerror2\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+119) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+119) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+119) + "\n" +
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
			"\nthe above error joins errors:\n\n" +
			"\terror1\n" +
			"\n" +
			"\terror2\n$",
	}, {
		errors.Join(errors.New("error1"), errors.New("error2")),
		"%#-.1v",
		"^error1\nerror2\n" +
			"the above error joins errors:\n" +
			"\terror1\n" +
			"\terror2\n$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+163) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+163) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+163) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+.3v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+183) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror1\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+183) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+183) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Join(pkgerrors.New("error1"), pkgerrors.New("error2")),
		"% #-+2.3v",
		"^error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  .+/format_test.go:" + strconv.Itoa(offset+201) + "\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"  error1\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  \t.+/format_test.go:" + strconv.Itoa(offset+201) + "\n" +
			"(.+\n  \t.+:\\d+\n)+" +
			"\n" +
			"  error2\n" +
			"  gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"  \t.+/format_test.go:" + strconv.Itoa(offset+201) + "\n" +
			"(.+\n  \t.+:\\d+\n)+$",
	}, {
		errors.WithMessage(errors.Join(errors.New("error1"), errors.New("error2")), "message"),
		"% #-+.1v",
		"^message: error1\nerror2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+219) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+219) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+219) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.WithDetails(errors.Join(errors.New("error1"), errors.New("error2")), "foo", "bar"),
		"% #-+.1v",
		"^error1\nerror2\n" +
			"foo=bar\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+239) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror1\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+239) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+" +
			"\n" +
			"\terror2\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestJoin\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+239) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, errors.Formatter{Error: tt.error})
			assert.Regexp(t, tt.want, got)

			if !strings.Contains(tt.format, ".3v") {
				err2 := copyThroughJSON(t, errors.Formatter{Error: tt.error})
				got2 := fmt.Sprintf(tt.format, err2)
				assert.Equal(t, got, got2)
			}
		})
	}
}

func TestFormatWrapWith(t *testing.T) {
	t.Parallel()

	offset := stackOffset(t)

	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"%s",
		"^error2$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"%v",
		"^error2$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+21) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+22) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(io.EOF, errors.Base("error")),
		"%s",
		"^error$",
	}, {
		errors.WrapWith(io.EOF, errors.Base("error")),
		"%v",
		"^error$",
	}, {
		errors.WrapWith(io.EOF, errors.Base("error")),
		"% +-.1v",
		"^error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+46) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.WrapWith(
			errors.WrapWith(io.EOF, errors.Base("error1")),
			errors.Base("error2"),
		),
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+56) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error1\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+57) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"EOF\n$",
	}, {
		errors.WrapWith(
			errors.New("error with space"),
			errors.Base("context"),
		),
		"%q",
		`^"context"$`,
	}, {
		errors.WrapWith(
			pkgerrors.New("error"),
			errors.Base("error2"),
		),
		"%s",
		"^error2$",
	}, {
		errors.WrapWith(
			pkgerrors.New("error"),
			errors.Base("error2"),
		),
		"%v",
		"^error2$",
	}, {
		errors.WrapWith(
			pkgerrors.New("error"),
			errors.Base("error2"),
		),
		"% +-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+96) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+97) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(
			pkgerrors.New("error"),
			errors.Base("error2"),
		),
		"% +-.3v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+113) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+114) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(
			pkgerrors.New("error"),
			errors.Base("error2"),
		),
		"% +-2.3v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"  .+/format_test.go:" + strconv.Itoa(offset+129) + "\n" +
			"(.+\n  .+:\\d+\n)+" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+130) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"%+-.1v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+145) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"the above error was caused by the following error:\n" +
			"error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+146) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"%+.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+162) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+163) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"% +.1v",
		"^error2\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+176) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\n" +
			"error\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+177) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"% +-#v",
		"^error2\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatWrapWith\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+191) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"% -.1v",
		"^error2\n" +
			"\nthe above error was caused by the following error:\n\n" +
			"error\n$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
		),
		"%-.1v",
		"^error2\n" +
			"the above error was caused by the following error:\n" +
			"error\n$",
	}, {
		errors.WrapWith(
			errors.New("error"),
			errors.Base("error2"),
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

type testStructJoined struct {
	msg     string  `exhaustruct:"optional"`
	cause   error   `exhaustruct:"optional"`
	parents []error `exhaustruct:"optional"`
}

func (e *testStructJoined) Error() string {
	return e.msg
}

func (e *testStructJoined) Cause() error {
	return e.cause
}

func (e *testStructJoined) Unwrap() []error {
	return e.parents
}

func TestFormatCustomError(t *testing.T) {
	t.Parallel()

	testErr := &testStructJoined{msg: "test2"}

	tests := []struct {
		error
		format string
		want   string
	}{{
		&testStructJoined{},
		"%s",
		"^$",
	}, {
		&testStructJoined{msg: "test"},
		"%s",
		"^test$",
	}, {
		&testStructJoined{msg: "test"},
		"%-.1v",
		"^test\n$",
	}, {
		&testStructJoined{msg: "test1", cause: testErr},
		"%-.1v",
		"^test1\nthe above error was caused by the following error:\ntest2\n$",
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{testErr}},
		"%-.1v",
		"^test1\nthe above error was caused by the following error:\ntest2\n$",
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{testErr, &testStructJoined{msg: "test3"}}},
		"%-.1v",
		"^test1\nthe above error joins errors:\n\ttest3\nthe above error was caused by the following error:\ntest2\n$",
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{testErr, &testStructJoined{msg: "test3"}, &testStructJoined{msg: "test4"}}},
		"%-.1v",
		"^test1\nthe above error joins errors:\n\ttest3\n\ttest4\nthe above error was caused by the following error:\ntest2\n$",
	}, {
		&testStructJoined{msg: "test1", cause: testErr, parents: []error{&testStructJoined{msg: "test3"}, &testStructJoined{msg: "test4"}}},
		"%-.1v",
		"^test1\nthe above error joins errors:\n\ttest3\n\ttest4\nthe above error was caused by the following error:\ntest2\n$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			got := fmt.Sprintf(tt.format, errors.Formatter{Error: tt.error})
			assert.Regexp(t, tt.want, got)

			err2 := copyThroughJSON(t, errors.Formatter{Error: tt.error})
			got2 := fmt.Sprintf(tt.format, err2)
			assert.Equal(t, got, got2)
		})
	}
}

func TestFormatPrefix(t *testing.T) {
	t.Parallel()

	offset := stackOffset(t)

	tests := []struct {
		error
		format string
		want   string
	}{{
		errors.Prefix(
			errors.New("error"),
			errors.Base("error2"),
		),
		"%s",
		"^error2: error$",
	}, {
		errors.Prefix(
			errors.New("error"),
			errors.Base("error2"),
		),
		"%v",
		"^error2: error$",
	}, {
		errors.Prefix(
			errors.New("error"),
			errors.Base("error2"),
		),
		"% +-.1v",
		"^error2: error\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatPrefix\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+22) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror2\n\n" +
			"\terror\n" +
			"\tstack trace \\(most recent call first\\):\n" +
			"\tgitlab.com/tozd/go/errors_test.TestFormatPrefix\n" +
			"\t\t.+/format_test.go:" + strconv.Itoa(offset+22) + "\n" +
			"(.+\n\t\t.+:\\d+\n)+$",
	}, {
		errors.Prefix(
			errors.Base("parent"),
			errors.Base(""),
		),
		"% +-.1v",
		"^parent\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatPrefix\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+39) + "\n" +
			"(.+\n\t.+:\\d+\n)+$",
	}, {
		errors.Prefix(io.EOF, errors.Base("error")),
		"%s",
		"^error: EOF$",
	}, {
		errors.Prefix(io.EOF, errors.Base("error")),
		"%v",
		"^error: EOF$",
	}, {
		errors.Prefix(io.EOF, errors.Base("error")),
		"% +-.1v",
		"^error: EOF\n" +
			"stack trace \\(most recent call first\\):\n" +
			"gitlab.com/tozd/go/errors_test.TestFormatPrefix\n" +
			"\t.+/format_test.go:" + strconv.Itoa(offset+58) + "\n" +
			"(.+\n\t.+:\\d+\n)+" +
			"\nthe above error joins errors:\n\n" +
			"\terror\n\n" +
			"\tEOF\n$",
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

func TestGetMessage(t *testing.T) {
	t.Parallel()

	got := fmt.Sprintf("%s", errors.Formatter{Error: errors.New("test"), GetMessage: func(err error) string {
		return fmt.Sprintf("X%sX", err.Error())
	}})
	assert.Equal(t, "XtestX", got)
}
