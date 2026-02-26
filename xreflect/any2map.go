package xreflect

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	timeType          = reflect.TypeOf(time.Time{})
	mapValuerType     = reflect.TypeOf((*MapValuer)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// MapValuer 是一个接口，类型可以实现它来控制在转换为map时的表现。
// 实现此接口的类型将被视为值类型，不会被递归分解为其内部字段。
// 例如，对于 type DateTime time.Time 这样的类型定义，
// 可以通过实现 MapValuer 接口来控制其在map转换时的值表示。
type MapValuer interface {
	MapValue() any
}

type StructOption func(*structOptions)
type structOptions struct {
	maxDepth int
}

func WithMaxDepth(maxDepth int) StructOption {
	return func(o *structOptions) {
		o.maxDepth = maxDepth
	}
}

// 将value 转换为map
func AnyToMap(value interface{}, opts ...StructOption) (map[string]any, error) {
	options := structOptions{
		maxDepth: 10, // default max depth
	}
	for _, o := range opts {
		o(&options)
	}

	result := make(map[string]interface{})

	// 使用反射获取值的类型信息
	val := reflect.ValueOf(value)

	// 如果是指针，检查是否为nil，并获取其指向的值
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result, nil // 如果是nil指针，直接返回空的map
		}
		val = val.Elem()
	}

	// 检查值的类型
	switch val.Kind() {
	case reflect.Map:
		// 如果是map类型，遍历键值对并转换为map[string]interface{}
		keys := val.MapKeys()
		for _, key := range keys {
			result[key.String()] = val.MapIndex(key).Interface()
		}
	case reflect.Struct:
		// 如果是struct类型，遍历字段并使用json标签作为键
		result = structToMapDepth(value, 0, options.maxDepth)
	default:
		return result, fmt.Errorf("unsupported input type: %v", val.Kind())
	}
	return result, nil
}

func structToMapDepth(obj interface{}, depth int, maxDepth int) map[string]interface{} {
	if depth >= maxDepth {
		return nil
	}

	result := make(map[string]interface{})
	value := reflect.ValueOf(obj)
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return result
		}
		value = value.Elem()
	}
	num := value.NumField()
	for i := 0; i < num; i++ {
		field := value.Field(i)

		fieldIsNull := false
		for field.Kind() == reflect.Ptr {
			if field.IsNil() {
				fieldIsNull = true
				break
			}
			field = field.Elem()
		}

		fieldType := value.Type().Field(i)
		fieldName := fieldType.Tag.Get("json")

		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}
		if fieldName == "-" {
			continue
		}

		if fieldName == "" {
			fieldName = fieldType.Name
		}
		if fieldIsNull {
			result[fieldName] = nil
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			// 检查是否实现了 MapValuer 接口，优先使用其返回值
			if mv, ok := getMapValuer(field); ok {
				result[fieldName] = mv.MapValue()
			} else if isSpecialValueType(field) {
				// 检查是否为特殊类型，如果是则直接作为值处理
				result[fieldName] = field.Interface()
			} else if fieldType.Anonymous {
				// Merge the fields of the anonymous struct
				anonymous := structToMapDepth(field.Interface(), depth+1, maxDepth)
				for k, v := range anonymous {
					result[k] = v
				}
			} else {
				result[fieldName] = structToMapDepth(field.Interface(), depth+1, maxDepth)
			}
		case reflect.Slice:
			field_len := field.Len()
			sliceResult := make([]any, field_len)
			for j := 0; j < field_len; j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.Struct && !isSpecialValueType(elem) {
					sliceResult[j] = structToMapDepth(elem.Interface(), depth+1, maxDepth)
				} else {
					sliceResult[j] = elem.Interface()
				}
			}
			result[fieldName] = sliceResult
		default:
			result[fieldName] = field.Interface()
		}
	}

	return result
}

// getMapValuer 检查值是否实现了 MapValuer 接口（值接收者或指针接收者）
func getMapValuer(field reflect.Value) (MapValuer, bool) {
	// 检查值接收者
	if mv, ok := field.Interface().(MapValuer); ok {
		return mv, true
	}
	// 检查指针接收者
	if field.CanAddr() {
		if mv, ok := field.Addr().Interface().(MapValuer); ok {
			return mv, true
		}
	}
	return nil, false
}

// 检查是否是应该被视为值类型处理的特殊类型
// 通过接口实现检查和底层类型检查来判断
func isSpecialValueType(field reflect.Value) bool {
	if field.Kind() != reflect.Struct {
		// 检查非结构体类型，如 type MyTime time.Time
		return isUnderlyingTimeType(field.Type())
	}

	// 获取字段的原始类型
	originalType := field.Type()
	ptrType := reflect.PointerTo(originalType)

	// 检查是否是 time.Time 类型或可以转换为 time.Time 的类型别名
	if originalType == timeType || originalType.ConvertibleTo(timeType) {
		return true
	}

	// 检查是否实现了 MapValuer 接口（值接收者或指针接收者）
	if originalType.Implements(mapValuerType) || ptrType.Implements(mapValuerType) {
		return true
	}

	// 检查是否实现了 stringer 接口 (fmt.Stringer)（值接收者或指针接收者）
	// 如果实现了，通常意味着它有一个有意义的整体字符串表示
	if originalType.Implements(stringerType) || ptrType.Implements(stringerType) {
		return true
	}

	// 检查是否实现了 encoding.TextMarshaler 接口（值接收者或指针接收者）
	if originalType.Implements(textMarshalerType) || ptrType.Implements(textMarshalerType) {
		return true
	}

	return false
}

// 获取类型的最底层类型（递归解包嵌套的类型定义）
// 使用 ConvertibleTo 进行更准确的类型别名检测
func getUnderlyingType(t reflect.Type) reflect.Type {
	// 处理指针类型
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 检查是否是 time.Time 类型或可以转换为 time.Time 的类型
	timeType := reflect.TypeOf(time.Time{})
	if t.Kind() == reflect.Struct && t.ConvertibleTo(timeType) {
		return timeType
	}

	return t
}

// 检查类型是否基于 time.Time 的类型别名
func isUnderlyingTimeType(t reflect.Type) bool {
	breakLoop := false
	// 获取底层类型
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Array || t.Kind() == reflect.Slice || t.Kind() == reflect.Chan || t.Kind() == reflect.Map {
		switch t.Kind() {
		case reflect.Ptr:
			t = t.Elem()
		case reflect.Array, reflect.Slice, reflect.Chan:
			t = t.Elem()
		case reflect.Map:
			// 对于 map 类型，检查键和值的类型
			return isUnderlyingTimeType(t.Key()) || isUnderlyingTimeType(t.Elem())
		default:
			breakLoop = true
		}
		if breakLoop {
			break
		}
	}

	// 对于结构体类型，使用 getUnderlyingType
	if t.Kind() == reflect.Struct {
		return getUnderlyingType(t).String() == "time.Time"
	}

	// 对于非结构体类型，检查其原始类型是否可以转换为 time.Time
	// 比如 type DateTime time.Time 这种情况
	return t.String() != "time.Time" &&
		t.ConvertibleTo(reflect.TypeOf(time.Time(time.Unix(0, 0))))
}
