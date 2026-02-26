package xreflect

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// DateTime 是 time.Time 的类型别名，不继承任何方法
type DateTime time.Time

// DateTimeWithMapValuer 是 time.Time 的类型别名，实现了 MapValuer 接口（值接收者）
type DateTimeWithMapValuer time.Time

func (d DateTimeWithMapValuer) MapValue() interface{} {
	return time.Time(d).Format("2006-01-02 15:04:05")
}

// DateTimeWithPtrMapValuer 是 time.Time 的类型别名，实现了 MapValuer 接口（指针接收者）
type DateTimeWithPtrMapValuer time.Time

func (d *DateTimeWithPtrMapValuer) MapValue() interface{} {
	return time.Time(*d).Format("2006/01/02")
}

// DateTimeWithStringer 是 time.Time 的类型别名，实现了 fmt.Stringer 接口（指针接收者）
type DateTimeWithStringer time.Time

func (d *DateTimeWithStringer) String() string {
	return time.Time(*d).Format("2006-01-02")
}

// CustomJsonMarshaler 实现了 json.Marshaler 接口（指针接收者）
type CustomJsonMarshaler struct {
	Value string
}

func (c *CustomJsonMarshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal("custom:" + c.Value)
}

// CustomTextMarshaler 实现了 encoding.TextMarshaler 接口（值接收者）
type CustomTextMarshaler struct {
	Value string
}

func (c CustomTextMarshaler) MarshalText() ([]byte, error) {
	return []byte("text:" + c.Value), nil
}

// NonTimeMapValuer 是一个非 time.Time 基础的 struct，实现了 MapValuer 接口
type NonTimeMapValuer struct {
	X int
	Y int
}

func (n NonTimeMapValuer) MapValue() interface{} {
	return map[string]int{"x": n.X, "y": n.Y}
}

// NonTimeStringer 是一个非 time.Time 基础的 struct，实现了 fmt.Stringer 接口
type NonTimeStringer struct {
	First string
	Last  string
}

func (n NonTimeStringer) String() string {
	return n.First + " " + n.Last
}

type TestStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TestStruct2 struct {
	ID   int    `json:"id"`
	Name string `json:"-"`
}

type TestStruct3 struct {
	ID   int
	Name string
}

type TestStruct4 struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

func TestAnyToMap(t *testing.T) {
	obj := TestStruct{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"id":   1,
		"name": "John",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should convert struct to map correctly")

	result, err = AnyToMap(expected)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should convert struct to map correctly")
}

func TestAnyToMapWithPointer(t *testing.T) {
	obj := &TestStruct{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"id":   1,
		"name": "John",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should convert struct pointer to map correctly")
}

func TestAnyToMapWithIgnoredField(t *testing.T) {
	obj := TestStruct2{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"id": 1,
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should ignore fields with '-' tag")
}

func TestAnyToMapWithEmptyTag(t *testing.T) {
	obj := TestStruct3{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"ID":   1,
		"Name": "John",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should use field name when json tag is empty")
}

func TestAnyToMapWithNestedAnonymousStruct(t *testing.T) {
	type Nested struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}
	obj := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Nested
	}{
		ID:   1,
		Name: "John",
		Nested: Nested{
			Field1: "NestedField",
			Field2: 2,
		},
	}

	expected := map[string]interface{}{
		"id":     1,
		"name":   "John",
		"field1": "NestedField",
		"field2": 2,
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "Test failed, expected: '%v', got:  '%v'", expected, result)
}

func TestAnyToMapWithNestedAnonymousStruct2(t *testing.T) {
	obj := struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Nested struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		} `json:"nested"`
	}{
		ID:   1,
		Name: "John",
		Nested: struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		}{
			Field1: "NestedField",
			Field2: 2,
		},
	}

	expected := map[string]interface{}{
		"id":   1,
		"name": "John",
		"nested": map[string]interface{}{
			"field1": "NestedField",
			"field2": 2,
		},
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "Test failed, expected: '%v', got:  '%v'", expected, result)
}

