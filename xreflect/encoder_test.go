package xreflect

import (
	"fmt"
	"reflect"
	"testing"
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
		if val, _, ok := f.EncoderV2(refval, StructOptions{MaxDepth: 1}); ok {
			params[f.Name] = val
		}
	}
	return
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
