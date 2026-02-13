package xtransform_test

import (
	"testing"

	"github.com/zhiyunliu/golibs/xtransform"
)

type StringerImpl string
type Inttpl int

func (s StringerImpl) String() string {
	return string(s)
}

func TestParseTemplate(t *testing.T) {
	tests := []struct {
		template string
		data     map[string]interface{}
		expected string
	}{
		{"aaa:bbb:@ccc", map[string]interface{}{"ccc": 100}, "aaa:bbb:100"},
		{"aaa:bbb:@{ccc}", map[string]interface{}{"ccc": 200}, "aaa:bbb:200"},
		{"aaa:@bbb:@{ccc}", map[string]interface{}{"bbb": "abcd", "ccc": 200}, "aaa:abcd:200"},
		{"@aaa:@bbb:@{ccc}", map[string]interface{}{"bbb": "abcd", "ccc": 200}, ":abcd:200"},
		{"@aaa:@Bbb:@{ccC}", map[string]interface{}{"Bbb": "abcd", "ccC": 200}, ":abcd:200"},
		{"@missing", map[string]interface{}{}, ""},
		{"@missing:@{alsoMissing}", map[string]interface{}{}, ":"},
		{"@{var}:@{var}", map[string]interface{}{"var": "value"}, "value:value"},
		{"@var:@{var}", map[string]interface{}{"var": "value"}, "value:value"},
		{"", map[string]interface{}{}, ""},
		{"noVarsHere", map[string]interface{}{}, "noVarsHere"},
		{"@{var}:@missing", map[string]interface{}{"var": "value"}, "value:"},
		{"@missing:@var", map[string]interface{}{"var": "value"}, ":value"},
		{"faaa:bbb:@ccc", map[string]interface{}{"ccc": 100.0}, "faaa:bbb:100"},
		{"uaaa:bbb:@ccc", map[string]interface{}{"ccc": uint(100)}, "uaaa:bbb:100"},
		{"baaa:bbb:@ccc", map[string]interface{}{"ccc": true}, "baaa:bbb:true"},
		{"taaa:bbb:@ccc", map[string]interface{}{"ccc": StringerImpl("tpl")}, "taaa:bbb:tpl"},
		{"saaa:bbb:@ccc", map[string]interface{}{"ccc": Inttpl(1011)}, "saaa:bbb:1011"},
	}

	for _, test := range tests {
		actual := xtransform.Translate(test.template, test.data)
		if actual != test.expected {
			t.Errorf("ParseTemplate(%q, %v) = %q; want %q", test.template, test.data, actual, test.expected)
		}
	}
}

func TestParseTemplateWithOptions(t *testing.T) {
	tests := []struct {
		template string
		data     map[string]interface{}
		opts     []xtransform.TranslateOption
		expected string
	}{
		{"aaa:bbb:@ccc", map[string]interface{}{"ccc": 100}, []xtransform.TranslateOption{xtransform.WithAtMode()}, "aaa:bbb:100"},
		{"aaa:bbb:@{ccc}", map[string]interface{}{"ccc": 200}, nil, "aaa:bbb:200"}, // 默认启用所有模式
		{"aaa:bbb:@{ccc}", map[string]interface{}{"ccc": 200}, []xtransform.TranslateOption{xtransform.WithAtBraceMode()}, "aaa:bbb:200"},
		{"aaa:bbb:{ccc}", map[string]interface{}{"ccc": 100}, []xtransform.TranslateOption{xtransform.WithBraceMode()}, "aaa:bbb:100"},
		{"@{var}:@{var}", map[string]interface{}{"var": "value"}, nil, "value:value"}, // 默认启用所有模式
		{"@var:@{var}", map[string]interface{}{"var": "value"}, []xtransform.TranslateOption{xtransform.WithAtMode(), xtransform.WithAtBraceMode()}, "value:value"},
		{"", map[string]interface{}{}, nil, ""},
		{"noVarsHere", map[string]interface{}{}, nil, "noVarsHere"},
		{"@{var}:@missing", map[string]interface{}{"var": "value"}, nil, "value:"}, // 默认启用所有模式
		{"@missing:@var", map[string]interface{}{"var": "value"}, []xtransform.TranslateOption{xtransform.WithAtMode()}, ":value"},
		{"@missing:@{var}", map[string]interface{}{"var": "value"}, nil, ":value"}, // 默认启用所有模式
		{"faaa:bbb:@ccc", map[string]interface{}{"ccc": 100.0}, []xtransform.TranslateOption{xtransform.WithAtMode()}, "faaa:bbb:100"},
		{"uaaa:bbb:@ccc", map[string]interface{}{"ccc": uint(100)}, []xtransform.TranslateOption{xtransform.WithAtMode()}, "uaaa:bbb:100"},
		{"taaa:bbb:@ccc", map[string]interface{}{"ccc": StringerImpl("tpl")}, []xtransform.TranslateOption{xtransform.WithAtMode()}, "taaa:bbb:tpl"},
		// 测试当只启用特定模式时，其他模式不应生效
		{"@ccc", map[string]interface{}{"ccc": 100}, []xtransform.TranslateOption{xtransform.WithAtBraceMode()}, "@ccc"}, // 只启用@{var}，@var不应匹配
		{"{ccc}", map[string]interface{}{"ccc": 100}, []xtransform.TranslateOption{xtransform.WithAtBraceMode()}, "{ccc}"}, // 只启用@{var}，{var}不应匹配
		{"@{ccc}", map[string]interface{}{"ccc": 100}, []xtransform.TranslateOption{xtransform.WithAtMode()}, "@{ccc}"}, // 只启用@var，@{var}不应匹配
	}

	for _, test := range tests {
		actual := xtransform.Translate(test.template, test.data, test.opts...)
		if actual != test.expected {
			t.Errorf("ParseTemplateWithOptions(%q, %v, %v) = %q; want %q", test.template, test.data, test.opts, actual, test.expected)
		}
	}
}