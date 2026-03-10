package xreflect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type MyStruct struct {
	BoolField bool
	IntField1 int
	IntField2 int
	IntField3 int
}

func Test_ReflectSet(t *testing.T) {

	// 创建结构体实例
	var myInstance *MyStruct = new(MyStruct)

	// 创建字段名称和值的映射
	fieldValues := map[string]interface{}{
		"BoolField": true,
		"IntField1": 42,
		"IntField2": 123,
		"IntField3": 987,
	}

	json.Unmarshal([]byte(`{}`), &fieldValues)

	// 使用反射循环为结构体字段赋值
	for fieldName, value := range fieldValues {
		setStructField(myInstance, fieldName, value)
	}
	// 输出赋值后的结构体实例
	fmt.Println(myInstance)
}

// setStructField 使用反射为结构体字段赋值
func setStructField(myStructPtr interface{}, fieldName string, value interface{}) {
	refval := reflect.ValueOf(myStructPtr)
	if refval.IsNil() {
		refval = reflect.New(refval.Type())
	}

	structValue := reflect.Indirect(refval)

	// 获取字段的 reflect.Value
	fieldValue := structValue.FieldByName(fieldName)

	// 如果字段存在且可设置，则设置字段的值
	if fieldValue.IsValid() && fieldValue.CanSet() {
		// 将传入的值转换为 reflect.Value
		newValue := reflect.ValueOf(value)

		// 确保值的类型与字段类型匹配
		if newValue.Type().AssignableTo(fieldValue.Type()) {
			// 设置字段的值
			fieldValue.Set(newValue)
		} else {
			fmt.Printf("Type mismatch for field %s\n", fieldName)
		}
	} else {
		fmt.Printf("Field %s not found or not settable\n", fieldName)
	}
}

func Test_Map(t *testing.T) {
	type val struct {
		Map *map[string]any `json:"map"`
	}
	result := &val{}
	reflectVal := reflect.ValueOf(result)

	fields := CachedTypeFields(reflectVal.Type())

	for i := range fields.List {
		ftype := fields.List[i].typ

		mapval := reflect.MakeMap(reflect.MapOf(ftype.Key(), ftype.Elem()))
		rv1 := reflect.New(ftype)
		mapval.SetMapIndex(reflect.ValueOf("aaaa"), reflect.ValueOf("bbb"))
		rv1.Elem().Set(mapval)

		fv := GetRealReflectVal(&fields.List[i], reflectVal)
		fv.Set(rv1)

	}

	if !reflect.DeepEqual(*result.Map, map[string]any{"aaaa": "bbb"}) {
		t.Errorf("反射失败:%+v", result.Map)
	}

}

func Test_Map2(t *testing.T) {
	type val struct {
		Map map[string]any `json:"map"`
	}
	result := &val{}
	reflectVal := reflect.ValueOf(result)

	fields := CachedTypeFields(reflectVal.Type())

	for i := range fields.List {
		ftype := fields.List[i].typ

		mapval := reflect.MakeMap(reflect.MapOf(ftype.Key(), ftype.Elem()))
		rv1 := reflect.New(ftype)
		mapval.SetMapIndex(reflect.ValueOf("aaaa"), reflect.ValueOf("bbb"))
		rv1.Elem().Set(mapval)

		fv := GetRealReflectVal(&fields.List[i], reflectVal)

		fv.Set(mapval)

	}

	if !reflect.DeepEqual(result.Map, map[string]any{"aaaa": "bbb"}) {
		t.Errorf("反射失败:%+v", result.Map)
	}

}

func Test_GetFieldType(t *testing.T) {
	type sample struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	fields := CachedTypeFields(reflect.TypeOf(sample{}))

	// existing field
	typ, ok := fields.GetFieldType("name")
	if !ok || typ.Kind() != reflect.String {
		t.Errorf("GetFieldType(name) = %v, %v; want string, true", typ, ok)
	}

	typ2, ok2 := fields.GetFieldType("count")
	if !ok2 || typ2.Kind() != reflect.Int {
		t.Errorf("GetFieldType(count) = %v, %v; want int, true", typ2, ok2)
	}

	// non-existing field
	_, ok3 := fields.GetFieldType("nonexist")
	if ok3 {
		t.Error("GetFieldType(nonexist) should return false")
	}
}

func Test_StructFields_Dencode(t *testing.T) {
	type sample struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	val := sample{}
	rv := reflect.ValueOf(&val).Elem()
	fields := CachedTypeFields(rv.Type())

	// Dencode existing field
	err := fields.Dencode(rv, "name", "hello")
	if err != nil {
		t.Errorf("Dencode(name) error = %v", err)
	}
	if val.Name != "hello" {
		t.Errorf("Dencode(name) = %v, want hello", val.Name)
	}

	// Dencode non-existing field (should not error)
	err = fields.Dencode(rv, "nonexist", "world")
	if err != nil {
		t.Errorf("Dencode(nonexist) error = %v", err)
	}

	// Dencode int field
	err = fields.Dencode(rv, "count", 42)
	if err != nil {
		t.Errorf("Dencode(count) error = %v", err)
	}
	if val.Count != 42 {
		t.Errorf("Dencode(count) = %v, want 42", val.Count)
	}
}

func Test_field_Dencoder(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}

	val := sample{}
	rv := reflect.ValueOf(&val).Elem()
	fields := CachedTypeFields(rv.Type())

	f := fields.ExactName["name"]
	err := f.Dencoder(rv, "world")
	if err != nil {
		t.Errorf("Dencoder error = %v", err)
	}
	if val.Name != "world" {
		t.Errorf("Dencoder = %v, want world", val.Name)
	}
}

