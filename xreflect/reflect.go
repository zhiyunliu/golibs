package xreflect

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"
)

var (
	fieldCache     sync.Map
	encoderCache   sync.Map
	dencoderCache  sync.Map
	stringerType   = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	scannerType    = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	mapScannerType = reflect.TypeOf((*MapScanner)(nil)).Elem()
)

type encoderFunc func(v reflect.Value) any
type dencoderFunc func(reflect.Value, any) error

func CachedTypeFields(t reflect.Type) *StructFields {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if f, ok := fieldCache.Load(t); ok {
		return f.(*StructFields)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.(*StructFields)
}

// typeFields returns a list of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typeFields(t reflect.Type) *StructFields {

	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	var count, nextCount map[reflect.Type]int

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i, cnt := 0, f.typ.NumField(); i < cnt; i++ {
				sf := f.typ.Field(i)
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}
				tag := sf.Tag.Get("db")
				if tag == "" {
					tag = sf.Tag.Get("json")
					if tag == "-" {
						continue
					}
				}

				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}

				index := make([]int, len(f.Index)+1)
				copy(index, f.Index)
				index[len(f.Index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					if name == "" {
						name = sf.Name
					}
					field := field{
						Name:      name,
						fieldName: sf.Name,
						Index:     index,
						typ:       ft,
						orgtyp:    sf.Type,
						omitEmpty: opts.Contains("omitempty"),
					}

					fields = append(fields, field)
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}
				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{Name: ft.Name(), typ: ft})
				}
			}
		}
	}

	// sort.Slice(fields, func(i, j int) bool {
	// 	x := fields
	// 	// sort field by name, breaking ties with depth, then
	// 	// breaking ties with "name came from json tag", then
	// 	// breaking ties with index sequence.
	// 	if x[i].Name != x[j].Name {
	// 		return x[i].Name < x[j].Name
	// 	}
	// 	if len(x[i].Index) != len(x[j].Index) {
	// 		return len(x[i].Index) < len(x[j].Index)
	// 	}

	// 	return byIndex(x).Less(i, j)
	// })
	exactName := make(map[string]*field, len(fields))
	for i := range fields {
		f := &fields[i]
		exactName[f.Name] = &fields[i]
		f.Encoder = typeEncoder(f.typ)
		f.Dencoder = typeDencoder(f.typ)
	}

	return &StructFields{List: fields, ExactName: exactName}
}

// func setReflectVal(field *field) {

// 	// ReflectValueOf returns field's reflect value
// 	fieldIndex := field.index
// 	switch {
// 	case len(field.StructField.Index) == 1 && fieldIndex > 0:
// 		field.ReflectValueOf = func(value reflect.Value) reflect.Value {
// 			return reflect.Indirect(value).Field(fieldIndex)
// 		}
// 	default:
// 		field.ReflectValueOf = func(v reflect.Value) reflect.Value {
// 			v = reflect.Indirect(v)
// 			for idx, fieldIdx := range field.StructField.Index {
// 				if fieldIdx >= 0 {
// 					v = v.Field(fieldIdx)
// 				} else {
// 					v = v.Field(-fieldIdx - 1)

// 					if v.IsNil() {
// 						v.Set(reflect.New(v.Type().Elem()))
// 					}

// 					if idx < len(field.StructField.Index)-1 {
// 						v = v.Elem()
// 					}
// 				}
// 			}
// 			return v
// 		}
// 	}

// }

func GetRealReflectVal(f *field, v reflect.Value) (subv reflect.Value) {
	subv = v
	for _, i := range f.Index {
		if subv.Kind() == reflect.Pointer {
			if subv.IsNil() {
				// If a struct embeds a pointer to an unexported type,
				// it is not possible to set a newly allocated value
				// since the field is unexported.
				//
				// See https://golang.org/issue/21357
				if !subv.CanSet() {
					// Invalidate subv to ensure d.value(subv) skips over
					// the JSON value without assigning it to subv.
					subv = reflect.Value{}
					break
				}
				subv.Set(reflect.New(subv.Type().Elem()))
			}
			subv = subv.Elem()
		}
		subv = subv.Field(i)
	}
	return
}
