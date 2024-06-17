package xtypes

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestSMap_Keys(t *testing.T) {
	// 初始化一个 SMap 实例
	sm := make(SMap)
	sm["key1"] = "value1"
	sm["key2"] = "value2"

	// 调用 Keys 方法
	keys := sm.Keys()

	// 检查返回的键是否包含预期的键
	expectedKeys := []string{"key1", "key2"}

	sort.Strings(keys)
	sort.Strings(expectedKeys)

	for _, expectedKey := range expectedKeys {
		keyExists := false
		for _, key := range keys {
			if key == expectedKey {
				keyExists = true
				break
			}
		}
		if !keyExists {
			t.Errorf("Keys 方法返回的键不包含预期的键: %s", expectedKey)
		}
	}

	// 检查返回的键的数量是否正确
	if len(keys) != len(expectedKeys) {
		t.Error("Keys 方法返回的键数量不正确")
	}
}

func TestSMap_Value(t *testing.T) {
	// 初始化一个 SMap 实例
	sm := make(SMap)
	sm["key1"] = "value1"
	sm["key2"] = "value2"

	// 调用 Value 方法
	value, err := sm.Value()

	// 检查返回的值是否符合预期
	expectedValue := `{"key1":"value1","key2":"value2"}`
	if value != expectedValue {
		t.Errorf("Value 方法返回的值不符合预期: %s", value)
	}

	// 检查返回的错误是否为 nil
	if err != nil {
		t.Errorf("Value 方法返回了错误: %v", err)
	}
}

func TestSMap_MarshalBinary(t *testing.T) {
	// 初始化一个 SMap 实例
	sm := make(SMap)
	sm["key1"] = "value1"
	sm["key2"] = "value2"

	// 调用 MarshalBinary 方法
	data, err := sm.MarshalBinary()

	// 检查返回的数据是否符合预期
	expectedData, _ := json.Marshal(map[string]string{"key1": "value1", "key2": "value2"})
	if string(data) != string(expectedData) {
		t.Errorf("MarshalBinary 方法返回的数据不符合预期: %s", string(data))
	}

	// 检查返回的错误是否为 nil
	if err != nil {
		t.Errorf("MarshalBinary 方法返回了错误: %v", err)
	}
}

func TestSMap_GetWithDefault(t *testing.T) {
	// 初始化一个 SMap 实例
	sm := make(SMap)
	sm["key1"] = "value1"
	sm["key2"] = "value2"

	// 调用 GetWithDefault 方法并传入存在的键
	result1 := sm.GetWithDefault("key1", "default1")
	if result1 != "value1" {
		t.Errorf("GetWithDefault 方法未返回预期的值: %s", result1)
	}

	// 调用 GetWithDefault 方法并传入不存在的键
	result2 := sm.GetWithDefault("key3", "default3")
	if result2 != "default3" {
		t.Errorf("GetWithDefault 方法未返回默认值: %s", result2)
	}
}

func TestSMap_Get(t *testing.T) {
	// 初始化一个 SMap 实例
	sm := make(SMap)
	sm["key1"] = "value1"
	sm["key2"] = "value2"

	// 调用 Get 方法并传入存在的键
	result1 := sm.Get("key1")
	if result1 != "value1" {
		t.Errorf("Get 方法未返回预期的值: %s", result1)
	}

	// 调用 Get 方法并传入不存在的键
	result2 := sm.Get("key3")
	if result2 != "" {
		t.Errorf("Get 方法未返回空字符串: %s", result2)
	}
}

// 创建一个结构体用于接收扫描的数据
type TestStruct struct {
	Name  string
	Value string
}

func TestSMap_Scan(t *testing.T) {
	// 初始化一个 SMap 实例
	sm := make(SMap)
	sm["Name"] = "John"
	sm["Value"] = "Doe"

	// 创建目标结构体实例
	target := TestStruct{}

	// 调用 Scan 方法
	err := sm.Scan(&target)

	// 检查返回的错误是否为 nil
	if err != nil {
		t.Errorf("Scan 方法返回了错误: %v", err)
	}

	// 检查扫描的数据是否赋值正确
	if target.Name != "John" || target.Value != "Doe" {
		t.Errorf("扫描的数据未能正确赋值到结构体中")
	}
}

func TestSMap_Read(t *testing.T) {
	// 初始化一个 SMap 实例
	sm := make(SMap)
	sm["key1"] = "value1"
	sm["key2"] = "value2"

	// 准备一个缓冲区用于读取
	buffer := make([]byte, 0, 100)

	// 调用 Read 方法
	n, err := sm.Read(buffer)

	// 检查返回的读取字节数是否符合预期
	if n != len(buffer) {
		t.Errorf("Read 方法返回的读取字节数不符合预期: %v", n)
	}

	// 如果有返回错误，进行错误处理
	if err != nil {
		t.Errorf("Read 方法返回了错误: %v", err)
	}
}
