package errors_test

import (
	"fmt"

	"gitlab.com/tozd/go/errors"
)

type MyError struct {
	code    int
	message string
}

func (e MyError) Error() string {
	return e.message
}

func (e MyError) Code() int {
	return e.code
}

var (
	ErrBadRequest = &MyError{400, "error"}
	ErrNotFound   = &MyError{404, "not found"}
)

func getMyErr() error {
	return ErrNotFound
}

func ExampleAs() {
	err := getMyErr()

	var myErr *MyError
	if errors.As(err, &myErr) {
		fmt.Printf("code: %d", myErr.Code())
	}
	// Output:
	// code: 404
}
