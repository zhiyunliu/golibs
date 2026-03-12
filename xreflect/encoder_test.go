package xreflect

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type encItem struct {
	Val string
}

func (i encItem) String() string {
	return i.Val
}

func Test_structEncoder_encode(t *testing.T) {
	paraTtest := struct {
		A encItem
		B *encItem
	}{
		A: encItem{Val: "A"},
		B: &encItem{Val: "B"},
	}

	tests := []struct {
		name  string
		input any
		want  any
	}{
		{name: "1", input: paraTtest, want: map[string]any{"A": "A", "B": "B"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotParam, err := teststructEncoder(tt.input)
			if err != nil {
				t.Errorf("testxx() error = %v", err)
				return
			}
			if !reflect.DeepEqual(gotParam, tt.want) {
				t.Errorf("testxx() = %v, want %v", gotParam, tt.want)
			}

		})
	}
}

func teststructEncoder(input any) (params map[string]any, err error) {
	params = make(map[string]any)
	refval := reflect.ValueOf(input)
	//获取最终的类型值
	for refval.Kind() == reflect.Pointer {
		refval = refval.Elem()
	}

	if refval.Kind() != reflect.Struct {
		return params, fmt.Errorf("只能接收struct; 实际是 %s", refval.Kind().String())
	}

	fields := CachedTypeFields(refval.Type())

	for _, f := range fields.ExactName {
		if val, _, ok := f.EncoderV2(refval); ok {
			params[f.Name] = val
		}
	}
	return
}

