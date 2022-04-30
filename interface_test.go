package unsafelib

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ii64/go-unsafe/unsafeheader"
)

func TestInterfaceNil(t *testing.T) {
	println(reflect.TypeOf((any)(nil)))
}

func TestInterfaceCastStruct(t *testing.T) {
	type implementer interface {
	}
	type m1 struct {
		Value string `json:"value"`
	}
	type m2 struct {
		Value string `json:"v"`
	}
	m2Type := reflect.TypeOf(m2{})
	m2ptrType := reflect.TypeOf(&m2{})
	_, _ = m2Type, m2ptrType
	var v implementer = &m1{Value: "sdw"}
	var v2 implementer = CastInterface(v, m2Type)
	var v3 implementer = CastInterface(v, m2ptrType)
	if _, ok := v2.(m2); !ok {
		t.Fail()
	}
	if _, ok := v3.(*m2); !ok {
		t.Fail()
	}
	fmt.Printf("%+#v\n", v)
	fmt.Printf("%+#v\n", v2)
	fmt.Printf("%+#v\n", v3)
}

func TestInterfaceCastSlice(t *testing.T) {
	type tsrc = []byte
	type tdst = []int8
	println(unsafe.Sizeof(tsrc{}), unsafe.Sizeof(tdst{}))
	var src any = tsrc{1, 3, 4, 5, 0xca, 0xfe, 0xba, 0xbe}
	typInt8Slice := reflect.TypeOf(tdst{})
	var dst = CastInterface(src, typInt8Slice).(tdst)
	for i, srcEl := range src.(tsrc) {
		var dstEl = dst[i]
		if srcEl != byte(dstEl) {
			t.Fail()
		}
	}
	fmt.Println(src)
	fmt.Println(dst)
}

func TestInterfaceCastMapString(t *testing.T) {
	type tsrc = map[string]string
	type tdst = map[string]unsafeheader.String
	println(unsafe.Sizeof(tsrc{}), unsafe.Sizeof(tdst{}))

	var src any = tsrc{
		"hey":   "hey",
		"hello": "hello",
	}
	typMapStrHdr := reflect.TypeOf(tdst{})
	var dst = CastInterface(src, typMapStrHdr).(tdst)

	for key, srcVal := range src.(tsrc) {
		dstVal := dst[key]
		dstCast := *(*string)(unsafe.Pointer(&dstVal)) // cast back to string
		if len(srcVal) != dstVal.Len {                 // match length
			t.Fail()
		}
		if dstCast != srcVal { // match as string
			t.Fail()
		}
	}
	fmt.Println(src)
	fmt.Println(dst)
}

func TestInterfaceCastMapIface(t *testing.T) {
	type tsrc = map[string]any
	type tdst = map[string]ifacetyp
	println(unsafe.Sizeof(tsrc{}), unsafe.Sizeof(tdst{}))

	var src any = tsrc{
		"hey":   []int64{1, 2, 3, 4, 5},
		"hello": []uint64{1, 2, 3, 4, 5},
	}
	typMapSliceHdr := reflect.TypeOf(tdst{})
	var dst = CastInterface(src, typMapSliceHdr).(tdst)

	for key, srcVal := range src.(tsrc) {
		dstVal := dst[key]
		dstCast := (*reflect.SliceHeader)(unsafe.Pointer(dstVal.word))

		switch srcCast := srcVal.(type) {
		case []int64:
			if len(srcCast) != dstCast.Len || cap(srcCast) != dstCast.Cap {
				t.Fail()
			} else if unsafe.Pointer(&srcCast[0]) != unsafe.Pointer(dstCast.Data) {
				t.Fail()
			}
		case []uint64:
			if len(srcCast) != dstCast.Len || cap(srcCast) != dstCast.Cap {
				t.Fail()
			} else if unsafe.Pointer(&srcCast[0]) != unsafe.Pointer(dstCast.Data) {
				t.Fail()
			}
		}

	}

	fmt.Println(src)
	fmt.Println(dst)
}
