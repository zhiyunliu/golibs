package xtypes

import (
	"reflect"
	"sort"
	"testing"
)

func TestXMap_Keys(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"name": "Alice", "age": 25, "isStudent": true}

	// 测试获取键的方法
	keys := testMap.Keys()
	// 验证预期的键
	expectedKeys := []string{"name", "age", "isStudent"}

	sort.Strings(keys)
	sort.Strings(expectedKeys)

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Keys() returned %v, expected %v", keys, expectedKeys)
	}
}

func TestXMap_Merge(t *testing.T) {
	// 初始化两个 XMap
	map1 := XMap{"name": "Alice", "age": 25}
	map2 := XMap{"isStudent": true, "city": "New York"}

	// 合并两个 XMap
	map1.Merge(map2)

	// 验证合并后的结果
	expectedMap := XMap{"name": "Alice", "age": 25, "isStudent": true, "city": "New York"}
	if !reflect.DeepEqual(map1, expectedMap) {
		t.Errorf("Merge() result is %v, expected %v", map1, expectedMap)
	}
}

func TestXMap_Get(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"name": "Alice", "age": 25, "isStudent": true}

	// 测试获取值的方法
	value, ok := testMap.Get("name")

	// 验证预期的值和存在性
	expectedValue := "Alice"
	if !ok || value != expectedValue {
		t.Errorf("Get(\"name\") returned value %v and ok value %v, expected %v and true", value, ok, expectedValue)
	}

	// 测试一个不存在的键
	_, ok = testMap.Get("city")
	if ok {
		t.Errorf("Get(\"city\") returned ok value true, expected false")
	}
}

func TestXMap_SMap(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"name": "Alice", "age": 25, "isStudent": true}

	// 测试获取转换为 SMap 的方法
	sMap := testMap.SMap()

	// 验证预期的 SMap 结果
	expectedSMap := SMap{"name": "Alice", "age": "25", "isStudent": "true"}
	if !reflect.DeepEqual(sMap, expectedSMap) {
		t.Errorf("SMap() returned %v, expected %v", sMap, expectedSMap)
	}
}

func TestXMap_GetBool(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"name": "Alice", "age": 25, "isStudent": true}

	// 测试获取布尔值的方法
	result1 := testMap.GetBool("isStudent")
	result2 := testMap.GetBool("name")
	result3 := testMap.GetBool("city") // 不存在的键

	// 验证预期的布尔值
	expected1 := true
	expected2 := false
	expected3 := false

	if result1 != expected1 {
		t.Errorf("GetBool(\"isStudent\") returned %v, expected %v", result1, expected1)
	}
	if result2 != expected2 {
		t.Errorf("GetBool(\"name\") returned %v, expected %v", result2, expected2)
	}
	if result3 != expected3 {
		t.Errorf("GetBool(\"city\") returned %v, expected %v", result3, expected3)
	}
}

func TestXMap_GetString(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"name": "Alice", "age": 25, "isStudent": true}

	// 测试获取字符串值的方法
	result1 := testMap.GetString("name")
	result2 := testMap.GetString("age")
	result3 := testMap.GetString("city") // 不存在的键

	// 验证预期的字符串值
	expected1 := "Alice"
	expected2 := "25"
	expected3 := ""

	if result1 != expected1 {
		t.Errorf("GetString(\"name\") returned %v, expected %v", result1, expected1)
	}
	if result2 != expected2 {
		t.Errorf("GetString(\"age\") returned %v, expected %v", result2, expected2)
	}
	if result3 != expected3 {
		t.Errorf("GetString(\"city\") returned %v, expected %v", result3, expected3)
	}
}

func TestXMap_GetInt(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"age": 25, "score": "90", "isStudent": true}

	// 测试获取整数值的方法
	result1, _ := testMap.GetInt("age")
	result3, _ := testMap.GetInt("isStudent")
	result2, _ := testMap.GetInt("score")
	result4, err4 := testMap.GetInt("name") // 不存在的键

	// 验证预期的整数值和错误
	expected1 := 25
	expected2 := 90
	expected3 := 0
	expected4 := 0
	var expectedErr4 error

	if result1 != expected1 {
		t.Errorf("GetInt(\"age\") returned %v, expected %v", result1, expected1)
	}
	if result2 != expected2 {
		t.Errorf("GetInt(\"score\") returned %v, expected %v", result2, expected2)
	}
	if result3 != expected3 {
		t.Errorf("GetInt(\"isStudent\") returned %v, expected %v", result3, expected3)
	}
	if result4 != expected4 || err4 != expectedErr4 {
		t.Errorf("GetInt(\"name\") returned %v and error %v, expected %v and %v", result4, err4, expected4, expectedErr4)
	}
}

func TestXMap_GetInt64(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"age": 25, "score": "90", "isStudent": true}

	// 测试获取整数值的方法
	result1, _ := testMap.GetInt64("age")
	result3, _ := testMap.GetInt64("isStudent")
	result2, _ := testMap.GetInt64("score")
	result4, err4 := testMap.GetInt64("name") // 不存在的键

	// 验证预期的整数值和错误
	var expected1 int64 = 25
	var expected2 int64 = 90
	var expected3 int64 = 0
	var expected4 int64 = 0
	var expectedErr4 error

	if result1 != expected1 {
		t.Errorf("GetInt(\"age\") returned %v, expected %v", result1, expected1)
	}
	if result2 != expected2 {
		t.Errorf("GetInt(\"score\") returned %v, expected %v", result2, expected2)
	}
	if result3 != expected3 {
		t.Errorf("GetInt(\"isStudent\") returned %v, expected %v", result3, expected3)
	}
	if result4 != expected4 || err4 != expectedErr4 {
		t.Errorf("GetInt(\"name\") returned %v and error %v, expected %v and %v", result4, err4, expected4, expectedErr4)
	}
}

func TestXMap_GetFloat(t *testing.T) {
	// 初始化一个 XMap
	testMap := XMap{"age": 25.5, "score": "90.5", "isStudent": true}

	// 测试获取浮点值的方法
	result1, _ := testMap.GetFloat64("age")
	result2, _ := testMap.GetFloat64("score")
	result3, _ := testMap.GetFloat64("isStudent")
	result4, err4 := testMap.GetFloat64("name") // 不存在的键

	// 验证预期的浮点值和错误
	expected1 := 25.5
	expected2 := 90.5
	expected3 := 0.0
	expected4 := 0.0
	var expectedErr4 error

	if result1 != expected1 {
		t.Errorf("GetFloat64(\"age\") returned %v, expected %v", result1, expected1)
	}
	if result2 != expected2 {
		t.Errorf("GetFloat(\"score\") returned %v, expected %v", result2, expected2)
	}
	if result3 != expected3 {
		t.Errorf("GetFloat(\"isStudent\") returned %v, expected %v", result3, expected3)
	}
	if result4 != expected4 || err4 != expectedErr4 {
		t.Errorf("GetFloat(\"name\") returned %v and error %v, expected %v and %v", result4, err4, expected4, expectedErr4)
	}
}
