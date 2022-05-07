// Source reflect/type.go
package unsafelib

import "unsafe"

// tflag is used by an rtype to signal what extra type information is
// available in the memory directly following the rtype value.
//
// tflag values must be kept in sync with copies in:
//	cmd/compile/internal/reflectdata/reflect.go
//	cmd/link/internal/ld/decodesym.go
//	runtime/type.go
type tflag uint8

type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype
type textOff int32 // offset from top of text section

type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte   // garbage collection data
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}

//

var (
	RTYPE_SIZE      = unsafe.Sizeof(rtype{})
	STRUCTTYPE_SIZE = unsafe.Sizeof(StructType{})
)

//

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}

//

// Struct field
type StructField struct {
	name        name    // name is always non-empty
	typ         *rtype  // type of field
	offsetEmbed uintptr // byte offset of field<<1 | isEmbedded
}

// StructType represents a struct type. (previously not exported.)
type StructType struct {
	rtype
	pkgPath name
	fields  []StructField // sorted by offset
}

// name is an encoded type name with optional extra data.
//
// The first byte is a bit field containing:
//
//	1<<0 the name is exported
//	1<<1 tag data follows the name
//	1<<2 pkgPath nameOff follows the name and tag
//
// Following that, there is a varint-encoded length of the name,
// followed by the name itself.
//
// If tag data is present, it also has a varint-encoded length
// followed by the tag itself.
//
// If the import path follows, then 4 bytes at the end of
// the data form a nameOff. The import path is only set for concrete
// methods that are defined in a different package than their type.
//
// If a name starts with "*", then the exported bit represents
// whether the pointed to type is exported.
//
// Note: this encoding must match here and in:
//   cmd/compile/internal/reflectdata/reflect.go
//   runtime/type.go
//   internal/reflectlite/type.go
//   cmd/link/internal/ld/decodesym.go

type name struct {
	bytes *byte
}

// Value is the reflection interface to a Go value.
//
// Not all methods apply to all kinds of values. Restrictions,
// if any, are noted in the documentation for each method.
// Use the Kind method to find out the kind of value before
// calling kind-specific methods. Calling a method
// inappropriate to the kind of type causes a run time panic.
//
// The zero Value represents no value.
// Its IsValid method returns false, its Kind method returns Invalid,
// its String method returns "<invalid Value>", and all other methods panic.
// Most functions and methods never return an invalid value.
// If one does, its documentation states the conditions explicitly.
//
// A Value can be used concurrently by multiple goroutines provided that
// the underlying Go value can be used concurrently for the equivalent
// direct operations.
//
// To compare two Values, compare the results of the Interface method.
// Using == on two Values does not compare the underlying values
// they represent.
type Value struct {
	// typ holds the type of the value represented by a Value.
	typ *rtype

	// Pointer-valued data or, if flagIndir is set, pointer to data.
	// Valid when either flagIndir is set or typ.pointers() is true.
	ptr unsafe.Pointer

	// flag holds metadata about the value.
	// The lowest bits are flag bits:
	//	- flagStickyRO: obtained via unexported not embedded field, so read-only
	//	- flagEmbedRO: obtained via unexported embedded field, so read-only
	//	- flagIndir: val holds a pointer to the data
	//	- flagAddr: v.CanAddr is true (implies flagIndir)
	//	- flagMethod: v is a method value.
	// The next five bits give the Kind of the value.
	// This repeats typ.Kind() except for method values.
	// The remaining 23+ bits give a method number for method values.
	// If flag.kind() != Func, code can assume that flagMethod is unset.
	// If ifaceIndir(typ), code can assume that flagIndir is set.
	flag

	// A method value represents a curried method invocation
	// like r.Read for some receiver r. The typ+val+flag bits describe
	// the receiver r, but the flag's Kind bits say Func (methods are
	// functions), and the top bits of the flag give the method number
	// in r's type's method table.
}

type flag uintptr

//go:linkname reflect_fnv1 reflect.fnv1
//go:noescape
func reflect_fnv1(x uint32, list ...byte) uint32

//go:linkname reflect_verifyNotInHeapPtr reflect.verifyNotInHeapPtr
func reflect_verifyNotInHeapPtr(p uintptr) bool

func reflect_verifyInHeapPtr(p uintptr) bool {
	return !reflect_verifyNotInHeapPtr(p)
}

//go:linkname reflect_escapes reflect.escapes
func reflect_escapes(x any)

//go:linkname reflect_resolveTypeOff reflect.resolveTypeOff
func reflect_resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

//go:linkname reflect_newName reflect.newName
func reflect_newName(n, tag string, exported bool) name

//go:linkname reflect_name_name reflect.name.name
func reflect_name_name(p name) (s string)

//go:linkname reflect_name_tag reflect.name.tag
func reflect_name_tag(p name) (s string)

//go:linkname reflect_name_pkgPath reflect.name.pkgPath
func reflect_name_pkgPath(p name) (s string)

//go:linkname reflect_name_isExported reflect.name.isExported
func reflect_name_isExported(p name) bool

//go:linkname reflect_name_hasTag reflect.name.hasTag
func reflect_name_hasTag(p name) bool

//go:linkname reflect_writeVarint reflect.writeVarint
//go:linkname reflect_name_readVarint reflect.name.readVarint
func reflect_name_readVarint(p name, off int) (int, int)
func reflect_writeVarint(buf []byte, n int) int
