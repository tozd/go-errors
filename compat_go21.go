//go:build go1.21
// +build go1.21

package errors

import (
	"unsafe"
)

func slicesEqual(a []uintptr, b []uintptr) bool {
	return len(a) == len(b) && unsafe.SliceData(a) == unsafe.SliceData(b)
}
