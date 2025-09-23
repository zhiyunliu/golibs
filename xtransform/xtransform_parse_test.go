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
