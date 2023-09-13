//go:build !go1.21
// +build !go1.21

package errors

import (
	"reflect"
	"unsafe"
)

func slicesEqual(a []uintptr, b []uintptr) bool {
	ah := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return ah.Len == bh.Len && ah.Data == bh.Data
}
