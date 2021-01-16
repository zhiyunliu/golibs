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

	fieldCount := rtype.NumField()
	for i := 0; i < fieldCount; i++ {
		fv := rtype.Field(i)
		tagName := fv.Tag.Get("xml")
		if tagName == "" {
			tagName = fv.Name
		}
		if rval.Field(i).IsZero() {
			continue
		}
		writer.WriteStart(fv.Name)
		err = marshal(writer, rval.Field(i), opts)
		if err != nil {
			return
		}
		writer.WriteEnd()
	}

	return
}
