package xtransform

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	replaceWord, _ = regexp.Compile(`@\{\w+\}|@\w+`)
)

//Translate 模板转换字符串
//eg1: str="aaa:bbb:@ccc"  data={ccc:100}  ==> "aaa:bbb:100"
//eg2: str="aaa:bbb:@{ccc}"  data={ccc:200}  ==> "aaa:bbb:200"
//eg3: str="aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> "aaa:abcd:200"
//eg4: str="@aaa:@bbb:@{ccc}"  data={bbb:"abcd",ccc:200}  ==> ":abcd:200"
//
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

func TranslateObject(tpl string, data interface{}) string {
	bytes, _ := json.Marshal(data)
	mapData := map[string]interface{}{}
	json.Unmarshal(bytes, &mapData)
	return Translate(tpl, mapData)
}
