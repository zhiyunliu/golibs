package xreflect

import (
	"reflect"
)

type StructFields struct {
	List      []field
	ExactName map[string]*field
	embedMap  map[string][]*field
}

// 反序列化结构体的对应属性
func (f *StructFields) Dencode(rv reflect.Value, name string, val any) (err error) {
	field, ok := f.ExactName[name]
	if !ok {
		return nil
	}
	err = field.Dencoder(rv, val)
	if err != nil {
		return
	}
	list, ok := f.embedMap[name]
	if !ok {
		return nil
	}
	for i := range list {
		list[i].Dencoder(rv, val)
	}
	return nil
}

type field struct {
	Name         string
	fieldName    string
	Index        []int
	typ          reflect.Type
	orgtyp       reflect.Type
	omitEmpty    bool
	encoderFunc  encoderFunc
	dencoderFunc dencoderFunc
}

func (f *field) Dencoder(rv reflect.Value, val any) (err error) {
	fv := GetRealReflectVal(f, rv) //fv := rv.Field(field.Index)
	if !fv.IsValid() {
		return nil
	}
	err = f.dencoderFunc(fv, val)
	return
}

func (f *field) Encoder(rv reflect.Value) (val any, ok bool) {
	fv := GetRealReflectVal(f, rv) //rv.Field(f.Index)
	if !fv.IsValid() {
		return nil, false
	}
	return f.encoderFunc(fv), true
}

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].Index {
		if k >= len(x[j].Index) {
			return false
		}
		if xik != x[j].Index[k] {
			return xik < x[j].Index[k]
		}
	}
	return len(x[i].Index) < len(x[j].Index)
}
