package xxml

import (
	"reflect"
	"strconv"
)

func marshalMap(writer *Writer, rval reflect.Value, opts *Options) (err error) {
	keys := rval.MapKeys()
	var kstr string
	for i := range keys {
		rkv := keys[i]
		switch rkv.Type().Kind() {
		case reflect.String:
			kstr = rkv.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			kstr = strconv.FormatInt(rkv.Int(), 10)
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			kstr = strconv.FormatUint(rkv.Uint(), 10)
		default:
			err = errorUnsupportType(rkv.Type().Name())
			return
		}
		rval := rval.MapIndex(rkv)
		if rval.IsZero() {
			continue
		}
		writer.WriteStart(kstr)

		err = marshal(writer, rval, opts)
		if err != nil {
			return
		}
		writer.WriteEnd()
	}
	return
}