type Match struct {
	MatchId   string   `json:"match_id"`
	MapName   string   `json:"map_name"`
	BeginTime string   `json:"begin_time"`
	EndTime   string   `json:"won_time"`
	WonDate   int      `json:"won_date"`
	TeamCount int      `json:"team_count"`
	GameType  string   `json:"game_type"`
	TeamSize  int      `json:"team_size"`
	Winners   []Player `json:"winners"`
}

type Player struct {
	Id    string `json:"playerId"`
	Name  string `json:"name"`
	Kills int    `json:"kills"`
}
type UserItem struct {
	Match
	UserId     int    `json:"user_id"`
	LoginId    string `json:"login_id"`
	PubgName   string `json:"pubg_nickname"`
	PubgUserId string `json:"pubg_id"`
}

func TestAnyToMapWithNestedAndAnonymousStruct(t *testing.T) {
	obj := UserItem{
		Match: Match{
			MatchId:   "1",
			MapName:   "Map1",
			BeginTime: "10:00",
			EndTime:   "11:00",
			WonDate:   20220202,
			TeamCount: 2,
			GameType:  "Type1",
			TeamSize:  5,
			Winners: []Player{
				{
					Id:    "Player1",
					Name:  "John",
					Kills: 10,
				},
				{
					Id:    "Player2",
					Name:  "Doe",
					Kills: 15,
				},
			},
		},
		UserId:     1,
		LoginId:    "User1",
		PubgName:   "User1",
		PubgUserId: "PubgUser1",
	}

	expected := map[string]interface{}{
		"match_id":   "1",
		"map_name":   "Map1",
		"begin_time": "10:00",
		"won_time":   "11:00",
		"won_date":   20220202,
		"team_count": 2,
		"game_type":  "Type1",
		"team_size":  5,
		"winners": []interface{}{
			map[string]interface{}{
				"playerId": "Player1",
				"name":     "John",
				"kills":    10,
			},
			map[string]interface{}{
				"playerId": "Player2",
				"name":     "Doe",
				"kills":    15,
			},
		},
		"user_id":       1,
		"login_id":      "User1",
		"pubg_nickname": "User1",
		"pubg_id":       "PubgUser1",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "Test failed, expected: '%v', got:  '%v'", expected, result)
}

func TestAnyToMapWithOmitEmpty(t *testing.T) {
	obj := TestStruct4{
		ID:       1,
		Name:     "John",
		Email:    "", // Empty value
		Password: "secret",
	}

	expected := map[string]interface{}{
		"id":       1,
		"name":     "John",
		"email":    "",
		"password": "secret",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMapWithOmitEmpty not error")
	assert.Equal(t, expected, result, "TestAnyToMapWithOmitEmpty should correctly parse json tags with omitempty")
}

// TestAnyToMapWithDateTimeAlias 测试 type DateTime time.Time 这样的类型别名
// 由于 ConvertibleTo 检测，应该被视为值类型而不是递归展开
func TestAnyToMapWithDateTimeAlias(t *testing.T) {
	now := time.Now()
	type Order struct {
		ID        int      `json:"id"`
		CreatedAt DateTime `json:"created_at"`
	}

	obj := Order{
		ID:        1,
		CreatedAt: DateTime(now),
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, 1, result["id"])
	// DateTime 应该作为整体值保留，而不是被拆解为 time.Time 的内部字段
	_, ok := result["created_at"]
	assert.True(t, ok, "created_at should exist in result")
	// 验证值是 DateTime 类型而不是 map
	_, isMap := result["created_at"].(map[string]interface{})
	assert.False(t, isMap, "created_at should NOT be a map (should be treated as value type)")
}

// TestAnyToMapWithMapValuer 测试实现了 MapValuer 接口的类型
func TestAnyToMapWithMapValuer(t *testing.T) {
	now := time.Date(2026, 2, 26, 10, 30, 0, 0, time.UTC)
	type Order struct {
		ID        int                   `json:"id"`
		CreatedAt DateTimeWithMapValuer `json:"created_at"`
	}

	obj := Order{
		ID:        1,
		CreatedAt: DateTimeWithMapValuer(now),
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, 1, result["id"])
	// MapValuer 返回的是格式化后的字符串
	assert.Equal(t, "2026-02-26 10:30:00", result["created_at"])
}

// TestAnyToMapWithStringerPointerReceiver 测试指针接收者实现 fmt.Stringer 的类型
func TestAnyToMapWithStringerPointerReceiver(t *testing.T) {
	now := time.Date(2026, 2, 26, 0, 0, 0, 0, time.UTC)
	type Order struct {
		ID        int                  `json:"id"`
		CreatedAt DateTimeWithStringer `json:"created_at"`
	}

	obj := Order{
		ID:        1,
		CreatedAt: DateTimeWithStringer(now),
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, 1, result["id"])
	// 因为实现了 Stringer（指针接收者），应该被视为值类型
	_, isMap := result["created_at"].(map[string]interface{})
	assert.False(t, isMap, "created_at should NOT be a map when type implements Stringer")
}

// TestAnyToMapWithTimeTime 测试原生 time.Time 仍然正常工作
func TestAnyToMapWithTimeTime(t *testing.T) {
	now := time.Now()
	type Event struct {
		Name      string    `json:"name"`
		StartTime time.Time `json:"start_time"`
	}

	obj := Event{
		Name:      "meeting",
		StartTime: now,
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, "meeting", result["name"])
	assert.Equal(t, now, result["start_time"])
}

// TestAnyToMapWithNilPointer 测试nil指针输入
func TestAnyToMapWithNilPointer(t *testing.T) {
	var obj *TestStruct
	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{}, result)
}

