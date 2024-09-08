package main

import (
	"gitlab.com/tozd/go/errors"
)

func main() {
	panic(errors.WithDetails(errors.Base("panic error"), "key", "value"))
}
