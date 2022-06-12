package unsafelib

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var (
	typeMap   = map[uintptr]reflect.Type{}
	typeMapMu sync.Mutex
)

type InspChangeFunc func(insp *Inspection) (doNext bool)

type Inspection struct {
	typ       *rtype
	inHeap    bool
	writable  bool
	writeMode bool
}

// Inspect type
func Inspect(typ reflect.Type) *Inspection {
	if typ == nil {
		panic("type not supplied")
	}
	if typ.Kind() == reflect.Ptr {
		panic(fmt.Sprintf("type kind not acceptable: %s of %s", typ.Kind(), typ))
	}

	insp := &Inspection{
		typ: reflectTypeToRtype(typ),
	}
	insp.inHeap = !reflect_verifyNotInHeapPtr(insp.getPtr())
	insp.writable = insp.inHeap
	return insp
}

// Type returns the reflect.Type.
func (i *Inspection) Type() reflect.Type {
	return rtypeToReflectType(i.typ)
}

// Change perform action.
func (i *Inspection) Change(fs ...InspChangeFunc) {
	defer func() {
		if !i.writable {
			i.writeMode = false
		}
	}()
	i.writeMode = true
	if !i.writable {
		// make page temporary writable during action execution.
		prf := new_mem_profile(i.getPtr(), 1, PROT_FLAG_RO, PROT_FLAG_RW)
		defer prf.toggle()
		prf.toggle()
	}
	for _, f := range fs {
		if f != nil && !f(i) {
			break
		}
	}
}

func (i *Inspection) Struct() *StructType {
	return ReinterpretPtr[StructType](i.typ)
}

func (i *Inspection) Fields() StructFields {
	fields := &i.Struct().fields
	if fields == nil {
		return nil
	}
	fieldsSh := ReinterpretPtr[reflect.SliceHeader](fields)

	// if not in heap, relocate slice data.
	// this will keep the original fields slice header cap, and len
	// so any size grow will not affect the original,
	// but can change any existing field attributes by writing new data pointer.
	if reflect_verifyNotInHeapPtr(fieldsSh.Data) && fieldsSh.Cap >= fieldsSh.Len {
		hFieldSh := make([]StructField, fieldsSh.Len, fieldsSh.Cap)
		for i := 0; i < len(*fields); i++ {
			hFieldSh[i] = (*fields)[i]
		}
		fields = &hFieldSh
		if i.writeMode {
			// write slice header back
			*fieldsSh = *ReinterpretPtr[reflect.SliceHeader](fields)
		}
	}
	return *fields
}

func (i *Inspection) IsWritable() bool {
	return i.writable
}

func (i *Inspection) IsInHeap() bool {
	return i.inHeap
}

// ---

// tbd.
func (i *Inspection) getRef() reflect.Type {
	typeMapMu.Lock()
	defer typeMapMu.Unlock()
	v, ok := typeMap[i.getPtr()]
	if !ok {
		newRef := i.newRef()
		typeMap[i.getPtr()] = newRef
		return newRef
	}
	return v
}

// tbd.
func (i *Inspection) newRef() reflect.Type {
	typ := i.Type()
	switch typ.Kind() {
	case reflect.Struct:
		// return reflect.StructOf()
	}
	return nil
}

// getPtr pointer of type ref (*reflect.rtype)
func (i *Inspection) getPtr() uintptr {
	return uintptr(unsafe.Pointer(i.typ))
}