// TestAnyToMapWithUnsupportedType 测试不支持的类型输入
func TestAnyToMapWithUnsupportedType(t *testing.T) {
	_, err := AnyToMap(123)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported input type")

	_, err = AnyToMap("hello")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported input type")
}

// TestAnyToMapWithMaxDepth 测试 WithMaxDepth 选项
func TestAnyToMapWithMaxDepth(t *testing.T) {
	type Inner struct {
		Value string `json:"value"`
	}
	type Middle struct {
		Inner Inner `json:"inner"`
	}
	type Outer struct {
		Middle Middle `json:"middle"`
	}

	obj := Outer{
		Middle: Middle{
			Inner: Inner{
				Value: "deep",
			},
		},
	}

	// 默认深度，应该可以递归到最深层
	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	middle, ok := result["middle"].(map[string]interface{})
	assert.True(t, ok)
	inner, ok := middle["inner"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "deep", inner["value"])
}

// TestAnyToMapWithMaxDepthExceeded 测试超过最大深度的情况
func TestAnyToMapWithMaxDepthExceeded(t *testing.T) {
	type Inner struct {
		Value string `json:"value"`
	}
	type Middle struct {
		Inner Inner `json:"inner"`
	}
	type Outer struct {
		Middle Middle `json:"middle"`
	}

	obj := Outer{
		Middle: Middle{
			Inner: Inner{
				Value: "deep",
			},
		},
	}

	// maxDepth=1 时，只递归一层
	result, err := AnyToMap(obj, WithMaxDepth(1))
	assert.Nil(t, err)
	// middle 在 depth=0 递归到 depth=1, depth=1 >= maxDepth=1，返回 nil
	assert.Nil(t, result["middle"])
}

// TestAnyToMapWithTextMarshaler 测试实现了 encoding.TextMarshaler 接口的类型
func TestAnyToMapWithTextMarshaler(t *testing.T) {
	type Config struct {
		Name   string              `json:"name"`
		Custom CustomTextMarshaler `json:"custom"`
	}

	obj := Config{
		Name:   "test",
		Custom: CustomTextMarshaler{Value: "world"},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, "test", result["name"])
	// CustomTextMarshaler 实现了 TextMarshaler，应该被视为值类型
	_, isMap := result["custom"].(map[string]interface{})
	assert.False(t, isMap, "custom should NOT be a map when type implements TextMarshaler")
}

