package unsafelib

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestString2ByteSlice(t *testing.T) {
	src := "hello world\n"
	bs := String2ByteSlice(src)

	dst := string(bs) // on the heap.

	srch := (*String)(unsafe.Pointer(&src))
	dsth := (*String)(unsafe.Pointer(&dst))

	fmt.Printf("%+#v\n%+#v\n", srch, dsth)
}
