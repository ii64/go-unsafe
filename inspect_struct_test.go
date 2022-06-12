package unsafelib

import (
	"reflect"
	"testing"
)

func TestInspectStructRecreateType(t *testing.T) {
	type meta struct {
		A string `json:"a"`
		B string `json:"b"`
		C string `json:"c"`
	}

	fields := Inspect(reflect.TypeOf(meta{})).
		Fields()

	field0 := fields[0].Meta()
	field0.Name = "A_etc"
	field0.Tag = field0.Tag + ` rt:"sec,code"`
	fields[0].SetMeta(field0)

	newTyp := fields.CreateStructType()

	val := reflect.New(newTyp)
	valMod := ReinterpretPtr[meta](val.UnsafePointer())
	valMod.A = "abcdef"
	valMod.B = "defgh"
	valMod.C = "ijklmno"

	exp := reflect.ValueOf(valMod)
	for i := 0; i < len(fields); i++ {
		if exp.Elem().Field(0).Interface() != val.Elem().Field(0).Interface() {
			t.Fail()
		}
	}
}
