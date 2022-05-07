package unsafelib

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestNameBytes(t *testing.T) {
	type test struct {
		K string `json:"helloWorld,omitempty"`
		M string `json:"m,omitempty"`

		C string
	}
	typ := reflectTypeToRtype(reflect.TypeOf(test{}))
	structType := (*StructType)(unsafe.Pointer(typ))

	for i := 0; i < len(structType.fields); i++ {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			n := structType.fields[i].name
			tn := (&TypeName{}).fromName(n)
			fmt.Printf("%+#v\n", tn)

			// if i == 2 {
			// 	tn.Tag = "sdf"
			// }

			res := tn.Bytes() // generate bytes.
			cmp := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
				Data: uintptr(unsafe.Pointer(n.bytes)),
				Len:  len(res),
				Cap:  len(res),
			}))
			t.Run("cmp", func(t *testing.T) {
				if bytes.Compare(res, cmp) != 0 {
					fmt.Printf("%+#v\n%+#v\n", res, cmp)
					t.Fail()
				} else {
					fmt.Printf("%+#v\n%+#v\n", res, cmp)
				}
			})
		})
	}

}
