package xreflect

import (
	"fmt"
	"reflect"
)

// MapValuer 是一个接口，类型可以实现它来控制在转换为map时的表现。
// 实现此接口的类型将被视为值类型，不会被递归分解为其内部字段。
// 例如，对于 type DateTime time.Time 这样的类型定义，
// 可以通过实现 MapValuer 接口来控制其在map转换时的值表示。

type StructOption func(*structOptions)
type structOptions struct {
	curDepth             int
	maxDepth             int
	structField          bool
	sliceItem            bool
	disableJSONMarshaler bool
}

func (o structOptions) IsValidDepth() bool {
	return o.curDepth < o.maxDepth
}
func WithMaxDepth(maxDepth int) StructOption {
	return func(o *structOptions) {
		o.maxDepth = maxDepth
	}
}

func WithStructField(state bool) StructOption {
	return func(o *structOptions) {
		o.structField = state
	}
}

func WithSliceItem(state bool) StructOption {
	return func(o *structOptions) {
		o.sliceItem = state
	}
}

func WithDisableJSONMarshaler(state bool) StructOption {
	return func(o *structOptions) {
		o.disableJSONMarshaler = state
	}
}

// 将value 转换为map
func AnyToMap(value interface{}, opts ...StructOption) (map[string]any, error) {
	options := structOptions{
		curDepth:             0,
		maxDepth:             10, // default max depth
		sliceItem:            true,
		structField:          true,
		disableJSONMarshaler: true,
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
		fields := CachedTypeFields(val.Type())
		for _, f := range fields.ExactName {
			if val, ok := f.Encoder(val, options); ok {
				result[f.Name] = val
			}
		}
		return result, nil
	default:
		return result, fmt.Errorf("unsupported input type: %v", val.Kind())
	}
	return result, nil
}
