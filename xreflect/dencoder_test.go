package xreflect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test_boolDecoder(t *testing.T) {

	val := struct {
		B1 bool
		B2 *bool
		B3 bool
		B4 *bool
	}{
		B1: false,
	}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()
	var refTrue *bool = new(bool)
	*refTrue = true

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("B1"), val: true, wantErr: false},
		{name: "2.", v: refVal.FieldByName("B2"), val: true, wantErr: false},
		{name: "3.", v: refVal.FieldByName("B3"), val: refTrue, wantErr: false},
		{name: "4.", v: refVal.FieldByName("B4"), val: refTrue, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := boolDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("boolDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_intDecoder(t *testing.T) {
	val := struct {
		B1 int
		B2 *int
		B3 int
		B4 *int
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()
	var refTrue *int = new(int)
	*refTrue = 1

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("B1"), val: 1, wantErr: false},
		{name: "2.", v: refVal.FieldByName("B2"), val: 1, wantErr: false},
		{name: "3.", v: refVal.FieldByName("B3"), val: refTrue, wantErr: false},
		{name: "4.", v: refVal.FieldByName("B4"), val: refTrue, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := intDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("intDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stringDecoder(t *testing.T) {
	val := struct {
		B1 string
		B2 *string
		B3 string
		B4 *string
		B5 string
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()
	var refTrue *string = new(string)
	*refTrue = "a"
	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("B1"), val: "a", wantErr: false},
		{name: "2.", v: refVal.FieldByName("B2"), val: "a", wantErr: false},
		{name: "3.", v: refVal.FieldByName("B3"), val: refTrue, wantErr: false},
		{name: "4.", v: refVal.FieldByName("B4"), val: refTrue, wantErr: false},
		{name: "5.", v: refVal.FieldByName("B5"), val: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := stringDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("stringDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_floatDecoder_dencode(t *testing.T) {
	val := struct {
		A1 float64
		A2 *float64
		A3 float64
		A4 *float64

		B1 float64
		B2 *float64
		B3 float64
		B4 *float64
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()
	var arefTrue *float64 = new(float64)
	*arefTrue = 1.1

	var brefTrue *float64 = new(float64)
	*brefTrue = 2.1

	tests := []struct {
		name    string
		bits    floatDecoder
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "a1.", bits: 32, v: refVal.FieldByName("A1"), val: 1.1, wantErr: false},
		{name: "a2.", bits: 32, v: refVal.FieldByName("A2"), val: 1.1, wantErr: false},
		{name: "a3.", bits: 32, v: refVal.FieldByName("A3"), val: arefTrue, wantErr: false},
		{name: "a4.", bits: 32, v: refVal.FieldByName("A4"), val: arefTrue, wantErr: false},
		{name: "b1.", bits: 64, v: refVal.FieldByName("B1"), val: 2.1, wantErr: false},
		{name: "b2.", bits: 64, v: refVal.FieldByName("B2"), val: 2.1, wantErr: false},
		{name: "b3.", bits: 64, v: refVal.FieldByName("B3"), val: brefTrue, wantErr: false},
		{name: "b4.", bits: 64, v: refVal.FieldByName("B4"), val: brefTrue, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.bits.dencode(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("floatDecoder.dencode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_interfaceDecoder(t *testing.T) {
	val := struct {
		V any
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("V"), val: 1, wantErr: false},
		{name: "2.", v: refVal.FieldByName("V"), val: "a", wantErr: false},
		{name: "3.", v: refVal.FieldByName("V"), val: true, wantErr: false},
		{name: "4.", v: refVal.FieldByName("V"), val: 1.1, wantErr: false},
		{name: "5.", v: refVal.FieldByName("V"), val: []int{1, 2, 3}, wantErr: false},
		{name: "6.", v: refVal.FieldByName("V"), val: map[string]int{"a": 1, "b": 2}, wantErr: false},
		{name: "7.", v: refVal.FieldByName("V"), val: struct{ A int }{A: 1}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := interfaceDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("interfaceDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type item struct {
	A int
}

func (i *item) MapScan(val any) error {
	if i == nil {
		i = new(item)
	}
	mp := val.(map[string]any)
	i.A = mp["A"].(int)
	return nil
}

func (i *item) Scan(val any) error {
	i.A, _ = val.(int)
	return nil
}

func Test_mapScanDecoder(t *testing.T) {

	val := struct {
		V item
		W *item
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("V"), val: map[string]any{"A": 1}},
		{name: "2.", v: refVal.FieldByName("W"), val: map[string]any{"A": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mapScanDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("mapScanDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MapItem map[string]string

func (i *MapItem) Scan(val any) error {
	if i == nil {
		*i = make(MapItem)
	}

	json.Unmarshal([]byte(fmt.Sprint(val)), &i)

	return nil
}

func Test_mapScanDecoder2(t *testing.T) {

	val := struct {
		V MapItem
		W *MapItem
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("V"), val: `{"V":"1"}`},
		{name: "2.", v: refVal.FieldByName("W"), val: `{"W":"2"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mapDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("mapScanDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MapArray []MapItem

func (i *MapArray) Scan(val any) error {
	if i == nil {
		*i = make(MapArray, 0)
	}

	json.Unmarshal([]byte(fmt.Sprint(val)), i)
	return nil
}

func Test_SliceDecoder3(t *testing.T) {

	val := struct {
		V MapArray
		W *MapArray
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("V"), val: `[{"V":"1"}]`},
		{name: "2.", v: refVal.FieldByName("W"), val: `[{"W":"2"}]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := newSliceDecoder(tt.v.Type())
			if err := decoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("mapScanDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

type XMapItem map[string]any

func (i *XMapItem) MapScan(val any) error {
	if i == nil {
		*i = make(XMapItem)
	}

	json.Unmarshal([]byte(fmt.Sprint(val)), i)
	return nil
}

type XMapArray []XMapItem

func (i *XMapArray) MapScan(val any) error {
	if i == nil {
		*i = make(XMapArray, 0)
	}

	json.Unmarshal([]byte(fmt.Sprint(val)), i)
	return nil
}

func Test_SliceDecoder4(t *testing.T) {

	val := struct {
		V XMapArray
		W *XMapArray
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("V"), val: `[{"V":"1"}]`},
		{name: "2.", v: refVal.FieldByName("W"), val: `[{"W":"2"}]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := newSliceDecoder(tt.v.Type())
			if err := decoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("mapScanDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

type IntArray []int

func (i *IntArray) Scan(val any) error {
	if i == nil {
		*i = make(IntArray, 0)
	}

	json.Unmarshal([]byte(fmt.Sprint(val)), i)
	return nil
}

func Test_SliceDecoder5(t *testing.T) {

	val := struct {
		V IntArray
		W *IntArray
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("V"), val: `[1,2,3]`},
		{name: "2.", v: refVal.FieldByName("W"), val: `[4,5,6]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := newSliceDecoder(tt.v.Type())
			if err := decoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("mapScanDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func Test_SliceDecoder6(t *testing.T) {

	val := struct {
		V []int
		W *[]int
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("V"), val: []int{1, 2, 3}},
		{name: "2.", v: refVal.FieldByName("W"), val: []int{4, 5, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := newSliceDecoder(tt.v.Type())
			if err := decoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("mapScanDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func Test_structDecoder(t *testing.T) {
	type structTestItem struct {
		U item
		V item
		W *item
		X *item
		Y *item
		Z *item
	}

	val := structTestItem{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()
	var Y *item = new(item)
	Y.A = 4

	expected := structTestItem{
		U: item{A: 0},
		V: item{A: 1},
		W: &item{A: 2},
		X: &item{A: 3},
		Y: &item{A: 4},
		Z: &item{A: 5},
	}

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "0.", v: refVal.FieldByName("V"), val: item{A: 0}, wantErr: false},
		{name: "1.", v: refVal.FieldByName("V"), val: 1, wantErr: false},
		{name: "2.", v: refVal.FieldByName("W"), val: 2, wantErr: false},
		{name: "3.", v: refVal.FieldByName("X"), val: &item{A: 3}, wantErr: false},
		{name: "4.", v: refVal.FieldByName("Y"), val: &Y, wantErr: false},
		{name: "4.", v: refVal.FieldByName("Z"), val: item{A: 5}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := structDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("structDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if !reflect.DeepEqual(val, expected) {
		t.Errorf("structDecoder() DeepEqual got= %v, expected %v", val, expected)
	}
}

func Test_uintDecoder(t *testing.T) {
	val := struct {
		B1 uint
		B2 *uint
		B3 uint
		B4 *uint
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()
	var refTrue *uint = new(uint)
	*refTrue = 1

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("B1"), val: uint(1), wantErr: false},
		{name: "2.", v: refVal.FieldByName("B2"), val: uint(2), wantErr: false},
		{name: "3.", v: refVal.FieldByName("B3"), val: refTrue, wantErr: false},
		{name: "4.", v: refVal.FieldByName("B4"), val: refTrue, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := uintDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("intDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mapDecoder(t *testing.T) {
	type mapTestItem struct {
		B1 map[string]int
		B2 map[string]int
		B3 *map[string]int
		B4 *map[string]int
		B5 *map[string]int
		B6 *map[string]int
		B7 *map[string]int
	}
	val := mapTestItem{}

	expected := mapTestItem{
		B1: map[string]int{"a": 1},
		B2: map[string]int{"a": 2},
		B3: &map[string]int{"a": 3},
		B4: &map[string]int{"a": 4},
		B5: &map[string]int{"a": 5},
		B6: nil,
		B7: &map[string]int{},
	}

	var Y *map[string]int = new(map[string]int)
	*Y = map[string]int{"a": 5}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "1.", v: refVal.FieldByName("B1"), val: map[string]int{"a": 1}},
		{name: "2.", v: refVal.FieldByName("B2"), val: &map[string]int{"a": 2}},
		{name: "3.", v: refVal.FieldByName("B3"), val: map[string]int{"a": 3}, wantErr: false},
		{name: "4.", v: refVal.FieldByName("B4"), val: &map[string]int{"a": 4}, wantErr: false},
		{name: "5.", v: refVal.FieldByName("B5"), val: &Y, wantErr: false},
		{name: "6.", v: refVal.FieldByName("B6"), val: nil, wantErr: false},
		{name: "7.", v: refVal.FieldByName("B7"), val: map[string]int{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mapDecoder(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("mapDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	if !reflect.DeepEqual(val, expected) {
		t.Errorf("mapDecoder()	DeepEqual got= %v, expected %v", val, expected)
	}

}

func Test_dencodeByteSlice(t *testing.T) {
	val := struct {
		B1 []byte
		B2 []byte
	}{}

	refVal := reflect.ValueOf(&val)
	refVal = refVal.Elem()

	tests := []struct {
		name    string
		v       reflect.Value
		val     any
		wantErr bool
	}{
		{name: "normal", v: refVal.FieldByName("B1"), val: []byte{1, 2, 3}, wantErr: false},
		{name: "nil", v: refVal.FieldByName("B2"), val: (*[]byte)(nil), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dencodeByteSlice(tt.v, tt.val); (err != nil) != tt.wantErr {
				t.Errorf("dencodeByteSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if !reflect.DeepEqual(val.B1, []byte{1, 2, 3}) {
		t.Errorf("dencodeByteSlice() B1 = %v, want [1 2 3]", val.B1)
	}
}

func Test_dencodeByteSlice_zeroValue(t *testing.T) {
	val := struct {
		B1 []byte
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	// zero value (empty []byte) should return nil
	err := dencodeByteSlice(refVal.FieldByName("B1"), []byte{})
	if err != nil {
		t.Errorf("dencodeByteSlice(zero) error = %v", err)
	}
}

func Test_defaultDecoder(t *testing.T) {
	val := struct {
		M map[string]int
	}{}
	refVal := reflect.ValueOf(&val).Elem()

	err := defaultDecoder(refVal.FieldByName("M"), map[string]int{"a": 1})
	if err != nil {
		t.Errorf("defaultDecoder(map) error = %v", err)
	}
	if !reflect.DeepEqual(val.M, map[string]int{"a": 1}) {
		t.Errorf("defaultDecoder(map) = %v, want map[a:1]", val.M)
	}
}

func Test_defaultDecoder_nonMap(t *testing.T) {
	val := struct {
		S string
	}{S: "hello"}
	refVal := reflect.ValueOf(&val).Elem()
	// default decoder with non-map kind does nothing
	err := defaultDecoder(refVal.FieldByName("S"), "world")
	if err != nil {
		t.Errorf("defaultDecoder(string) error = %v", err)
	}
}

func Test_boolDecoder_zeroValue(t *testing.T) {
	val := struct {
		B bool
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	// zero value should not change the field
	err := boolDecoder(refVal.FieldByName("B"), false)
	if err != nil {
		t.Errorf("boolDecoder(zero) error = %v", err)
	}
	if val.B != false {
		t.Errorf("boolDecoder(zero) = %v, want false", val.B)
	}
}

func Test_intDecoder_zeroValue(t *testing.T) {
	val := struct {
		I int
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := intDecoder(refVal.FieldByName("I"), 0)
	if err != nil {
		t.Errorf("intDecoder(zero) error = %v", err)
	}
}

func Test_uintDecoder_zeroValue(t *testing.T) {
	val := struct {
		U uint
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := uintDecoder(refVal.FieldByName("U"), uint(0))
	if err != nil {
		t.Errorf("uintDecoder(zero) error = %v", err)
	}
}

func Test_floatDecoder_zeroValue(t *testing.T) {
	val := struct {
		F float64
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := float64Decoder(refVal.FieldByName("F"), 0.0)
	if err != nil {
		t.Errorf("floatDecoder(zero) error = %v", err)
	}
}

func Test_stringDecoder_zeroValue(t *testing.T) {
	val := struct {
		S string
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := stringDecoder(refVal.FieldByName("S"), "")
	if err != nil {
		t.Errorf("stringDecoder(zero) error = %v", err)
	}
}

func Test_mapScanDecoder_nil(t *testing.T) {
	val := struct {
		V item
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := mapScanDecoder(refVal.FieldByName("V"), nil)
	if err != nil {
		t.Errorf("mapScanDecoder(nil) error = %v", err)
	}
}

func Test_mapDecoder_nil(t *testing.T) {
	val := struct {
		M map[string]int
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := mapDecoder(refVal.FieldByName("M"), nil)
	if err != nil {
		t.Errorf("mapDecoder(nil) error = %v", err)
	}
}

func Test_structDecoder_nil(t *testing.T) {
	val := struct {
		I item
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := structDecoder(refVal.FieldByName("I"), nil)
	if err != nil {
		t.Errorf("structDecoder(nil) error = %v", err)
	}
}

func Test_structDecoder_zeroRefVal(t *testing.T) {
	val := struct {
		I item
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	// zero reflect.Value of the same type
	err := structDecoder(refVal.FieldByName("I"), item{})
	if err != nil {
		t.Errorf("structDecoder(zero) error = %v", err)
	}
}

func Test_newTypeDencoder_allTypes(t *testing.T) {
	// Covers the newTypeDencoder switch branches
	types := []reflect.Type{
		reflect.TypeOf(true),
		reflect.TypeOf(int(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(uint(0)),
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
		reflect.TypeOf(""),
		reflect.TypeOf(struct{}{}),
		reflect.TypeOf(map[string]int{}),
		reflect.TypeOf([]int{}),
		reflect.TypeOf([3]int{}),
		reflect.TypeOf(make(chan int)), // default branch
	}

	for _, typ := range types {
		dec := newTypeDencoder(typ)
		if dec == nil {
			t.Errorf("newTypeDencoder(%v) returned nil", typ)
		}
	}
}

func Test_intDecoder_error(t *testing.T) {
	val := struct {
		I int
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	// string that can't be parsed as int
	err := intDecoder(refVal.FieldByName("I"), "abc")
	if err == nil {
		t.Error("intDecoder(\"abc\") expected error, got nil")
	}
}

func Test_uintDecoder_error(t *testing.T) {
	val := struct {
		U uint
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := uintDecoder(refVal.FieldByName("U"), "abc")
	if err == nil {
		t.Error("uintDecoder(\"abc\") expected error, got nil")
	}
}

func Test_floatDecoder_error(t *testing.T) {
	val := struct {
		F float64
	}{}
	refVal := reflect.ValueOf(&val).Elem()
	err := float64Decoder(refVal.FieldByName("F"), "abc")
	if err == nil {
		t.Error("floatDecoder(\"abc\") expected error, got nil")
	}
}
