package errors_test

import (
	"fmt"
	"runtime"

	"gitlab.com/tozd/go/errors"
)

func getErr() error {
	return errors.New("foobar")
}

func Example_stackTrace() {
	type stackTracer interface {
		StackTrace() []uintptr
	}

	var err stackTracer
	if !errors.As(getErr(), &err) {
		panic(errors.New("oops, err does not implement stackTracer"))
	}

	frames := runtime.CallersFrames(err.StackTrace())
	frame, _ := frames.Next()
	fmt.Printf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line)

	// Example output:
	// gitlab.com/tozd/go/errors_test.getErr
	//	/home/user/errors/example_test.go:11
}
