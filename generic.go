package unsafelib

import (
	"unsafe"
)

// CastPtr reinterpret ptr to perform type conversion.
func CastPtr[T any, U any](dst **T, src U) {
	*dst = *(**T)(unsafe.Pointer(&src))
}

// ReinterpretPtr reinterpret ptr to perform type conversion.
func ReinterpretPtr[T any, U any](src U) *T {
	return *(**T)(unsafe.Pointer(&src))
}
