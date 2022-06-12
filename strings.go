package unsafelib

import (
	"reflect"
	"unsafe"
)

// String2ByteSlice make new slice header from string data pointer, do not grow the size.
func String2ByteSlice(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	v := reflect.SliceHeader{
		Data: uintptr(sh.Data),
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&v))
}

// ByteSlice2String make string from byte slice data pointer.
func ByteSlice2String(s []byte) (ret string) {
	src := ReinterpretPtr[reflect.SliceHeader](&s)
	dst := ReinterpretPtr[reflect.StringHeader](&ret)
	dst.Data = src.Data
	dst.Len = src.Len
	return
}
