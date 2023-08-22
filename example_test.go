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
	fmt.Printf("%+-v", err)

	// Example output:
	// whoops
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleNew_printf
	// 	/home/user/errors/example_test.go:16
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
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
	fmt.Printf("%+-v", err)

	// Example Output:
	// oh noes: whoops
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWithMessage_printf
	// 	/home/user/errors/example_test.go:46
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
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
	fmt.Printf("%+-v", err)

	// Example output:
	// whoops
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWithStack_printf
	// 	/home/user/errors/example_test.go:85
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
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
	fmt.Printf("% +-.1v", err)

	// Example output:
	// oh noes
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWrap_printf
	// 	/home/user/errors/example_test.go:116
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
	//
	// the above error was caused by the following error:
	//
	// whoops
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWrap_printf
	// 	/home/user/errors/example_test.go:115
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
}

func ExampleWrapf() {
	base := errors.Base("whoops")
	err := errors.Wrapf(base, "oh noes #%d", 2)
	fmt.Println(err)
	// Output: oh noes #2
}

func ExampleErrorf() {
	err := errors.Errorf("whoops: %s", "foo")
	fmt.Printf("%+-v", err)

	// Example output:
	// whoops: foo
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleErrorf
	// 	/home/user/errors/example_test.go:165
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
}

func ExampleErrorf_wrap() {
	base := errors.Base("whoops")
	err := errors.Errorf("oh noes (%w)", base)
	fmt.Printf("%+-v", err)

	// Example output:
	// oh noes (whoops)
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleErrorf_wrap
	// 	/home/user/errors/example_test.go:189
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
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

//nolint:dupword
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

func ExampleJoin() {
	err1 := errors.New("error1")
	err2 := errors.New("error2")
	err := errors.Join(err1, err2)
	fmt.Printf("% +-.1v", err)

	// Example output:
	// error1
	// error2
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleJoin
	// 	/home/user/errors/example_test.go:265
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:131
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
	//
	// the above error joins multiple errors:
	//
	// 	error1
	// 	stack trace (most recent call first):
	// 	gitlab.com/tozd/go/errors_test.ExampleJoin
	// 		/home/user/errors/example_test.go:263
	// 	testing.runExample
	// 		/usr/local/go/src/testing/run_example.go:63
	// 	testing.runExamples
	// 		/usr/local/go/src/testing/example.go:44
	// 	testing.(*M).Run
	// 		/usr/local/go/src/testing/testing.go:1927
	// 	main.main
	// 		_testmain.go:131
	// 	runtime.main
	// 		/usr/local/go/src/runtime/proc.go:267
	// 	runtime.goexit
	// 		/usr/local/go/src/runtime/asm_amd64.s:1650
	//
	// 	error2
	// 	stack trace (most recent call first):
	// 	gitlab.com/tozd/go/errors_test.ExampleJoin
	// 		/home/user/errors/example_test.go:264
	// 	testing.runExample
	// 		/usr/local/go/src/testing/run_example.go:63
	// 	testing.runExamples
	// 		/usr/local/go/src/testing/example.go:44
	// 	testing.(*M).Run
	// 		/usr/local/go/src/testing/testing.go:1927
	// 	main.main
	// 		_testmain.go:131
	// 	runtime.main
	// 		/usr/local/go/src/runtime/proc.go:267
	// 	runtime.goexit
	// 		/usr/local/go/src/runtime/asm_amd64.s:1650
}
