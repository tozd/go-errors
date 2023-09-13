package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

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
	// the above error joins errors:
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

func ExampleAllDetails() {
	base := errors.Base("not found")
	err1 := errors.WithDetails(base, "file", "plans.txt")
	err2 := errors.WithDetails(err1, "user", "vader")
	fmt.Println(errors.AllDetails(err1))
	fmt.Println(errors.AllDetails(err2))
	// Output:
	// map[file:plans.txt]
	// map[file:plans.txt user:vader]
}

func ExampleDetails() {
	base := errors.Base("not found")
	err := errors.WithStack(base)
	errors.Details(err)["file"] = "plans.txt"
	errors.Details(err)["user"] = "vader"
	fmt.Println(errors.Details(err))
	// Output:
	// map[file:plans.txt user:vader]
}

func ExampleWithDetails_printf() {
	base := errors.Base("not found")
	err := errors.WithDetails(base, "file", "plans.txt", "user", "vader")
	fmt.Printf("%#v", err)
	// Output:
	// not found
	// file=plans.txt
	// user=vader
}

func ExampleStackFormatter_MarshalJSON() {
	const depth = 1
	var cs [depth]uintptr
	runtime.Callers(1, cs[:])
	data, err := json.Marshal(errors.StackFormatter{cs[:]})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// Example output:
	// [{"name":"gitlab.com/tozd/go/errors_test.ExampleStackFormatter_MarshalJSON","file":"/home/user/errors/example_test.go","line":360}]
}

func ExampleStackFormatter_Format() {
	const depth = 1
	var cs [depth]uintptr
	runtime.Callers(1, cs[:])
	fmt.Printf("%+v", errors.StackFormatter{cs[:]})

	// Example output:
	// gitlab.com/tozd/go/errors_test.ExampleStackFormatter_Format
	// 	/home/user/errors/example_test.go:374
}

func ExampleStackFormatter_Format_width() {
	const depth = 1
	var cs [depth]uintptr
	runtime.Callers(1, cs[:])
	fmt.Printf("%+2v", errors.StackFormatter{cs[:]})

	// Example output:
	// gitlab.com/tozd/go/errors_test.ExampleStackFormatter_Format
	//   /home/user/errors/example_test.go:385
}

func ExampleFormatter_Format() {
	base := errors.Base("not found")
	err := errors.Wrap(base, "image not found")
	errors.Details(err)["filename"] = "star.png"
	fmt.Printf("% #+-.1v", err)

	// Example output:
	// image not found
	// filename=star.png
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleFormatter_Format
	// 	/home/user/errors/example_test.go:395
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:137
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
	//
	// the above error was caused by the following error:
	//
	// not found
}

func ExampleFormatter_MarshalJSON() {
	base := errors.Base("not found")
	errE := errors.Wrap(base, "image not found")
	errors.Details(errE)["filename"] = "star.png"
	data, err := json.Marshal(errE)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	_ = json.Indent(&out, data, "", "\t")
	_, _ = out.WriteTo(os.Stdout)

	// Example output:
	// {
	// 	"cause": {
	// 		"error": "not found"
	// 	},
	// 	"error": "image not found",
	// 	"filename": "star.png",
	// 	"stack": [
	// 		{
	// 			"name": "gitlab.com/tozd/go/errors_test.ExampleFormatter_MarshalJSON",
	// 			"file": "/home/user/errors/example_test.go",
	// 			"line": 427
	// 		},
	// 		{
	// 			"name": "testing.runExample",
	// 			"file": "/usr/local/go/src/testing/run_example.go",
	// 			"line": 63
	// 		},
	// 		{
	// 			"name": "testing.runExamples",
	// 			"file": "/usr/local/go/src/testing/example.go",
	// 			"line": 44
	// 		},
	// 		{
	// 			"name": "testing.(*M).Run",
	// 			"file": "/usr/local/go/src/testing/testing.go",
	// 			"line": 1927
	// 		},
	// 		{
	// 			"name": "main.main",
	// 			"file": "_testmain.go",
	// 			"line": 145
	// 		},
	// 		{
	// 			"name": "runtime.main",
	// 			"file": "/usr/local/go/src/runtime/proc.go",
	// 			"line": 267
	// 		},
	// 		{
	// 			"name": "runtime.goexit",
	// 			"file": "/usr/local/go/src/runtime/asm_amd64.s",
	// 			"line": 1650
	// 		}
	// 	]
	// }
}

func ExampleUnmarshalJSON() {
	base := errors.Base("not found")
	errE := errors.Wrap(base, "image not found")
	errors.Details(errE)["filename"] = "star.png"
	data, err := json.Marshal(errE)
	if err != nil {
		panic(err)
	}
	errFromJSON, errE := errors.UnmarshalJSON(data)
	if errE != nil {
		panic(errE)
	}
	fmt.Printf("% #+-.1v", errFromJSON)

	// Example output:
	// image not found
	// filename=star.png
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleUnmarshalJSON
	// 	/home/user/errors/example_test.go:486
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:145
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
	//
	// the above error was caused by the following error:
	//
	// not found
}

func ExampleWrapWith() {
	cause := io.EOF
	baseErr := errors.Base("an error")

	err := errors.WrapWith(cause, baseErr)
	fmt.Println(err)
	// Output: an error
}

func ExampleWrapWith_printf() {
	cause := io.EOF
	baseErr := errors.Base("an error")

	err := errors.WrapWith(cause, baseErr)
	fmt.Printf("% +-.1v", err)

	// Example output:
	// an error
	// stack trace (most recent call first):
	// gitlab.com/tozd/go/errors_test.ExampleWrapWith_printf
	// 	/home/user/errors/example_test.go:536
	// testing.runExample
	// 	/usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	// 	/usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	// 	/usr/local/go/src/testing/testing.go:1927
	// main.main
	// 	_testmain.go:155
	// runtime.main
	// 	/usr/local/go/src/runtime/proc.go:267
	// runtime.goexit
	// 	/usr/local/go/src/runtime/asm_amd64.s:1650
	//
	// the above error was caused by the following error:
	//
	// EOF
}
