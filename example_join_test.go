package errors_test

import (
	"fmt"
	"os"

	"gitlab.com/tozd/go/errors"
)

func run() (errE errors.E) {
	file, err := os.CreateTemp("", "test")
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		errE = errors.Join(errE, errors.WithStack(os.Remove(file.Name()))) //nolint:gosec
	}()

	// Do something with the file...

	return nil
}

func ExampleJoin_defer() {
	errE := run()
	if errE != nil {
		fmt.Printf("error: %+v\n", errE)
	} else {
		fmt.Printf("success\n")
	}
	// Output:
	// success
}
