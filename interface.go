package unsafelib

import (
	"reflect"
	"unsafe"
)

type ifacetyp struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

func CastInterface(src interface{}, typ reflect.Type) interface{} {
	iface := (*ifacetyp)(unsafe.Pointer(&src))
	ifaceTyp := (*ifacetyp)(unsafe.Pointer(&typ))
	iface.typ = ifaceTyp.word
	return src
}
