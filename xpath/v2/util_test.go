package v2

import (
	"reflect"
	"testing"
)

func TestParseSegment(t *testing.T) {
	tests := []struct {
		segment     string
		expectedType NodeType
		expectedName string
	}{
		{"users", StaticNode, ""},
		{"{id}", ParamNode, "id"},
		{"{*}", ParamNode, "*"},
		{"*", WildcardNode, ""},
		{"**", CatchAllNode, ""},
		{"{user_id}", ParamNode, "user_id"},
		{"api", StaticNode, ""},
	}

	for _, test := range tests {
		nodeType, paramName := parseSegment(test.segment)
		if nodeType != test.expectedType {
			t.Errorf("For segment '%s': expected node type %v, got %v", test.segment, test.expectedType, nodeType)
		}
		if paramName != test.expectedName {
			t.Errorf("For segment '%s': expected param name '%s', got '%s'", test.segment, test.expectedName, paramName)
		}
	}
}

func TestSplitPath(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		path     string
		expected []string
	}{
		{"/", []string{}},
		{"/api", []string{"api"}},
		{"/api/users", []string{"api", "users"}},
		{"/api/users/", []string{"api", "users"}},
		{"/api/users/123", []string{"api", "users", "123"}},
		{"", []string{}},
	}

	for _, test := range tests {
		result := matcher.splitPath(test.path)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For path '%s': expected %v, got %v", test.path, test.expected, result)
		}
	}
}

func TestSplitPathWithCustomDelimiter(t *testing.T) {
	matcher := NewMatcher(WithDelimiter("."))

	tests := []struct {
		path     string
		expected []string
	}{
		{".", []string{}},
		{".config", []string{"config"}},
		{".config.database", []string{"config", "database"}},
		{".config.database.", []string{"config", "database"}},
		{".config.database.host", []string{"config", "database", "host"}},
		{"", []string{}},
	}

	for _, test := range tests {
		result := matcher.splitPath(test.path)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For path '%s': expected %v, got %v", test.path, test.expected, result)
		}
	}
}

func TestJoinPath(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		segments []string
		expected string
	}{
		{[]string{}, "/"},
		{[]string{""}, "/"},
		{[]string{"api"}, "/api"},
		{[]string{"api", "users"}, "/api/users"},
		{[]string{"", "api", "users"}, "/api/users"},
		{[]string{"api", "", "users"}, "/api/users"},
	}

	for _, test := range tests {
		result := matcher.joinPath(test.segments...)
		if result != test.expected {
			t.Errorf("For segments %v: expected '%s', got '%s'", test.segments, test.expected, result)
		}
	}
}

func TestJoinPathWithCustomDelimiter(t *testing.T) {
	matcher := NewMatcher(WithDelimiter("."))

	tests := []struct {
		segments []string
		expected string
	}{
		{[]string{}, "."},
		{[]string{""}, "."},
		{[]string{"config"}, ".config"},
		{[]string{"config", "database"}, ".config.database"},
		{[]string{"", "config", "database"}, ".config.database"},
		{[]string{"config", "", "database"}, ".config.database"},
	}

	for _, test := range tests {
		result := matcher.joinPath(test.segments...)
		if result != test.expected {
			t.Errorf("For segments %v: expected '%s', got '%s'", test.segments, test.expected, result)
		}
	}
}