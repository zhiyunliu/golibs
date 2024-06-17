package xtransform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zhiyunliu/golibs/xreflect"
)

var (
	replaceWord, _ = regexp.Compile(`@\{\w+\}|@\w+`)
)

// Translate 模板转换字符串
// eg1: str="aaa:bbb:@ccc"  data={ccc:100}  ==> "aaa:bbb:100"
// eg2: str="aaa:bbb:@{ccc}"  data={ccc:200}  ==> "aaa:bbb:200"
// eg3: str="aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> "aaa:abcd:200"
// eg4: str="@aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> ":abcd:200"
func Translate(tpl string, data map[string]interface{}) string {
	return replaceWord.ReplaceAllStringFunc(tpl, func(prop string) string {
		/*
			@{aa1}
			@aa1
		*/
		key := prop[1:]
		if strings.Contains(prop, "{") {
			key = prop[2 : len(prop)-1]
		}
		if data[key] != nil {
			return fmt.Sprint(data[key])
		}
		return ""
	})
}

func TranslateMap(tpl string, data map[string]string) string {
	return replaceWord.ReplaceAllStringFunc(tpl, func(prop string) string {
		/*
			@{aa1}
			@aa1
		*/
		key := prop[1:]
		if strings.Contains(prop, "{") {
			key = prop[2 : len(prop)-1]
		}
		return data[key]
	})
}

// Deprecated: As of Go v0.2.0, this function simply calls [GenericTranslate].
func TranslateObject(tpl string, data interface{}) string {
	return GenericTranslate(tpl, data)
}

func GenericTranslate[T any](tpl string, data T) string {
	mapVal, err := xreflect.AnyToMap(data, xreflect.WithMaxDepth(1))
	if err != nil {
		return ""
	}
	return Translate(tpl, mapVal)
}
