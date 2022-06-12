package unsafelib

import "reflect"

func (sf *StructField) Meta() *TypeName {
	return (&TypeName{}).fromName(sf.name)
}
func (sf *StructField) SetMeta(m *TypeName) {
	sf.name = m.toName()
}

func (sf StructField) ToReflectStructField() reflect.StructField {
	mfield := (&TypeName{}).fromName(sf.name)
	return reflect.StructField{
		Name:    mfield.Name,
		PkgPath: mfield.PkgPath,
		Type:    rtypeToReflectType(sf.typ),
		Tag:     mfield.Tag,
		Offset:  sf.offsetEmbed,
		// Index:     []int{},
		// Anonymous: false,
	}
}

type StructFields []StructField

func (sf StructFields) CreateStructType() reflect.Type {
	fields := make([]reflect.StructField, len(sf))
	for i := 0; i < len(sf); i++ {
		fields[i] = sf[i].ToReflectStructField()
	}
	return reflect.StructOf(fields)
}
