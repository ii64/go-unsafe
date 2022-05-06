package unsafelib

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type m[T comparable] struct {
	Value T
}

func TestGenericCastStruct(t *testing.T) {
	type ma struct {
		X int64
	}
	type mb struct {
		X [8]byte
	}
	type mc struct {
		X uint64
	}
	a := &ma{0xffffff}
	var b *mb
	CastPtr(&b, a)
	var c = ReinterpretPtr[mc](b)
	fmt.Printf("%+#v\n%+#v\n%+#v\n", a, b, c)

	t.Run("cmp", func(t *testing.T) {
		if uint64(a.X) != c.X {
			t.Fail()
		}
	})
}

func TestGenericCastCopy(t *testing.T) {
	type ma struct {
		A int64
		B uint64
	}
	type mb struct {
		A uint64
		B int64
	}
	var a ma = ma{0xff, 0xfe}
	b := *ReinterpretPtr[mb](&a)
	fmt.Printf("%+#v\n%+#v\n", a, b)
	a.A = 0xca
	b.A = uint64(a.A)
	fmt.Printf("%+#v\n%+#v\n", a, b)

	t.Run("cmp", func(t *testing.T) {
		t.Run("cmp-1", func(t *testing.T) {
			if uint64(a.A) != b.A {
				t.Fail()
			}
		})
		t.Run("cmp-2", func(t *testing.T) {
			if int64(a.B) != b.B {
				t.Fail()
			}
		})
	})
}

func TestGenericCastMap(t *testing.T) {
	orig := map[string]any{
		"hey":   1,
		"hello": "1234",
		"hi":    "1234",
	}
	casted := *ReinterpretPtr[map[string]Interface](&orig)
	fmt.Printf("%+#v\n", orig)
	fmt.Printf("%+#v\n", casted)
	t.Run("cmp", func(t *testing.T) {
		for k, v := range orig {
			if v == nil {
				continue
			}
			vx := reflect.ValueOf(v)
			vword := (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1]
			castVal := casted[k]
			t.Run("cmp1-"+k, func(t *testing.T) {
				if rtypeToReflectType(castVal.Type) != vx.Type() {
					t.Fail()
				}
			})
			t.Run("cmp2-"+k, func(t *testing.T) {
				if castVal.Word != vword {
					t.Fail()
				}
			})
		}
	})
}

func TestGenericCastSlice(t *testing.T) {
	orig := [][]string{{"123"}, {"456", "789"}, {}, nil}
	casted := *ReinterpretPtr[[]reflect.SliceHeader](&orig)
	fmt.Printf("%+#v\n", orig)
	fmt.Printf("%+#v\n", casted)
	t.Run("cmp", func(t *testing.T) {
		for i, src := range orig {
			cast := casted[i]
			t.Run(fmt.Sprintf("cmp-%d", i), func(t *testing.T) {
				if len(src) != cast.Len || cap(src) != cast.Cap {
					t.Fail()
				}
				if len(src) != 0 && uintptr(unsafe.Pointer(&src[0])) != cast.Data {
					t.Fail()
				}
			})
		}
	})
}

func TestGenericInvalidPtr(t *testing.T) {
	v := ReinterpretPtr[*int](0xff) // don't dereference !
	_ = v
	// println(v)
}

// ----

func TestX(t *testing.T) {
	type impl interface {
		X() int
	}
	var a impl = &testStruct1{
		K: 1234,
	}
	genericFun(a)
}

/*
Go1.18.1-linux-amd64
	LEAQ    type."".testStruct1(SB), AX
	PCDATA  $1, $0
	NOP
	CALL    runtime.newobject(SB) // allocate a new object
	LEAQ    go.itab.*"".testStruct1,"".impl·1(SB), BX
	MOVQ    AX, CX
	LEAQ    ""..dict.genericFun["".impl·1](SB), AX
	CALL    "".genericFun[go.shape.interface { X() int }_0](SB)
*/
func genericFun[T any](v T) {
	// AX, BX, CX, DI, SI, R8, R9, R10, R11, ... n(SP)

	// AX: ""..dict.Mmc["".impl·1](SB)
	// BX: go.itab.*"".testStruct1,"".impl·1(SB)

	args := (*[10]unsafe.Pointer)(unsafe.Pointer(&v))[:]
	fmt.Printf("%+#v\n", args)

	_ax := (*rtype)(args[0])
	ax := rtypeToReflectType(_ax)
	cx := rtypeToReflectType((*rtype)(args[2]))
	di := rtypeToReflectType((*rtype)(args[3]))

	ptyp := rtypeToReflectType((*rtype)(unsafe.Pointer(_ax.ptrdata)))

	fmt.Printf("%+#v %+#v %+#v %+#v\n",
		ax.Kind().String(),
		cx.Kind().String(),
		di.Kind().String(),
		ptyp.Kind().String())

	fmt.Printf("%+#v\n", ax)                      // AX
	fmt.Printf("%+#v\n", (*testStruct1)(args[1])) // BX
	fmt.Printf("%+#v\n", cx)                      // CX
	fmt.Printf("%+#v\n", di)                      // DI

}
