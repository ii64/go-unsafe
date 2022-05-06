package unsafelib

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func init() {
	// debug.SetPanicOnFault(true)
}

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
	type tdst = map[string]String
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
	type tdst = map[string]Interface
	println(unsafe.Sizeof(tsrc{}), unsafe.Sizeof(tdst{}))

	var src any = tsrc{
		"hey":   []int64{1, 2, 3, 4, 5},
		"hello": []uint64{1, 2, 3, 4, 5},
	}
	typMapSliceHdr := reflect.TypeOf(tdst{})
	var dst = CastInterface(src, typMapSliceHdr).(tdst)

	for key, srcVal := range src.(tsrc) {
		dstVal := dst[key]
		dstCast := (*reflect.SliceHeader)(unsafe.Pointer(dstVal.Word))

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

//

type testStruct1 struct {
	K int
}

func (m *testStruct1) X() int {
	m.K++
	return 0xfa
}
func (m *testStruct1) Y() int { return 0xfb }

func TestInterfaceCastImplementer(t *testing.T) {

	type testInterfaceImplementer interface {
		X() int
	}

	type testInterfaceImplementer2 interface {
		testInterfaceImplementer
		Y() int
	}

	orig := &testStruct1{123}
	var pl testInterfaceImplementer2

	var pl2 testInterfaceImplementer2 = &testStruct1{}
	var pl1 testInterfaceImplementer = &testStruct1{}

	// typ := reflect.TypeOf(pl1) // once `pl1` goes to reflect.TypeOf, it will be referenced as empty interface.
	mt := CastInterfacePtr(&pl, orig, (*Interface)(unsafe.Pointer(&pl1)).Type) // do this instead.

	// fmt.Printf("el %+#v\n", reflect.TypeOf(&pl).Elem().Elem().Kind().String())

	// --

	// var im any = orig
	// efo := (*ifacetyp)(unsafe.Pointer(&im))
	// fmt.Printf("* %+#v %p\n", (*rtype)(efo.typ), efo.typ)

	efc := (*Interface)(unsafe.Pointer(&pl))
	fmt.Printf("pl %+#v %p\n", (*rtype)(efc.Type), efc.Type)

	efc1 := (*Interface)(unsafe.Pointer(&pl1))
	fmt.Printf("pl1 %+#v %p\n", (*rtype)(efc1.Type), efc1.Type)

	// efc2 := (*ifacetyp)(unsafe.Pointer(&pl2))
	// fmt.Printf("* %+#v %p\n", (*rtype)(efc2.typ), efc2.typ)

	// --

	println(pl, pl2, mt)

	fmt.Printf("pl %+#v %T\n", pl, pl)
	fmt.Printf("mt %+#v %T\n", mt, mt)

	// -- check iface method lookup

	pl.X()
	fmt.Printf("pl post-call %+#v %+#v\n", pl, orig)

	// DON'T CALL Y() method !!
	// pl.Y() // SEGV, testInterfaceImplementer (pl1) casted to testInterfaceImplementer2 (pl)

}

func TestInterfaceCastNoType(t *testing.T) {

	type impl interface {
		X() int
	}

	type ma struct {
		X int
	}

	obj := &ma{0xfa}
	var x impl

	// invalid rtype
	_obj := (any)(obj)
	_iface := (*Interface)(unsafe.Pointer(&_obj))

	fmt.Printf("%+#v %+#v\n", _iface, _iface.Type)

	cast := (*Interface)(unsafe.Pointer(&x))
	typ := *_iface.Type
	typ.ptrdata = uintptr(unsafe.Pointer(_iface.Type))
	// cast.Type = &typ // need itab
	cast.Word = unsafe.Pointer(obj)

	// x.X()
}

// ---

func TestInterfaceCastChangePtrData(t *testing.T) {
	type testModelNoMethod struct {
		X [8]byte
	}
	type testInterfaceImplementer interface {
		X() int
	}
	m := &testStruct1{}
	var x testInterfaceImplementer = m

	x.X()
	x.X()
	x.X()

	var targ = &testModelNoMethod{}
	ChangeInterfacePtrData(&x, targ)

	x.X()
	x.X()
	x.X()

	fmt.Printf("%+#v\n%+#v\n", m, targ)
}

// ---

func TestInterfaceCastMethodMod(t *testing.T) {
	type testInterfaceImplementer interface {
		X() int
	}
	var x testInterfaceImplementer = &testStruct1{}
	val := reflect.ValueOf(x)

	for i := 0; i < val.NumMethod(); i++ {
		m := val.Method(i)
		fmt.Printf("%+#v %+#v\n", m, m.Bytes())
	}
}
