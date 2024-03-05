package xtypes

import (
	"reflect"
	"testing"
)

func TestXMaps_Scan_1(t *testing.T) {
	type Args struct {
		A string
		B int
	}
	tests := []struct {
		name    string
		ms      XMaps
		args    []Args
		wantErr bool
	}{
		{name: "1.", ms: XMaps{XMap{"A": "a1", "B": 1}}, args: []Args{}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.Scan(&tt.args); (err != nil) != tt.wantErr {
				t.Errorf("XMaps.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.args) != 1 {
				t.Error("反射失败")
			}
		})
	}
}

func TestXMaps_Scan_2(t *testing.T) {
	type Args struct {
		A string
		B int
	}
	tests := []struct {
		name    string
		ms      XMaps
		args    []*Args
		wantErr bool
	}{
		{name: "1.", ms: XMaps{XMap{"A": "a2", "B": 2}}, args: []*Args{}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.Scan(&tt.args); (err != nil) != tt.wantErr {
				t.Errorf("XMaps.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.args) != 1 {
				t.Error("反射失败")
			}
		})
	}
}

func TestXMaps_Get(t *testing.T) {
	// 初始化一个 XMaps
	testXMaps := XMaps{
		XMap{"name": "Alice", "age": 25},
		XMap{"name": "Bob", "age": 30},
	}

	// 测试获取指定索引的 XMap 的方法
	result1 := testXMaps.Get(0)
	result2 := testXMaps.Get(1)
	result3 := testXMaps.Get(2)

	// 验证预期的结果
	expected1 := XMap{"name": "Alice", "age": 25}
	expected2 := XMap{"name": "Bob", "age": 30}
	expected3 := XMap{}

	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Get(0) returned %v, expected %v", result1, expected1)
	}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Get(1) returned %v, expected %v", result2, expected2)
	}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Get(2) returned %v, expected %v", result3, expected3)
	}
}
