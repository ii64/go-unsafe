package unsafelib

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestInspectTransformStruct(t *testing.T) {
	type TestTransformStruct struct {
		A string `json:"a,omitempty" desc:"desc here"`
		B string `json:"b,omitempty" desc:"desc here"`
		C string `json:"c,omitempty" desc:"desc here"`
		D [4]int `json:"d,omitempty" desc:"desc here"`
		E int    `json:"e,omitempty" desc:"desc here"`
		F []int  `json:"f,omitempty" desc:"desc here"`
		G struct {
			Sub string
		} `json:"xsdf"`
	}

	typ := reflect.TypeOf(TestTransformStruct{})
	Inspect(typ).Change(TransformStructField(func(field *StructField, name *TypeName) {
		fmt.Printf("%+#v %+#v\n", field, name)

		tag := string(name.Tag)
		tag = strings.ReplaceAll(tag, `json:"`, `json:"inspected_`+name.Name+"_")

		typ := rtypeToReflectType(field.typ)
		tag = tag + fmt.Sprintf(` runtime_annotation:"name=%s,type=%s,alignment=%d,size=%d,json_tag=%s"`,
			name.Name,
			typ.Name(),
			typ.Align(),
			typ.Size(),
			name.Tag.Get("json"),
		)

		name.Tag = reflect.StructTag(tag)
	}))

	t.Run("typeof", func(t *testing.T) {
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			rt := f.Tag.Get("runtime_annotation")
			if rt == "" {
				t.Fail()
			}
			fmt.Println(rt)
		}
	})

	t.Run("json", func(t *testing.T) {
		bb, err := json.Marshal(&TestTransformStruct{
			A: "hello",
			B: "world",
			C: "yeah",
		})
		if err != nil {
			t.Fail()
		}
		fmt.Printf("%s\n", bb)
	})

}
