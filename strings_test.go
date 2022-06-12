package unsafelib

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestString2ByteSlice(t *testing.T) {
	src := "hello world\n"
	bs := String2ByteSlice(src)

	dst := string(bs) // on the heap.

	srch := ReinterpretPtr[reflect.StringHeader](&src)
	dsth := ReinterpretPtr[reflect.StringHeader](&dst)

	fmt.Printf("%+#v\n%+#v\n", srch, dsth)

	t.Run("cmp", func(t *testing.T) {
		if strings.Compare(src, dst) != 0 {
			t.Fail()
		}
	})
}

func TestByteSlice2String(t *testing.T) {
	src := []byte("hello world\n")
	bs := ByteSlice2String(src)
	dst := []byte(bs)

	srch := *ReinterpretPtr[[]byte](&src)
	dsth := *ReinterpretPtr[[]byte](&dst)

	fmt.Printf("%+#v\n%+#v\n", srch, dsth)

	t.Run("cmp", func(t *testing.T) {
		if bytes.Compare(srch, dsth) != 0 {
			t.Fail()
		}
	})
}
