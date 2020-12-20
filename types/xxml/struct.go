package xxml

import (
	"encoding/xml"
	"reflect"
)

var (
	structmarshalerType = reflect.TypeOf((*xml.Marshaler)(nil)).Elem()
)

func marshalStruct(writer *Writer, rval reflect.Value, opts *Options) (err error) {
	rtype := rval.Type()
	if rtype.Kind() == reflect.Ptr {
		rval = rval.Elem()
		rtype = rval.Type()
	}

	if opts.isCustomType(rtype) {
		var val string
		val, err = opts.customTypes[rtype].Marshal(rval.Interface())
		writer.WriteValue(val)
		return
	}

	if rtype.Implements(structmarshalerType) {
		var bytes []byte
		bytes, err = xml.Marshal(rval.Interface())
		writer.WriteValue(string(bytes))
		return
	}

	fieldCount := rtype.NumField()
	for i := 0; i < fieldCount; i++ {
		fv := rtype.Field(i)
		tagName := fv.Tag.Get("xml")
		if tagName == "" {
			tagName = fv.Name
		}
		writer.WriteStart(fv.Name)
		err = marshal(writer, reflect.ValueOf(fv.Type), opts)
		if err != nil {
			return
		}
		writer.WriteEnd()
	}

	return
}
