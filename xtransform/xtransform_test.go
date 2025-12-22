package xtransform

import (
	"strings"
	"testing"
)

func TestTranslate(t *testing.T) {

	tests := []struct {
		name string
		tpl  string
		data map[string]interface{}
		want string
	}{
		{name: "1.num.", tpl: "aaa:bbb:@ccc", data: map[string]interface{}{"ccc": 100}, want: "aaa:bbb:100"},
		{name: "2.str", tpl: "aaa:bbb:@ccc", data: map[string]interface{}{"ccc": "100"}, want: "aaa:bbb:100"},
		{name: "2.str", tpl: "aaa:bbb:@{ccc}", data: map[string]interface{}{"ccc": "200"}, want: "aaa:bbb:200"},
		{name: "2.str", tpl: "aaa:bbb:@{ccc}", data: map[string]interface{}{"ccc": 200}, want: "aaa:bbb:200"},
		{name: "2.str", tpl: "aaa:@bbb:@{ccc}", data: map[string]interface{}{"bbb": "abcd", "ccc": "200"}, want: "aaa:abcd:200"},
		{name: "2.str", tpl: "@aaa:@bbb:@{ccc}", data: map[string]interface{}{"bbb": "abcd", "ccc": "300"}, want: ":abcd:300"},
		// 新增测试用例：支持{variable}格式
		{name: "simple brace", tpl: "aaa:bbb:{ccc}", data: map[string]interface{}{"ccc": 100}, want: "aaa:bbb:100"},
		{name: "mixed formats", tpl: "aaa:{bbb}:@{ccc}", data: map[string]interface{}{"bbb": "xyz", "ccc": 200}, want: "aaa:xyz:200"},
		{name: "mixed formats2", tpl: "aaa:@bbb:{ccc}", data: map[string]interface{}{"bbb": "xyz", "ccc": 300}, want: "aaa:xyz:300"},
		{name: "special-char-1", tpl: "a-aa:@bbb:{ccc}", data: map[string]interface{}{"bbb": "xyz", "ccc": 300}, want: "a-aa:xyz:300"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Translate(tt.tpl, tt.data); got != tt.want {
				t.Errorf("Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslateMap(t *testing.T) {

	tests := []struct {
		name string
		tpl  string
		data map[string]string
		want string
	}{
		{name: "2.str", tpl: "aaa:bbb:@ccc", data: map[string]string{"ccc": "100"}, want: "aaa:bbb:100"},
		{name: "2.str", tpl: "aaa:bbb:@{ccc}", data: map[string]string{"ccc": "200"}, want: "aaa:bbb:200"},
		{name: "2.str", tpl: "aaa:@bbb:@{ccc}", data: map[string]string{"bbb": "abcd", "ccc": "200"}, want: "aaa:abcd:200"},
		{name: "2.str", tpl: "@aaa:@bbb:@{ccc}", data: map[string]string{"bbb": "abcd", "ccc": "300"}, want: ":abcd:300"},
		// 新增测试用例：支持{variable}格式
		{name: "simple brace", tpl: "aaa:bbb:{ccc}", data: map[string]string{"ccc": "100"}, want: "aaa:bbb:100"},
		{name: "mixed formats", tpl: "aaa:{bbb}:@{ccc}", data: map[string]string{"bbb": "xyz", "ccc": "200"}, want: "aaa:xyz:200"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateMap(tt.tpl, tt.data); got != tt.want {
				t.Errorf("TranslateMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslateCallback(t *testing.T) {
	tests := []struct {
		name     string
		tpl      string
		callback func(string) string
		want     string
	}{
		{
			name:     "uppercase conversion",
			tpl:      "aaa:@bbb:@{ccc}",
			callback: strings.ToUpper,
			want:     "aaa:BBB:CCC",
		},
		{
			name:     "prefix addition",
			tpl:      "hello:@{world}",
			callback: func(s string) string { return "my" + s },
			want:     "hello:myworld",
		},
		{
			name:     "empty callback",
			tpl:      "@aaa:@bbb:@{ccc}",
			callback: func(s string) string { return "" },
			want:     "::",
		},
		{
			name:     "no vars",
			tpl:      "hello:world",
			callback: strings.ToUpper,
			want:     "hello:world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateCallback(tt.tpl, tt.callback); got != tt.want {
				t.Errorf("TranslateCallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslateObject(t *testing.T) {

	tests := []struct {
		name string
		tpl  string
		data interface{}
		want string
	}{
		{name: "1.str", tpl: "aaa:bbb:@ccc", data: []byte("100"), want: ""},
		{name: "2.str", tpl: "aaa:bbb:@ccc", data: map[string]string{"ccc": "100"}, want: "aaa:bbb:100"},
		{name: "3.str", tpl: "aaa:bbb:@ccc", data: map[string]interface{}{"ccc": "100"}, want: "aaa:bbb:100"},
		{name: "4.str", tpl: "aaa:bbb:@ccc", data: struct {
			CCC string `json:"ccc"`
		}{
			CCC: "100",
		}, want: "aaa:bbb:100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenericTranslate(tt.tpl, tt.data); got != tt.want {
				t.Errorf("TranslateObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