func Test_boolEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    reflect.Value
		want any
	}{
		{name: "true", v: reflect.ValueOf(true), want: true},
		{name: "false", v: reflect.ValueOf(false), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boolEncoder(tt.v, StructOptions{})
			if got != tt.want {
				t.Errorf("boolEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boolEncoder_ptr(t *testing.T) {
	v := true
	rv := reflect.ValueOf(&v)
	got := boolEncoder(rv, StructOptions{})
	if got != true {
		t.Errorf("boolEncoder(ptr) = %v, want true", got)
	}

	// nil pointer
	var p *bool
	rv2 := reflect.ValueOf(&p).Elem()
	got2 := boolEncoder(rv2, StructOptions{})
	if got2 != nil {
		t.Errorf("boolEncoder(nil ptr) = %v, want nil", got2)
	}
}

func Test_intEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    reflect.Value
		want any
	}{
		{name: "int", v: reflect.ValueOf(42), want: 42},
		{name: "int8", v: reflect.ValueOf(int8(8)), want: int8(8)},
		{name: "int16", v: reflect.ValueOf(int16(16)), want: int16(16)},
		{name: "int32", v: reflect.ValueOf(int32(32)), want: int32(32)},
		{name: "int64", v: reflect.ValueOf(int64(64)), want: int64(64)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intEncoder(tt.v, StructOptions{})
			if got != tt.want {
				t.Errorf("intEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intEncoder_ptr(t *testing.T) {
	v := 42
	rv := reflect.ValueOf(&v)
	got := intEncoder(rv, StructOptions{})
	if got != 42 {
		t.Errorf("intEncoder(ptr) = %v, want 42", got)
	}
	var p *int
	rv2 := reflect.ValueOf(&p).Elem()
	got2 := intEncoder(rv2, StructOptions{})
	if got2 != nil {
		t.Errorf("intEncoder(nil ptr) = %v, want nil", got2)
	}
}

func Test_uintEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    reflect.Value
		want any
	}{
		{name: "uint", v: reflect.ValueOf(uint(42)), want: uint(42)},
		{name: "uint8", v: reflect.ValueOf(uint8(8)), want: uint8(8)},
		{name: "uint16", v: reflect.ValueOf(uint16(16)), want: uint16(16)},
		{name: "uint32", v: reflect.ValueOf(uint32(32)), want: uint32(32)},
		{name: "uint64", v: reflect.ValueOf(uint64(64)), want: uint64(64)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uintEncoder(tt.v, StructOptions{})
			if got != tt.want {
				t.Errorf("uintEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uintEncoder_ptr(t *testing.T) {
	v := uint(42)
	rv := reflect.ValueOf(&v)
	got := uintEncoder(rv, StructOptions{})
	if got != uint(42) {
		t.Errorf("uintEncoder(ptr) = %v, want 42", got)
	}
	var p *uint
	rv2 := reflect.ValueOf(&p).Elem()
	got2 := uintEncoder(rv2, StructOptions{})
	if got2 != nil {
		t.Errorf("uintEncoder(nil ptr) = %v, want nil", got2)
	}
}

func Test_floatEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    reflect.Value
		want any
	}{
		{name: "float32", v: reflect.ValueOf(float32(1.5)), want: float32(1.5)},
		{name: "float64", v: reflect.ValueOf(float64(2.5)), want: float64(2.5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			if tt.name == "float32" {
				got = float32Encoder(tt.v, StructOptions{})
			} else {
				got = float64Encoder(tt.v, StructOptions{})
			}
			if got != tt.want {
				t.Errorf("floatEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_floatEncoder_ptr(t *testing.T) {
	v := float64(3.14)
	rv := reflect.ValueOf(&v)
	got := float64Encoder(rv, StructOptions{})
	if got != 3.14 {
		t.Errorf("floatEncoder(ptr) = %v, want 3.14", got)
	}
	var p *float64
	rv2 := reflect.ValueOf(&p).Elem()
	got2 := float64Encoder(rv2, StructOptions{})
	if got2 != nil {
		t.Errorf("floatEncoder(nil ptr) = %v, want nil", got2)
	}
}

func Test_stringEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    reflect.Value
		want any
	}{
		{name: "non-empty", v: reflect.ValueOf("hello"), want: "hello"},
		{name: "empty", v: reflect.ValueOf(""), want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringEncoder(tt.v, StructOptions{})
			if got != tt.want {
				t.Errorf("stringEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringEncoder_ptr(t *testing.T) {
	v := "hello"
	rv := reflect.ValueOf(&v)
	got := stringEncoder(rv, StructOptions{})
	if got != "hello" {
		t.Errorf("stringEncoder(ptr) = %v, want hello", got)
	}
	var p *string
	rv2 := reflect.ValueOf(&p).Elem()
	got2 := stringEncoder(rv2, StructOptions{})
	if got2 != nil {
		t.Errorf("stringEncoder(nil ptr) = %v, want nil", got2)
	}
}

func Test_interfaceEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    reflect.Value
		want any
	}{
		{name: "int", v: reflect.ValueOf(42), want: 42},
		{name: "string", v: reflect.ValueOf("hello"), want: "hello"},
		{name: "bool", v: reflect.ValueOf(true), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := interfaceEncoder(tt.v, StructOptions{})
			if got != tt.want {
				t.Errorf("interfaceEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unsupportedTypeEncoder(t *testing.T) {
	ch := make(chan int)
	rv := reflect.ValueOf(ch)
	got := unsupportedTypeEncoder(rv, StructOptions{})
	if got != nil {
		t.Errorf("unsupportedTypeEncoder() = %v, want nil", got)
	}
}

func Test_encodeByteSlice(t *testing.T) {
	data := []byte{1, 2, 3}
	rv := reflect.ValueOf(data)
	got := encodeByteSlice(rv, StructOptions{})
	if !reflect.DeepEqual(got, data) {
		t.Errorf("encodeByteSlice() = %v, want %v", got, data)
	}
}

func Test_ptrEncoder(t *testing.T) {
	v := 42
	rv := reflect.ValueOf(&v)
	enc := newPtrEncoder(rv.Type())

	got := enc(rv, StructOptions{})
	if got != 42 {
		t.Errorf("ptrEncoder() = %v, want 42", got)
	}

	// nil pointer
	var p *int
	rv2 := reflect.ValueOf(&p).Elem()
	enc2 := newPtrEncoder(rv2.Type())
	got2 := enc2(rv2, StructOptions{})
	if got2 != nil {
		t.Errorf("ptrEncoder(nil) = %v, want nil", got2)
	}
}

func Test_sliceEncoder(t *testing.T) {
	// byte slice
	data := []byte{1, 2, 3}
	rv := reflect.ValueOf(data)
	enc := newSliceEncoder(rv.Type())
	got := enc(rv, StructOptions{SliceItem: true})
	if !reflect.DeepEqual(got, data) {
		t.Errorf("sliceEncoder(bytes) = %v, want %v", got, data)
	}

	// int slice
	intData := []int{1, 2, 3}
	rv2 := reflect.ValueOf(intData)
	enc2 := newSliceEncoder(rv2.Type())
	got2 := enc2(rv2, StructOptions{SliceItem: true})
	expected := []any{1, 2, 3}
	if !reflect.DeepEqual(got2, expected) {
		t.Errorf("sliceEncoder(int) = %v, want %v", got2, expected)
	}

	// SliceItem disabled
	got3 := enc2(rv2, StructOptions{SliceItem: false})
	if !reflect.DeepEqual(got3, intData) {
		t.Errorf("sliceEncoder(SliceItem=false) = %v, want %v", got3, intData)
	}
}

func Test_arrayEncoder(t *testing.T) {
	arr := [3]int{1, 2, 3}
	rv := reflect.ValueOf(arr)
	enc := newArrayEncoder(rv.Type())
	got := enc(rv, StructOptions{SliceItem: true})
	expected := []any{1, 2, 3}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("arrayEncoder() = %v, want %v", got, expected)
	}
}

// driverValuerEncItem implements driver.Valuer
type driverValuerEncItem struct {
	Val string
}

func (d driverValuerEncItem) Value() (driver.Value, error) {
	return "driver:" + d.Val, nil
}

// jsonMarshalerEncItem implements json.Marshaler
type jsonMarshalerEncItem struct {
	Val string
}

func (j jsonMarshalerEncItem) MarshalJSON() ([]byte, error) {
	return json.Marshal("json:" + j.Val)
}

// structValuerEncItem implements StructValuer
type structValuerEncItem struct {
	Val string
}

func (s structValuerEncItem) Value() any {
	return "valuer:" + s.Val
}

func Test_structEncoder_StructValuer(t *testing.T) {
	type testStruct struct {
		A structValuerEncItem `json:"a"`
	}
	input := testStruct{A: structValuerEncItem{Val: "test"}}
	result, err := teststructEncoder(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["a"] != "valuer:test" {
		t.Errorf("got %v, want valuer:test", result["a"])
	}
}

func Test_structEncoder_DriverValuer(t *testing.T) {
	type testStruct struct {
		A driverValuerEncItem `json:"a"`
	}
	input := testStruct{A: driverValuerEncItem{Val: "test"}}
	result, err := teststructEncoder(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["a"] != "driver:test" {
		t.Errorf("got %v, want driver:test", result["a"])
	}
}

func Test_structEncoder_TimeType(t *testing.T) {
	now := time.Now()
	type testStruct struct {
		T time.Time `json:"t"`
	}
	input := testStruct{T: now}
	result, err := teststructEncoder(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result["t"], now) {
		t.Errorf("got %v, want %v", result["t"], now)
	}
}

func Test_structEncoder_JsonMarshaler(t *testing.T) {
	type testStruct struct {
		A jsonMarshalerEncItem `json:"a"`
	}
	input := testStruct{A: jsonMarshalerEncItem{Val: "test"}}
	// Enable JSON marshaler
	params := make(map[string]any)
	refval := reflect.ValueOf(input)
	fields := CachedTypeFields(refval.Type())
	for _, f := range fields.ExactName {
		if val, _, ok := f.EncoderV2(refval); ok {
			params[f.Name] = val
		}
	}
	// json.Marshaler returns []byte
	_, isByte := params["a"].([]byte)
	if !isByte {
		t.Errorf("expected []byte, got %T", params["a"])
	}
}

// mapStringerEncItem implements fmt.Stringer for map encoder test
type mapStringerEncItem map[string]string

func (m mapStringerEncItem) String() string {
	return "map_stringer"
}

func Test_mapEncoder_Stringer(t *testing.T) {
	m := mapStringerEncItem{"a": "b"}
	rv := reflect.ValueOf(m)
	enc := newMapEncoder(rv.Type())
	got := enc(rv, StructOptions{})
	if got != "map_stringer" {
		t.Errorf("mapEncoder Stringer = %v, want map_stringer", got)
	}
}

// mapDriverValuerEncItem implements driver.Valuer
type mapDriverValuerEncItem map[string]string

func (m mapDriverValuerEncItem) Value() (driver.Value, error) {
	return "map_driver", nil
}

func Test_mapEncoder_DriverValuer(t *testing.T) {
	m := mapDriverValuerEncItem{"a": "b"}
	rv := reflect.ValueOf(m)
	enc := newMapEncoder(rv.Type())
	got := enc(rv, StructOptions{})
	if got != "map_driver" {
		t.Errorf("mapEncoder DriverValuer = %v, want map_driver", got)
	}
}

// mapJsonMarshalerEncItem implements json.Marshaler
type mapJsonMarshalerEncItem map[string]string

func (m mapJsonMarshalerEncItem) MarshalJSON() ([]byte, error) {
	return json.Marshal("map_json")
}

func Test_mapEncoder_JsonMarshaler(t *testing.T) {
	m := mapJsonMarshalerEncItem{"a": "b"}
	rv := reflect.ValueOf(m)
	enc := newMapEncoder(rv.Type())
	got := enc(rv, StructOptions{})
	_, isByte := got.([]byte)
	if !isByte {
		t.Errorf("mapEncoder JsonMarshaler = %T, want []byte", got)
	}
}

func Test_mapEncoder_unsupportedKey(t *testing.T) {
	// map with unsupported key type (float64)
	m := map[float64]string{1.0: "a"}
	rv := reflect.ValueOf(m)
	enc := newMapEncoder(rv.Type())
	got := enc(rv, StructOptions{})
	if got != nil {
		t.Errorf("mapEncoder unsupported key = %v, want nil", got)
	}
}

func Test_newTypeEncoder_allTypes(t *testing.T) {
	// Covers the newTypeEncoder switch branches
	types := []reflect.Type{
		reflect.TypeOf(true),
		reflect.TypeOf(int(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int64(0)),
		reflect.TypeOf(uint(0)),
		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint64(0)),
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
		reflect.TypeOf(""),
		reflect.TypeOf((*int)(nil)),      // pointer
		reflect.TypeOf([]int{}),          // slice
		reflect.TypeOf([3]int{}),         // array
		reflect.TypeOf(map[string]int{}), // map
		reflect.TypeOf(make(chan int)),   // unsupported
	}

	for _, typ := range types {
		enc := newTypeEncoder(typ)
		if enc == nil {
			t.Errorf("newTypeEncoder(%v) returned nil", typ)
		}
	}
}

func Test_structEncoder_Anonymous(t *testing.T) {
	type Anonymous struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}
	type Test struct {
		Anonymous
		Field3 string `json:"field3"`
	}

	paraTtest := Test{
		Anonymous: Anonymous{Field1: "field1", Field2: 2},
		Field3:    "field3",
	}

	tests := []struct {
		name  string
		input any
		want  any
	}{
		{name: "1", input: paraTtest, want: map[string]any{"field1": "field1", "field2": (2), "field3": "field3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotParam, err := teststructEncoder(tt.input)
			if err != nil {
				t.Errorf("testxx() error = %v", err)
				return
			}
			if !reflect.DeepEqual(gotParam, tt.want) {
				t.Errorf("testxx() = %v, want %v", gotParam, tt.want)
			}

		})
	}
}