// TestAnyToMapWithPtrMapValuer 测试实现了 MapValuer 接口（指针接收者）的类型
// 传入指针使字段可寻址，覆盖 getMapValuer 的 CanAddr 路径
func TestAnyToMapWithPtrMapValuer(t *testing.T) {
	now := time.Date(2026, 2, 26, 10, 30, 0, 0, time.UTC)
	type Order struct {
		ID        int                      `json:"id"`
		CreatedAt DateTimeWithPtrMapValuer `json:"created_at"`
	}

	obj := &Order{
		ID:        1,
		CreatedAt: DateTimeWithPtrMapValuer(now),
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, 1, result["id"])
	// 指针接收者的 MapValuer，通过传入指针使字段可寻址，走 CanAddr 路径
	assert.Equal(t, "2026/02/26", result["created_at"])
}

// TestAnyToMapSliceWithSpecialTypes 测试切片中包含特殊类型
func TestAnyToMapSliceWithSpecialTypes(t *testing.T) {
	now := time.Now()
	type Schedule struct {
		Name  string      `json:"name"`
		Times []time.Time `json:"times"`
	}

	obj := Schedule{
		Name:  "daily",
		Times: []time.Time{now, now.Add(time.Hour)},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, "daily", result["name"])
	times, ok := result["times"].([]any)
	assert.True(t, ok)
	assert.Len(t, times, 2)
	// time.Time 在切片中应该被视为特殊类型，保留原始值
	assert.Equal(t, now, times[0])
	assert.Equal(t, now.Add(time.Hour), times[1])
}

// TestAnyToMapSliceWithDateTimeAlias 测试切片中包含 DateTime 类型别名
func TestAnyToMapSliceWithDateTimeAlias(t *testing.T) {
	now := time.Now()
	type Schedule struct {
		Name  string     `json:"name"`
		Times []DateTime `json:"times"`
	}

	obj := Schedule{
		Name:  "daily",
		Times: []DateTime{DateTime(now), DateTime(now.Add(time.Hour))},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	times, ok := result["times"].([]any)
	assert.True(t, ok)
	assert.Len(t, times, 2)
	// DateTime 在切片中也应该被视为特殊值类型
	_, isMap := times[0].(map[string]interface{})
	assert.False(t, isMap, "DateTime in slice should NOT be a map")
}

// TestAnyToMapSliceWithMapValuer 测试切片中包含实现了 MapValuer 的类型
func TestAnyToMapSliceWithMapValuer(t *testing.T) {
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	type Schedule struct {
		Name  string                  `json:"name"`
		Times []DateTimeWithMapValuer `json:"times"`
	}

	obj := Schedule{
		Name:  "monthly",
		Times: []DateTimeWithMapValuer{DateTimeWithMapValuer(t1), DateTimeWithMapValuer(t2)},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	times, ok := result["times"].([]any)
	assert.True(t, ok)
	assert.Len(t, times, 2)
	// MapValuer 类型在切片中应该被视为特殊值类型（保留原始值）
	_, isMap := times[0].(map[string]interface{})
	assert.False(t, isMap, "MapValuer in slice should NOT be a map")
}

// TestAnyToMapEmptySlice 测试包含空切片的结构体
func TestAnyToMapEmptySlice(t *testing.T) {
	type Container struct {
		Name  string   `json:"name"`
		Items []string `json:"items"`
	}

	obj := Container{
		Name:  "empty",
		Items: []string{},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, "empty", result["name"])
	items, ok := result["items"].([]any)
	assert.True(t, ok)
	assert.Len(t, items, 0)
}

// TestGetUnderlyingType 直接测试 getUnderlyingType 内部函数
func TestGetUnderlyingType(t *testing.T) {
	timeType := reflect.TypeOf(time.Time{})

	// time.Time → time.Time
	assert.Equal(t, timeType, getUnderlyingType(reflect.TypeOf(time.Time{})))

	// DateTime (type DateTime time.Time) → time.Time
	assert.Equal(t, timeType, getUnderlyingType(reflect.TypeOf(DateTime{})))

	// *time.Time → time.Time (解包指针)
	assert.Equal(t, timeType, getUnderlyingType(reflect.TypeOf((*time.Time)(nil))))

	// int → int (非struct类型直接返回)
	intType := reflect.TypeOf(0)
	assert.Equal(t, intType, getUnderlyingType(intType))

	// 普通 struct → 返回自身
	structType := reflect.TypeOf(TestStruct{})
	assert.Equal(t, structType, getUnderlyingType(structType))
}

// TestIsUnderlyingTimeType 直接测试 isUnderlyingTimeType 内部函数
func TestIsUnderlyingTimeType(t *testing.T) {
	// time.Time → true (struct 分支，getUnderlyingType 返回 time.Time)
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf(time.Time{})))

	// DateTime → true
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf(DateTime{})))

	// *DateTime → true (指针解包)
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf((*DateTime)(nil))))

	// []DateTime → true (slice元素检查)
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf([]DateTime{})))

	// int → false
	assert.False(t, isUnderlyingTimeType(reflect.TypeOf(0)))

	// string → false
	assert.False(t, isUnderlyingTimeType(reflect.TypeOf("")))

	// map[string]DateTime → true (map值检查)
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf(map[string]DateTime{})))
}

