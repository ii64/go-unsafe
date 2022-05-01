package unsafelib

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes(t *testing.T) {
	type a struct{}
	type b a

	type malias map[string]any
	type aalias []string
	type balias []string

	fmt.Printf("%+#v\n", reflect.TypeOf(a{}))
	fmt.Printf("%+#v\n", reflect.TypeOf(b{}))

	fmt.Printf("%+#v\n", reflect.TypeOf(map[string]any{}))
	fmt.Printf("%+#v\n", reflect.TypeOf(malias{}))
	fmt.Printf("%+#v\n", reflect.TypeOf([]string{}))
	fmt.Printf("%+#v\n", reflect.TypeOf(aalias{}))
	fmt.Printf("%+#v\n", reflect.TypeOf(balias{}))

}
