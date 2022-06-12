package unsafelib

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestMprotect(t *testing.T) {
	var src any = &testStruct1{
		K: 12345,
	}
	iface := (*Interface)(unsafe.Pointer(&src))
	typ := rtypeToReflectType(iface.Type).Elem()
	structType := ReinterpretPtr[StructType](reflectTypeToRtype(typ))
	fmt.Printf("st: %+#v\n", structType)

	var sfs []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		sfs = append(sfs, typ.Field(i))
	}
	fmt.Printf("fields: %+#v\n", sfs)

	// change fields tag
	for i, sf := range sfs {
		tag := string(sf.Tag)
		tag = strings.ReplaceAll(tag, `json:"`, `json:"replaced_`)
		sfs[i].Tag = reflect.StructTag(tag)
	}

	// create new type
	newTyp := reflect.StructOf(sfs)
	newStructType := ReinterpretPtr[StructType](reflectTypeToRtype(newTyp))
	fmt.Printf("newSt: %+#v\n", newStructType)

	t.Run("mprotect", func(t *testing.T) {
		prf := new_mem_profile(uintptr(unsafe.Pointer(structType)),
			1, PROT_FLAG_RO, PROT_FLAG_RW)
		prf.toggle()
		defer prf.toggle()

		// write fields meta
		structType.fields = newStructType.fields
	})

	fmt.Printf("oldType represents: %+#v\n  %+#v\n", typ.String(), structType)
	fmt.Printf("newType represents: %+#v\n  %+#v\n", newTyp.String(), newStructType)

	//
	t.Run("json-encode", func(t *testing.T) {
		bb, err := json.Marshal(src)
		if err != nil {
			t.Fail()
		}
		fmt.Printf("output type: %T, json: %s\n", src, bb)
	})

}
