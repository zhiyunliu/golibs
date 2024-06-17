package datetime

import (
	"reflect"
	"testing"
	"time"
)

func TestNewDateTime(t *testing.T) {
	testTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	d := New(testTime)

	if !reflect.DeepEqual(d.Time, testTime) {
		t.Errorf("New 函数未正确设置时间")
	}

	// 测试默认格式
	if d.Format() != DefaultTimeformat {
		t.Errorf("New 函数未正确设置默认格式")
	}
}

func TestMarshalJSON(t *testing.T) {
	testTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	d := New(testTime)

	expectedResult := []byte(`"2022-01-01 12:00:00"`)
	result, _ := d.MarshalJSON()

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("MarshalJSON 函数返回结果错误")
	}
}

// 编写其他函数的测试用例...
func TestDateTime_Scan(t *testing.T) {
	// 测试传入 time.Time 类型的参数
	inputTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	dt := &DateTime{}
	err := dt.Scan(inputTime)
	if err != nil {
		t.Errorf("Scan 函数未能正确处理 time.Time 类型参数: %v", err)
	}
	expectedTime := New(inputTime.Local())
	if !reflect.DeepEqual(dt, expectedTime) {
		t.Errorf("Scan 函数未能正确转换为 DateTime 结构体")
	}

	// 测试传入字符串参数
	inputString := "2022-01-01 12:00:00"
	err = dt.Scan(inputString)
	if err != nil {
		t.Errorf("Scan 函数未能正确处理字符串类型参数: %v", err)
	}
	expectedTimeStr, _ := time.Parse(DefaultTimeformat, inputString)
	expectedTimeStr = expectedTimeStr.Local()
	expectedTime = New(expectedTimeStr)
	if !reflect.DeepEqual(dt, expectedTime) {
		t.Errorf("Scan 函数未能正确转换为 DateTime 结构体")
	}

	// 测试传入无法处理的参数
	err = dt.Scan(123) //假设传入一个无法处理的参数
	if err == nil {
		t.Errorf("Scan 函数未返回错误信息来处理无法处理的参数")
	}

	err = dt.Scan("123") //假设传入一个无法处理的参数
	if err == nil {
		t.Errorf("Scan 函数未返回错误信息来处理无法处理的参数")
	}
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	// 测试正常情况下的 JSON 反序列化
	bytes := []byte(`"2022-01-01 12:00:00"`)
	dt := &DateTime{}
	err := dt.UnmarshalJSON(bytes)
	if err != nil {
		t.Errorf("UnmarshalJSON 函数未能正确处理有效的 JSON 数据: %v", err)
	}
	expectedTime, _ := time.Parse(DefaultTimeformat, "2022-01-01 12:00:00")
	expectedTime = expectedTime.Local()
	expectedDateTime := New(expectedTime)
	if !reflect.DeepEqual(dt, expectedDateTime) {
		t.Errorf("UnmarshalJSON 函数未能正确反序列化 JSON 数据")
	}

	//测试无效的 JSON 数据情况
	invalidBytes := []byte(`"invalid-time-format"`)
	err = dt.UnmarshalJSON(invalidBytes)
	if err == nil {
		t.Errorf("UnmarshalJSON 函数未返回错误信息来处理无效的 JSON 数据")
	}
}

func TestDateTime_String(t *testing.T) {
	testTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	d := New(testTime)

	expectedStr := "2022-01-01 12:00:00"
	result := d.String()

	if result != expectedStr {
		t.Errorf("String 方法返回结果错误")
	}
}

func TestDateTime_Value(t *testing.T) {
	testTime := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	d := New(testTime)

	expectedValue := "2022-01-01 12:00:00"
	result, err := d.Value()

	if err != nil {
		t.Errorf("Value 方法返回错误: %v", err)
	}

	if result != expectedValue {
		t.Errorf("Value 方法返回结果错误")
	}
}
