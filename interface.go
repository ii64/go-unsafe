package unsafelib

import (
	"reflect"
	"unsafe"
)

type Interface struct {
	Type *rtype
	Word unsafe.Pointer
}

var (
	// reflect.TypeOf (*rtype)
	typInterfaceHdr           = reflect.TypeOf(Interface{})
	TypReflectTypeOf_rtypePtr = (*Interface)(unsafe.Pointer(&typInterfaceHdr)).Type
)

// rtypeToReflectType wraps `*rtype` with `reflect.Type & *reflect.rtype` itab.
func rtypeToReflectType(typ *rtype) (v reflect.Type) {
	ret := (*Interface)(unsafe.Pointer(&v))
	ret.Type = TypReflectTypeOf_rtypePtr
	ret.Word = unsafe.Pointer(typ)
	return
}

func reflectTypeToRtype(typ reflect.Type) *rtype {
	iface := (*Interface)(unsafe.Pointer(&typ))
	return (*rtype)(iface.Word)
}

// ChangeInterfacePtrData change ptr data, leaving type referencec untouched.
//   dst: *(any)(nil/obj) - must NOT be nil
//   src: (any)
// See also: CastInterface, CastInterfacePtr
func ChangeInterfacePtrData(dst any, src any) {
	if dst == nil || src == nil {
		return
	}
	ifaceDst := (*Interface)(unsafe.Pointer(&dst))
	iface := (*Interface)(ifaceDst.Word)

	ifaceSrc := (*Interface)(unsafe.Pointer(&src))
	iface.Word = ifaceSrc.Word
}

// CastInterface casts interface type
// See also: ChangeInterfacePtrData
func CastInterface(src any, typ reflect.Type) any {
	iface := *(*Interface)(unsafe.Pointer(&src))
	ifaceTyp := (*Interface)(unsafe.Pointer(&typ))
	iface.Type = (*rtype)(ifaceTyp.Word)
	return *(*any)(unsafe.Pointer(&iface))
}

// CastInterfacePtr dst is a pointer to a interface
//   dst: *(any)(nil/obj)
//   src: (any)(obj)
//   typ: (any)(*rtype)
// See also: ChangeInterfacePtrData
func CastInterfacePtr(dst any, src any, typ any) (ret any) {
	if dst == nil || src == nil || typ == nil {
		return
	}

	// dst any
	ifaceDst := (*Interface)(unsafe.Pointer(&dst))
	if ifaceDst == nil {
		return
	}
	//*dst
	iface := (*Interface)(ifaceDst.Word)

	// src any
	ifaceSrc := (*Interface)(unsafe.Pointer(&src))
	if ifaceSrc == nil || ifaceSrc.Word == nil { // source data ptr
		return
	}

	// typ any
	ifaceTypRef := (*Interface)(unsafe.Pointer(&typ))
	if ifaceTypRef == nil || ifaceTypRef.Word == nil {
		return
	}
	rtyp := (*rtype)(ifaceTypRef.Word)

	// fmt.Printf("%+#v %+#v\n", reflect.ValueOf(dst).Elem().UnsafeAddr(), iface.word)

	// fmt.Printf("%+#v\n", (*rtype)(ifaceDst.typ))  // dst typ
	// fmt.Printf("* %+#v\n", (*rtype)(iface.typ)) // dst data typ
	// fmt.Printf("%+#v\n", (*rtype)(ifaceSrc.typ)) // src typ
	// fmt.Printf("%+#v\n", (*rtype)(ifaceTypRef.word)) // typ val

	// fmt.Printf("ifaceDst %+#v\n", ifaceDst)
	// fmt.Printf("iface %+#v\n", iface)
	// fmt.Printf("ifaceTypRef %+#v %+#v\n", ifaceTypRef, ifaceTypRef.Type)
	// fmt.Printf("rtyp %+#v\n", rtyp)

	iface.Type = rtyp
	iface.Word = ifaceSrc.Word

	ifaceRet := (*Interface)(unsafe.Pointer(&ret))

	// check ptrdata
	ifaceElTyp := (*rtype)(unsafe.Pointer(rtyp.ptrdata))
	elTyp := rtypeToReflectType(ifaceElTyp)
	if elTyp.Kind() != reflect.Invalid { // interface implements method
		ifaceRet.Type = ifaceElTyp
		ifaceRet.Word = iface.Word
	} else {
		*ifaceRet = *iface
	}
	return
}
