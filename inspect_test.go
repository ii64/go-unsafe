package unsafelib

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInspect(t *testing.T) {
	ts := []reflect.Type{
		reflect.TypeOf(int64(123)),
		reflect.TypeOf("sdf"),
		reflect.TypeOf(testStruct1{}),
		reflect.StructOf([]reflect.StructField{
			{
				Name: "TestX",
				Type: reflect.TypeOf("sdf"),
				Tag:  `json:"TestX"`,
			},
		}),

		// invalid:
		// reflect.TypeOf(new(int)),
		// nil,
	}

	for _, ts := range ts {
		Inspect(ts).Change(
			func(insp *Inspection) bool {
				// println(insp, insp.writable, insp.Type().Kind().String())
				if insp.Type().Kind() == reflect.Struct {
					k := insp.Fields()
					m := ReinterpretPtr[reflect.SliceHeader](&k)
					fmt.Printf("%+#v %+#v\n", k, m)
				}
				return true
			},
		)

		fmt.Println("represents: ", ts.String())
	}

}
