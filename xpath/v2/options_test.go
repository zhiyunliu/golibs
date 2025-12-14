package v2

import (
	"reflect"
	"testing"
)

func TestWithName(t *testing.T) {
	option := WithName("测试名称")
	info := &NodeInfo{}

	option(info)

	if info.Name != "测试名称" {
		t.Errorf("Expected name '测试名称', got '%s'", info.Name)
	}
}

func TestWithDesc(t *testing.T) {
	option := WithDesc("测试描述")
	info := &NodeInfo{}

	option(info)

	if info.Desc != "测试描述" {
		t.Errorf("Expected description '测试描述', got '%s'", info.Desc)
	}
}

func TestWithMeta(t *testing.T) {
	initialMeta := map[string]any{
		"key1": "value1",
		"key2": 42,
	}

	option := WithMeta(initialMeta)
	info := &NodeInfo{}

	option(info)

	if info.Meta == nil {
		t.Error("Expected meta to be initialized, but it's nil")
	}

	if !reflect.DeepEqual(info.Meta, initialMeta) {
		t.Errorf("Expected meta %v, got %v", initialMeta, info.Meta)
	}
}

func TestWithMetaMerge(t *testing.T) {
	info := &NodeInfo{
		Meta: map[string]any{
			"existing": "value",
		},
	}

	additionalMeta := map[string]any{
		"new": "value",
	}

	option := WithMeta(additionalMeta)
	option(info)

	expected := map[string]any{
		"existing": "value",
		"new":      "value",
	}

	if !reflect.DeepEqual(info.Meta, expected) {
		t.Errorf("Expected merged meta %v, got %v", expected, info.Meta)
	}
}

func TestWithDelimiter(t *testing.T) {
	option := WithDelimiter(".")
	matcher := &Matcher{
		delimiter: "/",
	}

	option(matcher)

	if matcher.delimiter != "." {
		t.Errorf("Expected delimiter '.', got '%s'", matcher.delimiter)
	}
}