//go:build go1.18

package errors_test

// Tests in this file copied from errors/wrap_test.go available from Go 1.26 on.

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"gitlab.com/tozd/go/errors"
)

type errorT struct{ s string } //nolint:errname

func (e errorT) Error() string { return fmt.Sprintf("errorT(%s)", e.s) }

type wrapped struct {
	msg string
	err error
}

func (e wrapped) Error() string { return e.msg }
func (e wrapped) Unwrap() error { return e.err }

type multiErr []error //nolint:errname

func (m multiErr) Error() string   { return "multiError" }
func (m multiErr) Unwrap() []error { return []error(m) }

type poser struct {
	msg string
	f   func(error) bool
}

var poserPathErr = &fs.PathError{Op: "poser"} //nolint:errname,exhaustruct,gochecknoglobals

func (p *poser) Error() string     { return p.msg }
func (p *poser) Is(err error) bool { return p.f(err) }
func (p *poser) As(err any) bool {
	switch x := err.(type) {
	case **poser:
		*x = p
	case *errorT:
		*x = errorT{"poser"}
	case **fs.PathError:
		*x = poserPathErr
	default:
		return false
	}
	return true
}

func TestAsType(t *testing.T) {
	t.Parallel()

	var errT errorT
	var errP *fs.PathError
	type timeout interface {
		Timeout() bool
		error
	}
	_, errF := os.Open("non-existing")
	poserErr := &poser{"oh no", nil}

	testAsType(t,
		nil,
		errP,
		false,
	)
	testAsType(t,
		wrapped{"pitied the fool", errorT{"T"}},
		errorT{"T"},
		true,
	)
	testAsType(t,
		errF,
		errF,
		true,
	)
	testAsType(t,
		errT,
		errP,
		false,
	)
	testAsType(t,
		wrapped{"wrapped", nil},
		errT,
		false,
	)
	testAsType(t,
		&poser{"error", nil},
		errorT{"poser"},
		true,
	)
	testAsType(t,
		&poser{"path", nil},
		poserPathErr,
		true,
	)
	testAsType(t,
		poserErr,
		poserErr,
		true,
	)
	testAsType(t,
		errors.New("err"),
		timeout(nil),
		false,
	)
	testAsType(t,
		errF,
		func() timeout {
			var target timeout
			_ = errors.As(errF, &target)
			return target
		}(),
		true)
	testAsType(t,
		wrapped{"path error", errF},
		func() timeout {
			var target timeout
			_ = errors.As(errF, &target)
			return target
		}(),
		true,
	)
	testAsType(t,
		multiErr{},
		errT,
		false,
	)
	testAsType(t,
		multiErr{errors.New("a"), errorT{"T"}},
		errorT{"T"},
		true,
	)
	testAsType(t,
		multiErr{errorT{"T"}, errors.New("a")},
		errorT{"T"},
		true,
	)
	testAsType(t,
		multiErr{errorT{"a"}, errorT{"b"}},
		errorT{"a"},
		true,
	)
	testAsType(t,
		multiErr{multiErr{errors.New("a"), errorT{"a"}}, errorT{"b"}},
		errorT{"a"},
		true,
	)
	testAsType(t,
		multiErr{wrapped{"path error", errF}},
		func() timeout {
			var target timeout
			_ = errors.As(errF, &target)
			return target
		}(),
		true,
	)
	testAsType(t,
		multiErr{nil},
		errT,
		false,
	)
}

type compError interface {
	comparable
	error
}

func testAsType[E compError](t *testing.T, err error, want E, wantOK bool) {
	t.Helper()
	name := fmt.Sprintf("AsType[%T](Errorf(..., %v))", want, err)
	t.Run(name, func(t *testing.T) {
		got, gotOK := errors.AsType[E](err)
		if gotOK != wantOK || !errors.Is(got, want) {
			t.Fatalf("got %v, %t; want %v, %t", got, gotOK, want, wantOK)
		}
	})
}
