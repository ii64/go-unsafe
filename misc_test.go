package unsafelib

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

type testMtype struct {
	//	state         protoimpl.MessageState
	//	sizeCache     protoimpl.SizeCache
	//	unknownFields protoimpl.UnknownFields

	field0 string `mt:"field0,omitempty"`
	Field1 string `protobuf:"bytes,1,opt,name=field1,proto3" json:"field1,omitempty"`
	Field2 string `protobuf:"bytes,2,opt,name=field2,proto3" json:"field2,omitempty"`
}

type testMtyp string

func TestStringCompare(t *testing.T) {
	val := "test"
	val2 := testMtyp(val)
	valHdr := (*String)(unsafe.Pointer(&val))
	valHdr2 := (*String)(unsafe.Pointer(&val2))
	val3 := testMtyp(val)
	valHdr3 := (*String)(unsafe.Pointer(&val3))
	println(valHdr, valHdr2, valHdr3)
	fmt.Println(valHdr, valHdr2, valHdr3)
}

type memwritertest struct {
	byt *byte
	cap int
}

func (w *memwritertest) mprotect() {
}

func (w *memwritertest) WriteString(s string) int {
	for i, b := range s {
		if i > w.cap {
			return i
		}
		*w.byt = byte(b)
		w.byt = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(w.byt)) + 1))
	}
	return len(s)
}

func getStructType(v reflect.Type) *StructType {
	// reflect.Type
	iface := (*emptyInterface)(unsafe.Pointer(&v))
	return (*StructType)(iface.word)
}
func getTyp(v reflect.Type) *rtype {
	// reflect.Type
	iface := (*emptyInterface)(unsafe.Pointer(&v))
	return (*rtype)(iface.word)
}

// switch field
func try(dst any, src any) {
	typDst, typSrc := reflect.TypeOf(dst), reflect.TypeOf(src)
	if typDst.Kind() != reflect.Struct || typSrc.Kind() != reflect.Struct {
		return
	}
	rtDst, rtSrc := getStructType(typDst), getStructType(typSrc)
	var m = make([]structField, len(rtDst.fields))
	for i, f := range rtSrc.fields {
		m[i] = f
	}
	println(rtDst.fields[0].name.bytes)
	rtdstfhdr := (*reflect.SliceHeader)(unsafe.Pointer(&rtDst.fields))
	println(rtdstfhdr)
	fmt.Println(rtdstfhdr, len(rtDst.fields))

	// rodata :/
	mm := (*reflect.SliceHeader)(unsafe.Pointer(&m))
	fmt.Println(mm)
	// rtdstfhdr.Data = mm.Data

	_ = rtDst
	// rtDst.fields = rtSrc.fields
}

func TestRefIfaceTyp(t *testing.T) {
	var m any = testMtype{}

	typ := reflect.TypeOf(m)
	st := getStructType(typ)

	iface := (*emptyInterface)(unsafe.Pointer(&m))
	println(iface.typ, st)
}

func TestSwitchStructFieldMeta(t *testing.T) {

	try(testMtype{}, struct {
		field0 string `asd:"" mt:"field0,omitempty"`
		Field1 string `asd:"" protobuf:"bytes,1,opt,name=field1,proto3" json:"field1,omitempty"`
		Field2 string `asd:"" protobuf:"bytes,2,opt,name=field2,proto3" json:"field2,omitempty"`
	}{})

}

func TestGetStructField(t *testing.T) {
	typ := reflect.TypeOf(struct {
		Mc string `json:"asdf,omitempty"`
	}{})
	rtyp := getStructType(typ)
	fmt.Printf("%+#v %p\n", rtyp, rtyp)

	typ = reflect.TypeOf(testMtype{})
	rtyp = getStructType(typ)
	fmt.Printf("%+#v %p\n", rtyp, rtyp)

	typ = reflect.StructOf([]reflect.StructField{
		{Name: "Mc", Type: reflect.TypeOf(""), Tag: `json:"mc,omitempty"`},
	})
	rtyp = getStructType(typ)
	fmt.Printf("%+#v %p\n", rtyp, rtyp)

}

func TestIfaceTypReplace(t *testing.T) {

	var m any = testMtype{
		field0: "1",
		Field1: "2",
		Field2: "3",
	}
	mhdr := (*emptyInterface)(unsafe.Pointer(&m))
	_ = mhdr

	dojson := func(v any) {
		bb, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", bb)
	}

	dojson(m)

	// create type in heap.
	typ := reflect.TypeOf(m)
	fields := make([]reflect.StructField, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		fields[i] = typ.Field(i)
		newTag := string(fields[i].Tag)
		newTag = strings.ReplaceAll(newTag, `json:"field1,`, `    json:"f1,`)
		newTag = strings.ReplaceAll(newTag, `json:"field2,`, `    json:"f2,`)
		fields[i].Tag = reflect.StructTag(newTag)
	}
	newTyp := reflect.StructOf(fields)
	newHeapTyp := getTyp(newTyp)

	var m2 any
	m2hdr := (*emptyInterface)(unsafe.Pointer(&m2))
	_ = m2hdr
	m2hdr.typ = newHeapTyp
	m2hdr.word = mhdr.word
	fmt.Println(m2hdr)

	for i := 0; i < 5; i++ {
		dojson(m2)
		runtime.GC()
	}
}

func TestTypeTabStructTag(t *testing.T) {
	mf := func(m *testMtype) {
		typ := reflect.TypeOf(testMtype{})

		rtp := ((*emptyInterface)(unsafe.Pointer(&typ))).word
		rt := (*StructType)(rtp)

		fmt.Printf("%+#v %p\n", rt, rt)

		for i := 0; i < typ.NumField(); i++ {
			val := typ.Field(i).Tag
			hdr := (*String)(unsafe.Pointer(&val))

			if false {
				writer := &memwritertest{(*byte)(hdr.Data), hdr.Len}
				mm := strings.ReplaceAll(string(val), `json:"field1,`, `    json:"f1,`)
				writer.WriteString(mm)
			}

			println(hdr)
			fmt.Println(hdr, val)
		}
	}
	mf(new(testMtype))
	mf(new(testMtype))
}
