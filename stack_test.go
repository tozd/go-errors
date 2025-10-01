package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	initpc = caller()
	zeropc = func() frame {
		frames := runtime.CallersFrames([]uintptr{0})
		f, _ := frames.Next()
		return frame(f)
	}()
)

type X struct{} //nolint:recvcheck

// val returns a frame pointing to itself.
//
//go:noinline
func (x X) val() frame {
	return caller()
}

// ptr returns a frame pointing to itself.
//
//go:noinline
func (x *X) ptr() frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		frame

		format string
		want   string
	}{{
		initpc,
		"%s",
		"^stack_test.go$",
	}, {
		initpc,
		"%+s",
		"^gitlab.com/tozd/go/errors.init\n" +
			"\t.+/stack_test.go$",
	}, {
		zeropc,
		"%s",
		"^unknown$",
	}, {
		zeropc,
		"%+s",
		"^unknown\n\tunknown$",
	}, {
		initpc,
		"%d",
		"^15$",
	}, {
		zeropc,
		"%d",
		"^0$",
	}, {
		initpc,
		"%n",
		"^init$",
	}, {
		func() frame {
			var x X
			return x.ptr()
		}(),
		"%n",
		`^\(\*X\).ptr$`,
	}, {
		func() frame {
			var x X
			return x.val()
		}(),
		"%n",
		"^X.val$",
	}, {
		zeropc,
		"%n",
		"^unknown$",
	}, {
		initpc,
		"%v",
		"^stack_test.go:15$",
	}, {
		initpc,
		"%+v",
		"^gitlab.com/tozd/go/errors.init\n" +
			"\t.+/stack_test.go:15$",
	}, {
		zeropc,
		"%v",
		"^unknown:0$",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.frame))
		})
	}
}

func TestFuncname(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, want string
	}{
		{"", ""},
		{"runtime.main", "main"},
		{"gitlab.com/tozd/go/errors.funcname", "funcname"},
		{"funcname", "funcname"},
		{"io.copyBuffer", "copyBuffer"},
		{"main.(*R).Write", "(*R).Write"},
	}
	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, funcname(tt.name))
		})
	}
}

func TestStackFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err    error
		format string
		want   string
	}{{
		New("ooh"),
		"%+v",
		"^gitlab.com/tozd/go/errors.TestStackFormatter\n" +
			"\t.+/stack_test.go:151\n",
	}, {
		Wrap(
			New("ooh"),
			"ahh",
		),
		"%+v",
		"^gitlab.com/tozd/go/errors.TestStackFormatter\n" +
			"\t.+/stack_test.go:156\n",
	}, {
		func() error {
			noinline()
			return New("ooh")
		}(),
		"%+v",
		"^gitlab.com/tozd/go/errors.TestStackFormatter.func1\n" +
			"\t.+/stack_test.go:166\n" +
			"gitlab.com/tozd/go/errors.TestStackFormatter\n" +
			"\t.+/stack_test.go:167\n",
	}, {
		func() error {
			return func() error {
				noinline()
				return Errorf("hello %s", "world: ooh")
			}()
		}(),
		"%+v",
		// Nested function names changed in Go 1.21: https://github.com/golang/go/issues/62132
		"^gitlab.com/tozd/go/errors.(TestStackFormatter.){1,2}func2.(1|func5)\n" +
			"\t.+/stack_test.go:177\n" +
			"gitlab.com/tozd/go/errors.TestStackFormatter.func2\n" +
			"\t.+/stack_test.go:178\n" +
			"gitlab.com/tozd/go/errors.TestStackFormatter\n" +
			"\t.+/stack_test.go:179\n",
	}}

	for k, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()

			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, StackFormatter{tt.err.(stackTracer).StackTrace()})) //nolint:forcetypeassert,errcheck
		})
	}

	stack := func() []uintptr {
		return func() []uintptr {
			noinline()
			return callers(0)
		}()
	}()

	assert.Regexp(t, "^gitlab.com/tozd/go/errors.TestStackFormatter.func4\n"+
		"\t.+/stack_test.go:204\n"+
		"gitlab.com/tozd/go/errors.TestStackFormatter\n"+
		"\t.+/stack_test.go:205\n", fmt.Sprintf("%+v", StackFormatter{stack}))

	assert.Regexp(t, "^gitlab.com/tozd/go/errors.TestStackFormatter.func4\n"+
		"\t.+/stack_test.go\n"+
		"gitlab.com/tozd/go/errors.TestStackFormatter\n"+
		"\t.+/stack_test.go\n", fmt.Sprintf("%+s", StackFormatter{stack}))

	assert.Regexp(t, "^gitlab.com/tozd/go/errors.TestStackFormatter.func4\n"+
		"  .+/stack_test.go:204\n"+
		"gitlab.com/tozd/go/errors.TestStackFormatter\n"+
		"  .+/stack_test.go:205\n", fmt.Sprintf("%+2v", StackFormatter{stack}))

	assert.Regexp(t, "^gitlab.com/tozd/go/errors.TestStackFormatter.func4\n"+
		"  .+/stack_test.go\n"+
		"gitlab.com/tozd/go/errors.TestStackFormatter\n"+
		"  .+/stack_test.go\n", fmt.Sprintf("%+2s", StackFormatter{stack}))

	assert.Empty(t, fmt.Sprintf("%+v", StackFormatter{nil}))

	assert.Regexp(t, "^%!f\\(errors.frame=stack_test.go:204\\)\n"+
		"%!f\\(errors.frame=stack_test.go:205\\)\n", fmt.Sprintf("%f", StackFormatter{stack}))

	assert.Regexp(t, "^stack_test.go\n"+
		"stack_test.go\n", fmt.Sprintf("%s", StackFormatter{stack}))

	assert.Regexp(t, "^204\n"+
		"205\n", fmt.Sprintf("%d", StackFormatter{stack}))

	assert.Regexp(t, "^TestStackFormatter.func4\n"+
		"TestStackFormatter\n", fmt.Sprintf("%n", StackFormatter{stack}))

	assert.Regexp(t, "^stack_test.go:204\n"+
		"stack_test.go:205\n", fmt.Sprintf("%v", StackFormatter{stack}))
}

func TestStackMarshalJSON(t *testing.T) {
	t.Parallel()

	stack := func() []uintptr {
		return func() []uintptr {
			noinline()
			return callers(0)
		}()
	}()
	j, err := json.Marshal(StackFormatter{stack})
	require.NoError(t, err)
	var d []struct {
		Name string `json:"name"`
		File string `json:"file"`
		Line int    `json:"line"`
	}
	decoder := json.NewDecoder(bytes.NewReader(j))
	decoder.DisallowUnknownFields()
	e := decoder.Decode(&d)
	require.NoError(t, e)
	assert.Equal(t, 252, d[0].Line)
	assert.Equal(t, 253, d[1].Line)

	j, err = json.Marshal(StackFormatter{nil})
	require.NoError(t, err)
	assert.Equal(t, "[]", string(j))
}

// A version of runtime.Caller that returns a frame, not a uintptr.
func caller() frame {
	var pcs [1]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	f, _ := frames.Next()
	return frame(f)
}

// noinline prevents the caller being inlined.
//
//go:noinline
func noinline() {}
