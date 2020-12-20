package xxml

import (
	"encoding/xml"
	"reflect"
	"strconv"
)

//Marshal Marshal
func Marshal(obj interface{}, opts ...Option) (val string, err error) {
	newOpts := defaults
	writer := NewWriter()
	for i := range opts {
		opts[i](newOpts)
	}
	rootEle := newOpts.RootElement
	refVal := reflect.ValueOf(obj)

	writer.WriteStart(rootEle)
	if refVal.IsNil() {
		writer.WriteEnd()
		return writer.String(), nil
	}

	err = marshal(writer, refVal, newOpts)
	if err != nil {
		return
	}
	writer.WriteEnd()
	val = writer.String()
	return
}

func marshal(writer *Writer, refVal reflect.Value, opts *Options) (err error) {
	if isSimpleType(refVal) {
		marshalSimple(writer, refVal)
	}
	switch refVal.Kind() {
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

func marshalSimple(writer *Writer, val reflect.Value) (string, error) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Array:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			break
		}
		// [...]byte
		var bytes []byte
		if val.CanAddr() {
			bytes = val.Slice(0, val.Len()).Bytes()
		} else {
			bytes = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(bytes), val)
		}
		return string(bytes), nil
	case reflect.Slice:
		if val.Type().Elem().Kind() != reflect.Uint8 {
			break
		}
		// []byte
		return string(val.Bytes()), nil
	}
	return "", &xml.UnsupportedTypeError{val.Type()}
}

func isSimpleType(v reflect.Value) bool {
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
	default:
		return false
	}
	return false
}