// TestIsSpecialValueType 直接测试 isSpecialValueType 内部函数
func TestIsSpecialValueType(t *testing.T) {
	// time.Time → true
	assert.True(t, isSpecialValueType(reflect.ValueOf(time.Time{})))

	// DateTime → true (ConvertibleTo time.Time)
	assert.True(t, isSpecialValueType(reflect.ValueOf(DateTime{})))

	// 普通struct → false
	assert.False(t, isSpecialValueType(reflect.ValueOf(TestStruct{})))

	// CustomTextMarshaler → true (实现了 TextMarshaler)
	assert.True(t, isSpecialValueType(reflect.ValueOf(CustomTextMarshaler{Value: "x"})))

}

// TestAnyToMapWithMapInput 测试 map 类型输入
func TestAnyToMapWithMapInput(t *testing.T) {
	input := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	result, err := AnyToMap(input)
	assert.Nil(t, err)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, 42, result["key2"])
}

// TestAnyToMapWithMultipleAnonymousStructs 测试多层匿名结构体嵌套
func TestAnyToMapWithMultipleAnonymousStructs(t *testing.T) {
	type Base struct {
		ID int `json:"id"`
	}
	type Middle struct {
		Base
		Name string `json:"name"`
	}
	type Top struct {
		Middle
		Extra string `json:"extra"`
	}

	obj := Top{
		Middle: Middle{
			Base: Base{ID: 1},
			Name: "test",
		},
		Extra: "info",
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, 1, result["id"])
	assert.Equal(t, "test", result["name"])
	assert.Equal(t, "info", result["extra"])
}

// TestAnyToMapWithNonTimeMapValuer 测试非 time.Time 基础的 struct 实现 MapValuer
// 覆盖 isSpecialValueType 中 MapValuer Implements 检查分支
func TestAnyToMapWithNonTimeMapValuer(t *testing.T) {
	type Canvas struct {
		Name   string           `json:"name"`
		Origin NonTimeMapValuer `json:"origin"`
	}

	obj := Canvas{
		Name:   "canvas1",
		Origin: NonTimeMapValuer{X: 10, Y: 20},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, "canvas1", result["name"])
	// NonTimeMapValuer 实现了 MapValuer，应该调用 MapValue() 返回自定义值
	expected := map[string]int{"x": 10, "y": 20}
	assert.Equal(t, expected, result["origin"])
}

// TestAnyToMapWithNonTimeStringer 测试非 time.Time 基础的 struct 实现 fmt.Stringer
// 覆盖 isSpecialValueType 中 Stringer Implements 检查分支
func TestAnyToMapWithNonTimeStringer(t *testing.T) {
	type Person struct {
		Age      int             `json:"age"`
		FullName NonTimeStringer `json:"full_name"`
	}

	obj := Person{
		Age:      30,
		FullName: NonTimeStringer{First: "John", Last: "Doe"},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, 30, result["age"])
	// NonTimeStringer 实现了 Stringer，应该被视为值类型
	_, isMap := result["full_name"].(map[string]interface{})
	assert.False(t, isMap, "full_name should NOT be a map when type implements Stringer")
	// 保留原始值（NonTimeStringer 类型本身）
	val, ok := result["full_name"].(NonTimeStringer)
	assert.True(t, ok, "full_name should be NonTimeStringer type")
	assert.Equal(t, "John", val.First)
}

