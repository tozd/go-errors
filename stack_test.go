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

type X struct{}

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
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, tt.frame))
		})
	}
}

func TestFuncname(t *testing.T) {
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
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Equal(t, tt.want, funcname(tt.name))
		})
	}
}

func TestStackFormat(t *testing.T) {
	tests := []struct {
		err    error
		format string
		want   string
	}{{
		New("ooh"),
		"%+v",
		"^gitlab.com/tozd/go/errors.TestStackFormat\n" +
			"\t.+/stack_test.go:137\n",
	}, {
		Wrap(
			New("ooh"),
			"ahh",
		),
		"%+v",
		"^gitlab.com/tozd/go/errors.TestStackFormat\n" +
			"\t.+/stack_test.go:142\n",
	}, {
		func() error {
			noinline()
			return New("ooh")
		}(),
		"%+v",
		"^gitlab.com/tozd/go/errors.TestStackFormat.func1\n" +
			"\t.+/stack_test.go:152\n" +
			"gitlab.com/tozd/go/errors.TestStackFormat\n" +
			"\t.+/stack_test.go:153\n",
	}, {
		func() error {
			return func() error {
				noinline()
				return Errorf("hello %s", fmt.Sprintf("world: %s", "ooh"))
			}()
		}(),
		"%+v",
		"^gitlab.com/tozd/go/errors.TestStackFormat.func2.1\n" +
			"\t.+/stack_test.go:163\n" +
			"gitlab.com/tozd/go/errors.TestStackFormat.func2\n" +
			"\t.+/stack_test.go:164\n" +
			"gitlab.com/tozd/go/errors.TestStackFormat\n" +
			"\t.+/stack_test.go:165\n",
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.Regexp(t, tt.want, fmt.Sprintf(tt.format, stack(tt.err.(stackTracer).StackTrace())))
		})
	}
}

func TestStackFormatFunc(t *testing.T) {
	stack := func() []uintptr {
		return func() []uintptr {
			noinline()
			return callers().StackTrace()
		}()
	}()
	output := StackFormat("%+v", stack)
	assert.Regexp(t, "^gitlab.com/tozd/go/errors.TestStackFormatFunc.func1\n"+
		"\t.+/stack_test.go:187\n"+
		"gitlab.com/tozd/go/errors.TestStackFormatFunc\n"+
		"\t.+/stack_test.go:188\n", output)
}

func TestStackMarshalJSON(t *testing.T) {
	stack := func() []uintptr {
		return func() []uintptr {
			noinline()
			return callers().StackTrace()
		}()
	}()
	j, err := StackMarshalJSON(stack)
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
	assert.Equal(t, 201, d[0].Line)
	assert.Equal(t, 202, d[1].Line)
}

// A version of runtime.Caller that returns a frame, not a uintptr.
func caller() frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	f, _ := frames.Next()
	return frame(f)
}

// noinline prevents the caller being inlined.
//
//go:noinline
func noinline() {}
