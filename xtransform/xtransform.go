package xtransform

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/zhiyunliu/golibs/xreflect"
)

// translateOptions holds configuration for the translate function
type translateOptions struct {
	enableAtMode      bool // @aaa
	enableAtBraceMode bool // @{aaa}
	enableBraceMode   bool // {aaa}
	explicit          bool // 是否显式设置了选项
}

// TranslateOption defines a function to configure translate options
type TranslateOption func(*translateOptions)

// WithAtMode enables @aaa pattern
func WithAtMode() TranslateOption {
	return func(o *translateOptions) {
		o.enableAtMode = true
		o.explicit = true
	}
}

// WithAtBraceMode enables @{aaa} pattern (default enabled)
func WithAtBraceMode() TranslateOption {
	return func(o *translateOptions) {
		o.enableAtBraceMode = true
		o.explicit = true
	}
}

// WithBraceMode enables {aaa} pattern
func WithBraceMode() TranslateOption {
	return func(o *translateOptions) {
		o.enableBraceMode = true
		o.explicit = true
	}
}

// getOptions 获取翻译选项，如果没有显式指定，则使用向后兼容的默认值
func getOptions(opts []TranslateOption) translateOptions {
	defaultOpts := translateOptions{
		enableAtMode:      false, // 默认不启用 @var 模式
		enableAtBraceMode: false, // 默认不启用 @{var} 模式（除非显式调用选项）
		enableBraceMode:   false, // 默认不启用 {var} 模式
		explicit:          false,
	}

	// 应用用户提供的选项
	for _, opt := range opts {
		opt(&defaultOpts)
	}

	// 如果没有显式设置选项，则启用所有模式以保持向后兼容性
	if !defaultOpts.explicit {
		defaultOpts.enableAtMode = true
		defaultOpts.enableAtBraceMode = true
		defaultOpts.enableBraceMode = false
	} else {
		// 如果显式设置了选项但没有设置任何一种模式，则默认启用AtBraceMode
		if !defaultOpts.enableAtMode && !defaultOpts.enableAtBraceMode && !defaultOpts.enableBraceMode {
			defaultOpts.enableAtBraceMode = true
		}
	}

	return defaultOpts
}

// Translate 模板转换字符串
// eg1: str="aaa:bbb:@ccc"  data={ccc:100}  ==> "aaa:bbb:100"
// eg2: str="aaa:bbb:@{ccc}"  data={ccc:200}  ==> "aaa:bbb:200"
// eg3: str="aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> "aaa:abcd:200"
// eg4: str="@aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> ":abcd:200"
func Translate(tpl string, data map[string]interface{}, opts ...TranslateOption) string {
	optsObj := getOptions(opts)
	return xTranslateWithOptions(tpl, data, optsObj)
}

func TranslateMap(tpl string, data map[string]string, opts ...TranslateOption) string {
	optsObj := getOptions(opts)
	return xTranslateWithOptions(tpl, data, optsObj)
}

// TranslateCallback 使用回调函数进行模板转换
// eg: str="aaa:@bbb:@{ccc}" callback=function(key) { return strings.ToUpper(key) } ==> "aaa:BBB:CCC"
func TranslateCallback(tpl string, callback func(param string) string, opts ...TranslateOption) string {
	optsObj := getOptions(opts)
	return translateWithCallbackWithOptions(tpl, callback, optsObj)
}

// GenericTranslate 泛型模板转换字符串
func GenericTranslate[T any](tpl string, data T, opts ...TranslateOption) string {
	var tmpData any = data

	optsObj := getOptions(opts)

	if v, ok := tmpData.(map[string]any); ok {
		return xTranslateWithOptions(tpl, v, optsObj)
	}

	if v, ok := tmpData.(map[string]string); ok {
		return xTranslateWithOptions(tpl, v, optsObj)
	}

	mapVal, err := xreflect.AnyToMap(data, xreflect.WithMaxDepth(1))
	if err != nil {
		return ""
	}
	return xTranslateWithOptions(tpl, mapVal, optsObj)
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

func xTranslateWithOptions[T any](template string, data map[string]T, opts translateOptions) string {
	return translateWithCallbackWithOptions(template, func(key string) string {
		if value, ok := data[key]; ok {
			return toString(value)
		}
		return ""
	}, opts)
}

// translateWithCallbackWithOptions 使用回调函数处理模板的核心逻辑，带选项控制
func translateWithCallbackWithOptions(template string, callback func(param string) string, opts translateOptions) string {
	var builder strings.Builder
	var keyBuilder strings.Builder
	inVar := false
	inBrace := false
	inSimpleBrace := false // 标记是否处于简单大括号模式
	cs := byte('a')
	ucs := byte('A')
	ce := byte('z')
	uce := byte('Z')

	for i, cnt := 0, len(template); i < cnt; i++ {
		c := template[i]

		switch {
		case c == '@' && !inVar && opts.enableAtMode && i+1 < cnt && template[i+1] != '{':
			// 开始 @var 模式（如果下一个字符不是{）
			inVar = true

		case c == '@' && !inVar && opts.enableAtBraceMode && i+1 < cnt && template[i+1] == '{':
			// 开始 @{var} 模式
			inVar = true
			inBrace = true
			i++ // 跳过'{'

		case c == '{' && !inVar && opts.enableBraceMode:
			// 开始简单大括号模式 {variable}
			inVar = true
			inSimpleBrace = true

		case inVar && inBrace && c == '}':
			// 结束@{var}模式
			key := keyBuilder.String()
			value := callback(key)
			builder.WriteString(value)
			keyBuilder.Reset()
			inVar = false
			inBrace = false

		case inVar && inSimpleBrace && c == '}':
			// 结束{var}模式
			key := keyBuilder.String()
			value := callback(key)
			builder.WriteString(value)
			keyBuilder.Reset()
			inVar = false
			inSimpleBrace = false

		case inVar && !inBrace && !inSimpleBrace && !((cs <= c && c <= ce) || (ucs <= c && c <= uce)) && opts.enableAtMode:
			// 结束@var模式（遇到分隔符）
			key := keyBuilder.String()
			value := callback(key)
			builder.WriteString(value)
			keyBuilder.Reset()
			inVar = false
			builder.WriteByte(c) // 写入当前字符

		case inVar && !inBrace && !inSimpleBrace && i == cnt-1 && opts.enableAtMode:
			// 字符串末尾的@var模式
			keyBuilder.WriteByte(c)
			key := keyBuilder.String()
			value := callback(key)
			builder.WriteString(value)
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
