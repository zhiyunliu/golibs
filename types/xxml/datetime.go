package xxml

import (
	"reflect"
	"time"
)

func marshalTime(writer *Writer, rval reflect.Value, opts *Options) (err error) {
	if rval.Kind() == reflect.Ptr {
		rval = rval.Elem()
	}
	timeVal := rval.Interface().(time.Time)
	writer.WriteValue(timeVal.Format(opts.TimeFormat))
	return nil
}
