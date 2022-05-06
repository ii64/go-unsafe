package unsafelib

import (
	"reflect"
	"unsafe"
)

func String2ByteSlice(s string) []byte {
	sh := (*String)(unsafe.Pointer(&s))
	v := reflect.SliceHeader{
		Data: uintptr(sh.Data),
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&v))
}
