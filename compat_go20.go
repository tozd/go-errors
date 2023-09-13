//go:build go1.20
// +build go1.20

package errors

import (
	"fmt"
	"unsafe"
)

var formatString = fmt.FormatString //nolint:gochecknoglobals

func slicesEqual(a []uintptr, b []uintptr) bool {
	return len(a) == len(b) && unsafe.SliceData(a) == unsafe.SliceData(b)
}
