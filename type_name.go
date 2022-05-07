package unsafelib

import "reflect"

type TypeName struct {
	PkgPath  string
	Name     string
	Tag      reflect.StructTag
	Exported bool
}

func (t *TypeName) fromName(n name) *TypeName {
	t.PkgPath = reflect_name_pkgPath(n)
	t.Name = reflect_name_name(n)
	t.Tag = reflect.StructTag(reflect_name_tag(n))
	t.Exported = reflect_name_isExported(n)
	return t
}

func (t TypeName) HasTag() bool {
	return len(t.Tag) > 0
}
func (t TypeName) HasPkgPath() bool {
	return len(t.PkgPath) > 0
}

func (t TypeName) Equal(v TypeName) bool {
	return t.Name == v.Name && t.Tag == v.Tag && t.Exported == v.Exported
}

func (t TypeName) toName() name {
	return name{bytes: &t.Bytes()[0]}
}

func (t TypeName) Bytes() []byte {
	var bits byte
	var nameLen [10]byte
	var tagLen [10]byte
	var pkgPathLen [10]byte
	nameLenWritten := reflect_writeVarint(nameLen[:], len(t.Name))
	tagLenWritten := reflect_writeVarint(tagLen[:], len(t.Tag))
	pkgPathLenWritten := reflect_writeVarint(pkgPathLen[:], len(t.PkgPath))

	l := 1 + nameLenWritten + len(t.Name)
	if t.Exported {
		bits |= 1 << 0
	}
	if t.HasTag() {
		l += tagLenWritten + len(t.Tag)
		bits |= 1 << 1
	}
	if t.HasPkgPath() {
		l += pkgPathLenWritten + len(t.PkgPath)
		bits |= 1 << 2
	}

	ret := make([]byte, l)
	ret[0] = bits

	off := 1
	// write name
	copy(ret[off:], nameLen[:nameLenWritten])
	copy(ret[off+nameLenWritten:], t.Name)
	off = off + nameLenWritten + len(t.Name)

	// write tag
	if t.HasTag() {
		tb := ret[off:]
		copy(tb, tagLen[:tagLenWritten])
		copy(tb[tagLenWritten:], t.Tag)
		off = off + tagLenWritten + len(t.Tag)
	}

	// write pkgPath
	if t.HasPkgPath() {
		tb := ret[off:]
		copy(tb, pkgPathLen[:pkgPathLenWritten])
		copy(tb[pkgPathLenWritten:], t.PkgPath)
		off = off + pkgPathLenWritten + len(t.PkgPath)
	}
	_ = off
	return ret
}
