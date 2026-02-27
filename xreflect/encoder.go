package xreflect

import (
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(t, encoderFunc(func(v reflect.Value, opts StructOptions) any {
		wg.Wait()
		return f(v, opts)
	}))
	if loaded {
		return fi.(encoderFunc)
	}

	// Compute the real encoder and replace the indirect func with it.
	f = newTypeEncoder(t)
	wg.Done()
	encoderCache.Store(t, f)
	return f
}

// newTypeEncoder constructs an encoderFunc for a type.
// The returned encoder only checks CanAddr when allowAddr is true.
func newTypeEncoder(t reflect.Type) encoderFunc {
	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Pointer:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}

func unsupportedTypeEncoder(v reflect.Value, _ StructOptions) any {
	return nil
}

func boolEncoder(v reflect.Value, _ StructOptions) any {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	return v.Bool()
}

func intEncoder(v reflect.Value, _ StructOptions) any {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	return v.Interface()
}

func uintEncoder(v reflect.Value, _ StructOptions) any {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	return v.Interface()
}

type floatEncoder int // number of bits

func (bits floatEncoder) encode(v reflect.Value, _ StructOptions) any {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	return v.Interface()
}

var (
	float32Encoder = (floatEncoder(32)).encode
	float64Encoder = (floatEncoder(64)).encode
)

func stringEncoder(v reflect.Value, _ StructOptions) any {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	return v.String()
}

func interfaceEncoder(v reflect.Value, _ StructOptions) any {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	return v.Interface()
}

type structEncoder struct {
	fields *StructFields
}

func (se structEncoder) encode(v reflect.Value, opts StructOptions) any {

	tmpv := v
	for tmpv.Kind() == reflect.Pointer {
		if tmpv.IsNil() {
			return nil
		}
		tmpv = tmpv.Elem()
	}

	rft := tmpv.Type()

	if rft.Implements(structValuerType) {
		return tmpv.Interface().(StructValuer).Value()
	}

	if tmpv.CanAddr() && tmpv.Addr().Type().Implements(structValuerType) {
		return tmpv.Addr().Interface().(StructValuer).Value()
	}

	// 检查是否是 time.Time 类型或可以转换为 time.Time 的类型别名
	if rft == timeType || rft.ConvertibleTo(timeType) {
		return tmpv.Interface()
	}

	if rft.Implements(stringerType) {
		return tmpv.Interface().(fmt.Stringer).String()
	}

	if tmpv.CanAddr() && tmpv.Addr().Type().Implements(stringerType) {
		return tmpv.Addr().Interface().(fmt.Stringer).String()
	}

	if rft.Implements(driverValuerType) {
		tmpVal, _ := tmpv.Interface().(driver.Valuer).Value()
		return tmpVal
	}

	if tmpv.CanAddr() && tmpv.Addr().Type().Implements(driverValuerType) {
		tmpVal, _ := tmpv.Addr().Interface().(driver.Valuer).Value()
		return tmpVal
	}

	if !opts.DisableJSONMarshaler {
		if rft.Implements(jsonMarshalerType) {
			tmpVal, _ := tmpv.Interface().(json.Marshaler).MarshalJSON()
			return tmpVal
		}

		if tmpv.CanAddr() && tmpv.Addr().Type().Implements(jsonMarshalerType) {
			tmpVal, _ := tmpv.Addr().Interface().(json.Marshaler).MarshalJSON()
			return tmpVal
		}
	}

	if rft.Implements(textMarshalerType) {
		tmpVal, _ := tmpv.Interface().(encoding.TextMarshaler).MarshalText()
		return tmpVal
	}

	if tmpv.CanAddr() && tmpv.Addr().Type().Implements(textMarshalerType) {
		tmpVal, _ := v.Addr().Interface().(encoding.TextMarshaler).MarshalText()
		return tmpVal
	}

	if opts.StructField && opts.IsValidDepth(1) {
		params := make(map[string]any)
		for _, f := range se.fields.ExactName {
			if val, ok := f.Encoder(v, opts); ok {
				params[f.Name] = val
			}
		}
		return params
	}
	return unsupportedTypeEncoder(v, opts)
}

func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{fields: CachedTypeFields(t)}
	return se.encode
}

type mapEncoder struct {
	elemEnc encoderFunc
}

func (me mapEncoder) encode(v reflect.Value, opts StructOptions) any {
	rft := v.Type()

	if rft.Implements(stringerType) {
		return v.Interface().(fmt.Stringer).String()
	}

	if rft.Implements(driverValuerType) {
		tmpVal, _ := v.Interface().(driver.Valuer).Value()
		return tmpVal
	}

	if rft.Implements(jsonMarshalerType) {
		tmpVal, _ := v.Interface().(json.Marshaler).MarshalJSON()
		return tmpVal
	}

	return unsupportedTypeEncoder(v, opts)
}

func newMapEncoder(t reflect.Type) encoderFunc {
	switch t.Key().Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		if !t.Key().Implements(stringerType) {
			return unsupportedTypeEncoder
		}
	}
	me := mapEncoder{elemEnc: typeEncoder(t.Elem())}
	return me.encode
}

func encodeByteSlice(v reflect.Value, opts StructOptions) any {
	return v.Bytes()
}

// sliceEncoder just wraps an arrayEncoder, checking to make sure the value isn't nil.
type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se sliceEncoder) encode(v reflect.Value, opts StructOptions) any {
	return se.arrayEnc(v, opts)
}

func newSliceEncoder(t reflect.Type) encoderFunc {
	// Byte slices get special treatment; arrays don't.
	if t.Elem().Kind() == reflect.Uint8 {
		p := reflect.PointerTo(t.Elem())
		if !p.Implements(stringerType) {
			return encodeByteSlice
		}
	}
	enc := sliceEncoder{arrayEnc: newArrayEncoder(t)}
	return enc.encode
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae arrayEncoder) encode(v reflect.Value, opts StructOptions) any {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if opts.SliceItem {
		arrLen := v.Len()
		arrVal := make([]any, arrLen)
		for i := 0; i < arrLen; i++ {
			arrVal[i] = ae.elemEnc(v.Index(i), opts)
		}
		return arrVal
	}
	return v.Interface()
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	enc := arrayEncoder{elemEnc: typeEncoder(t.Elem())}
	return enc.encode
}

type ptrEncoder struct {
	elemEnc encoderFunc
}

func (pe ptrEncoder) encode(v reflect.Value, opts StructOptions) any {
	if v.IsNil() {
		return nil
	}
	return pe.elemEnc(v.Elem(), opts)
}

func newPtrEncoder(t reflect.Type) encoderFunc {
	enc := ptrEncoder{elemEnc: typeEncoder(t.Elem())}
	return enc.encode
}
