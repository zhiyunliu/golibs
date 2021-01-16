package xxml

import (
	"encoding/xml"
	"reflect"
	"strconv"
	"time"
)

var (
	custommarshalerType = reflect.TypeOf((*Customize)(nil)).Elem()
	timemarshalerType   = reflect.TypeOf((*time.Time)(nil)).Elem()
)

//Marshal Marshal
func Marshal(obj interface{}, opts ...Option) (val string, err error) {
	var newOpts Options = defaults
	newOpts.customMarshalers = map[string]XxmlMarshaler{}
	newOpts.customUnmarshalers = map[string]XxmlUnmarshaler{}

	for i := range opts {
		opts[i](&newOpts)
	}
	writer := NewWriter(newOpts.Compress)
	rootEle := newOpts.RootElement
	refVal := reflect.ValueOf(obj)

	writer.WriteStart(rootEle)

	err = marshal(writer, refVal, &newOpts)
	if err != nil {
		return
	}
	writer.WriteEnd()
	val = writer.String()
	return
}

func marshal(writer *Writer, refVal reflect.Value, opts *Options) (err error) {
	if refVal.IsZero() {
		return nil
	}
	if refVal.Kind() == reflect.Interface {
		refVal = reflect.ValueOf(refVal.Interface())
	}

	if isSimpleType(refVal) {
		return marshalSimple(writer, refVal)
	}

	if refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}
	if opts.IsCustomMarshaler(refVal.Type()) {
		obj := refVal.Interface()
		cusMarshaler := obj.(XxmlMarshaler)
		val, err1 := cusMarshaler.Marshal()
		if err1 != nil {
			err = err1
			return
		}
		writer.WriteValue(val)
		return
	}

	switch refVal.Kind() {
	case reflect.Map:
		return marshalMap(writer, refVal, opts)
	case reflect.Array, reflect.Slice:
		return marshalArray(writer, refVal, opts)
	case reflect.Struct:

		return marshalStruct(writer, refVal, opts)
	}

	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func marshalSimple(writer *Writer, val reflect.Value) (err error) {
	if val.Kind() == reflect.Interface {
		val = reflect.ValueOf(val.Interface())
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		writer.WriteValue(strconv.FormatInt(val.Int(), 10))
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		writer.WriteValue(strconv.FormatUint(val.Uint(), 10))
		return
	case reflect.Float32, reflect.Float64:
		writer.WriteValue(strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()))
		return
	case reflect.String:
		writer.WriteValue(val.String())
		return
	case reflect.Bool:
		writer.WriteValue(strconv.FormatBool(val.Bool()))
		return
	case reflect.Array:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			break
		}
		var bytes []byte
		if val.CanAddr() {
			bytes = val.Slice(0, val.Len()).Bytes()
		} else {
			bytes = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(bytes), val)
		}
		writer.WriteValue(string(bytes))
		return
	case reflect.Slice:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			break
		}
		writer.WriteValue(string(val.Bytes()))
		return
	}
	return &xml.UnsupportedTypeError{val.Type()}
}

func isSimpleType(v reflect.Value) bool {

	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Array:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return false
		}
		return true
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return false
		}
		return true

	}
	return false
}
