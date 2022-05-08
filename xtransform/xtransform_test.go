package xtransform

import (
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateMap(tt.tpl, tt.data); got != tt.want {
				t.Errorf("TranslateMap() = %v, want %v", got, tt.want)
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
		{name: "2.str", tpl: "aaa:bbb:@ccc", data: map[string]string{"ccc": "100"}, want: "aaa:bbb:100"},
		{name: "2.str", tpl: "aaa:bbb:@ccc", data: map[string]interface{}{"ccc": "100"}, want: "aaa:bbb:100"},
		{name: "2.str", tpl: "aaa:bbb:@ccc", data: struct {
			CCC string `json:"ccc"`
		}{
			CCC: "100",
		}, want: "aaa:bbb:100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateObject(tt.tpl, tt.data); got != tt.want {
				t.Errorf("TranslateObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
