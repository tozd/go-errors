package errors_test

import (
	"fmt"

	"gitlab.com/tozd/go/errors"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)
	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// Stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleNew_printf
	// 	/home/user/errors/example_test.go:16
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:87
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1581
}

func ExampleWithMessage() {
	cause := errors.New("whoops")
	err := errors.WithMessage(cause, "oh noes")
	fmt.Println(err)
	// Output: oh noes: whoops
}

func ExampleWithMessage_printf() {
	cause := errors.New("whoops")
	err := errors.WithMessage(cause, "oh noes")
	fmt.Printf("%+v", err)

	// Example output:
	// oh noes: whoops
	// Stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWithMessage_printf
	// 	/home/user/errors/example_test.go:46
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:97
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1581
}

func ExampleWithMessagef() {
	cause := errors.New("whoops")
	err := errors.WithMessagef(cause, "oh noes #%d", 2)
	fmt.Println(err)
	// Output: oh noes #2: whoops
}

func ExampleWithStack() {
	base := errors.Base("whoops")
	err := errors.WithStack(base)
	fmt.Println(err)
	// Output: whoops
}

func ExampleWithStack_printf() {
	base := errors.Base("whoops")
	err := errors.WithStack(base)
	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// Stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWithStack_printf
	// 	/home/user/errors/example_test.go:54
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:91
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1581
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")
	fmt.Println(err)
	// Output: oh noes
}

func ExampleWrap_printf() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")
	fmt.Printf("%+v", err)

	// Example output:
	// oh noes
	// Stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWrap_printf
	// 	/home/user/errors/example_test.go:86
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:93
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1581
	//
	// The above error was caused by the following error:
	//
	// whoops
	// Stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWrap_printf
	// 	/home/user/errors/example_test.go:85
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:93
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1581
}

func ExampleWrapf() {
	base := errors.Base("whoops")
	err := errors.Wrapf(base, "oh noes #%d", 2)
	fmt.Println(err)
	// Output: oh noes #2
}

func ExampleErrorf() {
	err := errors.Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops: foo
	// Stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleErrorf
	// 	/home/user/errors/example_test.go:134
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:95
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1581
}

func ExampleErrorf_wrap() {
	base := errors.Base("whoops")
	err := errors.Errorf("oh noes (%w)", base)
	fmt.Printf("%+v", err)

	// Example output:
	// oh noes (whoops)
	// Stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleErrorf_wrap
	// 	/home/user/errors/example_test.go:189
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:64
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1505
	// main.main
	// 	_testmain.go:99
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:255
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1581
}

func ExampleBase() {
	err := errors.Base("whoops")
	fmt.Println(err)
	// Output: whoops
}

func ExampleBasef() {
	err := errors.Basef("whoops #%d", 2)
	fmt.Println(err)
	// Output: whoops #2
}

func ExampleBaseWrap() {
	base := errors.Base("error")
	valueError := errors.BaseWrap(base, "value")
	fmt.Println(valueError)
	// Output: value
}

func ExampleBaseWrapf() {
	base := errors.Base("error")
	valueError := errors.BaseWrapf(base, "value %d", 2)
	fmt.Println(valueError)
	// Output: value 2
}

func ExampleCause() {
	base := errors.Base("error")
	wrapped := errors.Wrap(base, "wrapped")
	fmt.Println(errors.Cause(wrapped))
	// Output: error
}

func ExampleUnwrap() {
	base := errors.Base("error")
	withPrefix := errors.WithMessage(base, "prefix")
	fmt.Println(withPrefix)
	fmt.Println(errors.Unwrap(withPrefix))
	// Output:
	// prefix: error
	// error
}

func ExampleIs() {
	base := errors.Base("error")
	valueError := errors.BaseWrap(base, "value")
	fmt.Println(errors.Is(valueError, base))
	// Output: true
}