func Test_field_Encoder(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}
	val := sample{Name: "hello"}
	rv := reflect.ValueOf(&val).Elem()
	fields := CachedTypeFields(rv.Type())

	f := fields.ExactName["name"]
	v, ok := f.Encoder(rv, StructOptions{MaxDepth: 5})
	if !ok {
		t.Error("Encoder should return ok=true")
	}
	if v != "hello" {
		t.Errorf("Encoder = %v, want hello", v)
	}
}

func Test_field_Encoder_exceedDepth(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}
	val := sample{Name: "hello"}
	rv := reflect.ValueOf(&val).Elem()
	fields := CachedTypeFields(rv.Type())

	f := fields.ExactName["name"]
	// MaxDepth=0, curDepth will become 1 after increment -> exceeds
	_, ok := f.Encoder(rv, StructOptions{MaxDepth: 0})
	if ok {
		t.Error("Encoder should return ok=false when depth exceeded")
	}
}

func Test_byIndex_sort(t *testing.T) {
	fields := byIndex{
		{Name: "b", Index: []int{1, 0}},
		{Name: "a", Index: []int{0}},
		{Name: "c", Index: []int{1, 1}},
	}

	if fields.Len() != 3 {
		t.Errorf("Len() = %d, want 3", fields.Len())
	}

	fields.Swap(0, 1)
	if fields[0].Name != "a" || fields[1].Name != "b" {
		t.Errorf("Swap failed: %v", fields)
	}

	// Test Less
	if !fields.Less(0, 1) { // [0] < [1,0]
		t.Error("Less(0,1) should be true")
	}
	if fields.Less(1, 0) {
		t.Error("Less(1,0) should be false")
	}
}

func Test_GetRealReflectVal_nilPtr(t *testing.T) {
	type Inner struct {
		Value string `json:"value"`
	}
	type Outer struct {
		Inner *Inner `json:"inner"`
	}

	val := Outer{Inner: nil}
	rv := reflect.ValueOf(&val).Elem()
	fields := CachedTypeFields(rv.Type())
	f := fields.ExactName["value"]
	if f != nil {
		// value is inside Inner which is a named (non-anonymous) field
		// It shouldn't be accessible as a top-level field
	}

	// Test the field directly with inner
	innerField := fields.ExactName["inner"]
	if innerField != nil {
		subv := GetRealReflectVal(innerField, rv)
		if subv.IsValid() && subv.IsNil() {
			// inner pointer is nil, which is expected
		}
	}
}

func Test_CachedTypeFields_pointer(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
	}
	// Pass pointer type
	fields := CachedTypeFields(reflect.TypeOf(&sample{}))
	if _, ok := fields.ExactName["name"]; !ok {
		t.Error("CachedTypeFields should unwrap pointer type")
	}
}

func Test_getFieldTag_xdb(t *testing.T) {
	type sample struct {
		Name string `xdb:"xdb_name"`
	}
	sf := reflect.TypeOf(sample{}).Field(0)
	tag := getFieldTag(sf.Tag)
	if tag != "xdb_name" {
		t.Errorf("getFieldTag(xdb) = %v, want xdb_name", tag)
	}
}

func Test_getFieldTag_db(t *testing.T) {
	type sample struct {
		Name string `db:"db_name"`
	}
	sf := reflect.TypeOf(sample{}).Field(0)
	tag := getFieldTag(sf.Tag)
	if tag != "db_name" {
		t.Errorf("getFieldTag(db) = %v, want db_name", tag)
	}
}

func Test_getFieldTag_json(t *testing.T) {
	type sample struct {
		Name string `json:"json_name"`
	}
	sf := reflect.TypeOf(sample{}).Field(0)
	tag := getFieldTag(sf.Tag)
	if tag != "json_name" {
		t.Errorf("getFieldTag(json) = %v, want json_name", tag)
	}
}

func Test_getFieldTag_priority(t *testing.T) {
	type sample struct {
		Name string `xdb:"xdb_name" db:"db_name" json:"json_name"`
	}
	sf := reflect.TypeOf(sample{}).Field(0)
	tag := getFieldTag(sf.Tag)
	if tag != "xdb_name" {
		t.Errorf("getFieldTag(priority) = %v, want xdb_name (xdb > db > json)", tag)
	}
}

func Test_getFieldTag_empty(t *testing.T) {
	type sample struct {
		Name string
	}
	sf := reflect.TypeOf(sample{}).Field(0)
	tag := getFieldTag(sf.Tag)
	if tag != "" {
		t.Errorf("getFieldTag(empty) = %v, want empty", tag)
	}
}

// func Test_Map3(t *testing.T) {
// 	type val struct {
// 		Map map[string]any `json:"xxmap"`
// 	}
// 	result := &val{
// 		Map: map[string]any{"aaaa": "bbb"},
// 	}
// 	refval := reflect.ValueOf(result)

// 	fields := CachedTypeFields(refval.Type())

// 	mapresult := map[string]any{}

// 	for _, f := range fields.ExactName {
// 		if val, ok := f.Encoder(refval); ok {
// 			mapresult[f.Name] = val
// 		}
// 	}
// 	expectMap := map[string]any{
// 		"xxmap": `{"aaaa": "bbb"}`,
// 	}

// 	if !reflect.DeepEqual(expectMap, mapresult) {
// 		t.Errorf("反射失败:result:%+v, expect:%+v", mapresult, expectMap)
// 	}

// }