// TestIsSpecialValueTypeNonStruct 测试 isSpecialValueType 非 struct 分支
func TestIsSpecialValueTypeNonStruct(t *testing.T) {
	// int → false
	assert.False(t, isSpecialValueType(reflect.ValueOf(42)))

	// string → false
	assert.False(t, isSpecialValueType(reflect.ValueOf("hello")))

	// bool → false
	assert.False(t, isSpecialValueType(reflect.ValueOf(true)))
}

// TestIsUnderlyingTimeTypeWithArrayAndChan 测试 array 和 chan 类型
func TestIsUnderlyingTimeTypeWithArrayAndChan(t *testing.T) {
	// [1]DateTime → true (array 元素检查)
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf([1]DateTime{})))

	// chan DateTime → true (chan 元素检查)
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf(make(chan DateTime))))

	// [2]int → false
	assert.False(t, isUnderlyingTimeType(reflect.TypeOf([2]int{})))

	// map[DateTime]string → true (map 键检查)
	assert.True(t, isUnderlyingTimeType(reflect.TypeOf(map[DateTime]string{})))

	// map[string]int → false
	assert.False(t, isUnderlyingTimeType(reflect.TypeOf(map[string]int{})))
}

// TestAnyToMapSliceWithNonTimeMapValuer 测试切片中包含非 time 基础的 MapValuer 类型
func TestAnyToMapSliceWithNonTimeMapValuer(t *testing.T) {
	type Canvas struct {
		Name   string             `json:"name"`
		Points []NonTimeMapValuer `json:"points"`
	}

	obj := Canvas{
		Name:   "shape",
		Points: []NonTimeMapValuer{{X: 1, Y: 2}, {X: 3, Y: 4}},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	points, ok := result["points"].([]any)
	assert.True(t, ok)
	assert.Len(t, points, 2)
	// MapValuer 类型在切片中也应该被视为特殊类型
	_, isMap := points[0].(map[string]interface{})
	assert.False(t, isMap, "NonTimeMapValuer in slice should NOT be recursively expanded")
}

// TestAnyToMapSliceWithNonTimeStringer 测试切片中包含非 time 基础的 Stringer 类型
func TestAnyToMapSliceWithNonTimeStringer(t *testing.T) {
	type Team struct {
		Name    string            `json:"name"`
		Members []NonTimeStringer `json:"members"`
	}

	obj := Team{
		Name:    "team1",
		Members: []NonTimeStringer{{First: "John", Last: "Doe"}, {First: "Jane", Last: "Smith"}},
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	members, ok := result["members"].([]any)
	assert.True(t, ok)
	assert.Len(t, members, 2)
	// Stringer 类型在切片中应该被视为特殊类型
	_, isMap := members[0].(map[string]interface{})
	assert.False(t, isMap, "NonTimeStringer in slice should NOT be recursively expanded")
}

// TestAnyToMapWithMultipleAnonymousStructs 测试多层匿名结构体嵌套
func TestAnyToMapWithMultipleAnonymousStructsPtr(t *testing.T) {
	type Base struct {
		ID int `json:"id"`
	}
	type Middle struct {
		*Base
		Name string `json:"name"`
	}
	type Top struct {
		*Middle
		Extra string `json:"extra"`
	}

	obj := Top{
		Middle: &Middle{
			Base: &Base{ID: 1},
			Name: "test",
		},
		Extra: "info",
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, 1, result["id"])
	assert.Equal(t, "test", result["name"])
	assert.Equal(t, "info", result["extra"])
}

// TestAnyToMapWithMultipleAnonymousStructs 测试多层匿名结构体嵌套
func TestAnyToMapWithMultipleAnonymousStructsPtr2(t *testing.T) {
	type Base struct {
		ID int `json:"id"`
	}
	type Middle struct {
		B    *Base  `json:"base"`
		Name string `json:"name"`
	}

	obj := &Middle{
		B:    nil, // Base 是 nil，测试指针解包时的健壮性
		Name: "test",
	}

	result, err := AnyToMap(obj)
	assert.Nil(t, err)
	assert.Equal(t, nil, result["base"])
	assert.Equal(t, "test", result["name"])
}
