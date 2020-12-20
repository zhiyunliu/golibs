package xxml

import "reflect"

func marshalArray(writer *Writer, rval reflect.Value, opts *Options) (err error) {
	n := rval.Len()
	for i := 0; i < n; i++ {
		writer.WriteStart(opts.DefaultElement)
		marshal(writer, rval.Index(i), opts)
		writer.WriteEnd()
	}
	return
}
