package xreflect

import (
	"fmt"
	"reflect"
)

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
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	num := value.NumField()
	for i := 0; i < num; i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)
		fieldName := fieldType.Tag.Get("json")
		if fieldName == "" {
			fieldName = fieldType.Name
		}
		if fieldName == "-" {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			if fieldType.Anonymous {
				// Merge the fields of the anonymous struct
				for k, v := range structToMapDepth(field.Interface(), depth+1, maxDepth) {
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
				if elem.Kind() == reflect.Struct {
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
