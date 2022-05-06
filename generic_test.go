package unsafelib

import (
	"fmt"
	"testing"
)

type m[T comparable] struct {
	Value T
}

func TestGenericShapes(t *testing.T) {
	type ma struct {
		X int64
	}
	type mb struct {
		X [8]byte
	}

	a := &ma{0xffffff}
	var b *mb
	CastPtr(&b, a)

	var c = ReinterpretPtr[mb](a)

	fmt.Printf("%+#v\n%+#v\n%+#v\n", a, b, c)
}
