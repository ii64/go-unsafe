package unsafelib

import "reflect"

type InspTransformStructFieldFunc func(field *StructField, name *TypeName)

func TransformStructField(f InspTransformStructFieldFunc) InspChangeFunc {
	return func(insp *Inspection) (doNext bool) {
		if insp.Type().Kind() == reflect.Struct {
			fields := insp.Fields()
			for i, _ := range fields {
				rfield := &fields[i]
				mod := (&TypeName{}).fromName(rfield.name)
				if f != nil {
					f(rfield, mod) // exec
					rfield.name = mod.toName()
				}
			}
		}
		return true
	}
}
