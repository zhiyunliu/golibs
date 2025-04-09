package xtransform

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/zhiyunliu/golibs/xreflect"
)

// Translate 模板转换字符串
// eg1: str="aaa:bbb:@ccc"  data={ccc:100}  ==> "aaa:bbb:100"
// eg2: str="aaa:bbb:@{ccc}"  data={ccc:200}  ==> "aaa:bbb:200"
// eg3: str="aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> "aaa:abcd:200"
// eg4: str="@aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> ":abcd:200"
func Translate(tpl string, data map[string]interface{}) string {
	return xTranslate(tpl, data)
}

func TranslateMap(tpl string, data map[string]string) string {
	return xTranslate(tpl, data)
}

// GenericTranslate 泛型模板转换字符串
func GenericTranslate[T any](tpl string, data T) string {
	var tmpData any = data

	if v, ok := tmpData.(map[string]any); ok {
		return xTranslate(tpl, v)
	}

	if v, ok := tmpData.(map[string]string); ok {
		return xTranslate(tpl, v)
	}

	mapVal, err := xreflect.AnyToMap(data, xreflect.WithMaxDepth(1))
	if err != nil {
		return ""
	}
	return xTranslate(tpl, mapVal)
}

// toString 高性能类型转换（覆盖常见类型）
func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(x).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(x).Uint(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(x).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprint(x) // 兜底处理
	}
}

func xTranslate[T any](template string, data map[string]T) string {
	var builder strings.Builder
	var keyBuilder strings.Builder
	inVar := false
	inBrace := false
	cs := byte('a')
	ucs := byte('A')
	ce := byte('z')
	uce := byte('Z')

	for i := 0; i < len(template); i++ {
		c := template[i]

		switch {
		case c == '@' && !inVar:
			// 开始变量
			inVar = true
			if i+1 < len(template) && template[i+1] == '{' {
				inBrace = true
				i++ // 跳过'{'
			}

		case inVar && inBrace && c == '}':
			// 结束@{var}模式
			key := keyBuilder.String()
			if value, ok := data[key]; ok {
				builder.WriteString(toString(value))
			}
			// 不存在的值不做任何输出（即替换为空字符串）
			keyBuilder.Reset()
			inVar = false
			inBrace = false

		case inVar && !inBrace && !((cs <= c && c <= ce) || (ucs <= c && c <= uce)):
			// 结束@var模式（遇到分隔符）
			key := keyBuilder.String()
			if value, ok := data[key]; ok {
				builder.WriteString(toString(value))
			}
			// 不存在的值不做任何输出
			keyBuilder.Reset()
			inVar = false
			builder.WriteByte(c) // 写入当前字符

		case inVar && !inBrace && i == len(template)-1:
			// 字符串末尾的@var模式
			keyBuilder.WriteByte(c)
			key := keyBuilder.String()
			if value, ok := data[key]; ok {
				builder.WriteString(toString(value))
			}
			// 不存在的值不做任何输出
			keyBuilder.Reset()
			inVar = false

		case inVar:
			// 收集变量名
			keyBuilder.WriteByte(c)

		default:
			// 普通字符
			builder.WriteByte(c)
		}
	}

	return builder.String()
}
