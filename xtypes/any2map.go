package xtypes

import (
	"reflect"
)

type Option func(*options)
type options struct {
	maxDepth int
}

func WithMaxDepth(maxDepth int) Option {
	return func(o *options) {
		o.maxDepth = maxDepth
	}
}

// 将value 转换为map
func AnyToMap(obj any, opts ...Option) map[string]any {
	options := options{
		maxDepth: 100, // default max depth
	}
	for _, o := range opts {
		o(&options)
	}

	return structToMapDepth(obj, 0, options.maxDepth)
}

func structToMapDepth(obj any, depth int, maxDepth int) map[string]any {
	if depth > maxDepth {
		return nil
	}

	result := make(map[string]any)
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
