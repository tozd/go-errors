package errors

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var initpc = caller()
var zeropc = func() frame {
	frames := runtime.CallersFrames([]uintptr{0})
	f, _ := frames.Next()
	return frame(f)
}()

type X struct{}

// val returns a frame pointing to itself.
//go:noinline
func (x X) val() frame {
	return caller()
}

// ptr returns a frame pointing to itself.
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
		"^13$",
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
		"^stack_test.go:13$",
	}, {
		initpc,
		"%+v",
		"^gitlab.com/tozd/go/errors.init\n" +
			"\t.+/stack_test.go:13$",
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

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want []string
	}{{
		New("ooh"), []string{
			"^gitlab.com/tozd/go/errors.TestStackTrace\n" +
				"\t.+/stack_test.go:131$",
		},
	}, {
		Wrap(New("ooh"), "ahh"), []string{
			"^gitlab.com/tozd/go/errors.TestStackTrace\n" +
				"\t.+/stack_test.go:136$", // this is the stack of Wrap, not New
		},
	}, {
		func() error { noinline(); return New("ooh") }(), []string{
			`^gitlab.com/tozd/go/errors.TestStackTrace.func1` +
				"\n\t.+/stack_test.go:141$", // this is the stack of New
			"^gitlab.com/tozd/go/errors.TestStackTrace\n" +
				"\t.+/stack_test.go:141$", // this is the stack of New's caller
		},
	}, {
		func() error {
			return func() error {
				noinline()
				return Errorf("hello %s", fmt.Sprintf("world: %s", "ooh"))
			}()
		}(), []string{
			`^gitlab.com/tozd/go/errors.TestStackTrace.func2.1` +
				"\n\t.+/stack_test.go:151$", // this is the stack of Errorf
			`^gitlab.com/tozd/go/errors.TestStackTrace.func2` +
				"\n\t.+/stack_test.go:152$", // this is the stack of Errorf's caller
			"^gitlab.com/tozd/go/errors.TestStackTrace\n" +
				"\t.+/stack_test.go:153$", // this is the stack of Errorf's caller's caller
		},
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			require.Implements(t, (*stackTracer)(nil), tt.err)
			st := fmt.Sprintf("%+v", stack(tt.err.(stackTracer).StackTrace()))
			stackLines := strings.Split(st, "\n")[1:]
			for i := 0; i < len(tt.want); i++ {
				assert.Regexp(t, tt.want[i], stackLines[2*i]+"\n"+stackLines[2*i+1])
			}
		})
	}
}

// A version of runtime.Caller that returns a frame, not a uintptr.
func caller() frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	f, _ := frames.Next()
	return frame(f)
}

//go:noinline
// noinline prevents the caller being inlined
func noinline() {}
