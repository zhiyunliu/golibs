package v2

import (
	"testing"
)

func TestValidatePattern(t *testing.T) {
	matcher := NewMatcher()

	// 测试有效的模式
	validPatterns := []string{
		"/api/users",
		"/api/users/{id}",
		"/api/assets/*",
		"/api/docs/**",
		"/",
	}

	for _, pattern := range validPatterns {
		err := matcher.ValidatePattern(pattern)
		if err != nil {
			t.Errorf("Expected pattern '%s' to be valid, but got error: %v", pattern, err)
		}
	}

	// 测试无效的模式
	invalidPatterns := []struct {
		pattern string
		reason  string
	}{
		{"", "empty pattern"},
		{"api/users", "must start with delimiter"},
		{"/api//users", "empty segment"},
		{"/api/users/**/docs", "catch-all must be last"},
		{"/api/users/{invalid-param}", "invalid parameter name"},
		{"/api/users/{}", "empty parameter name"},
	}

	for _, test := range invalidPatterns {
		err := matcher.ValidatePattern(test.pattern)
		if err == nil {
			t.Errorf("Expected pattern '%s' to be invalid (%s), but got no error", test.pattern, test.reason)
		}
	}
}

func TestValidatePatternWithCustomDelimiter(t *testing.T) {
	matcher := NewMatcher(WithDelimiter("."))

	// 测试有效的模式
	validPatterns := []string{
		".config.database",
		".config.database.{table}",
		".assets.*",
		".logs.**",
		".",
	}

	for _, pattern := range validPatterns {
		err := matcher.ValidatePattern(pattern)
		if err != nil {
			t.Errorf("Expected pattern '%s' to be valid, but got error: %v", pattern, err)
		}
	}

	// 测试无效的模式
	invalidPatterns := []struct {
		pattern string
		reason  string
	}{
		{"", "empty pattern"},
		{"config.database", "must start with delimiter"},
		{".config..database", "empty segment"},
		{".config.logs.**.archive", "catch-all must be last"},
		{".config.table.{invalid-table}", "invalid parameter name"},
		{".config.table.{}", "empty parameter name"},
	}

	for _, test := range invalidPatterns {
		err := matcher.ValidatePattern(test.pattern)
		if err == nil {
			t.Errorf("Expected pattern '%s' to be invalid (%s), but got no error", test.pattern, test.reason)
		}
	}
}

func TestAddPathWithValidation(t *testing.T) {
	matcher := NewMatcher()

	// 测试有效的路径添加
	validPaths := []string{
		"/api/users",
		"/api/users/{id}",
		"/api/assets/*",
		"/api/docs/**",
	}

	for _, path := range validPaths {
		err := matcher.AddPathWithValidation(path)
		if err != nil {
			t.Errorf("Expected to add path '%s' successfully, but got error: %v", path, err)
		}
	}

	// 测试无效的路径添加
	invalidPaths := []string{
		"",
		"api/users",
		"/api//users",
		"/api/users/**/docs",
		"/api/users/{invalid-param}",
		"/api/users/{}",
	}

	for _, path := range invalidPaths {
		err := matcher.AddPathWithValidation(path)
		if err == nil {
			t.Errorf("Expected error when adding invalid path '%s', but got no error", path)
		}
	}
}