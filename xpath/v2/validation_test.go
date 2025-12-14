package v2

import (
	"testing"
)

func TestValidatePattern(t *testing.T) {
	matcher := NewMatcher(nil)

	// 测试有效的模式
	validPatterns := []string{
		"/api/users",
		"/api/users/{id}",
		"/api/users/{id}/profile",
		"/assets/*",
		"/docs/**",
	}

	for _, pattern := range validPatterns {
		err := matcher.ValidatePattern(pattern)
		if err != nil {
			t.Errorf("Expected pattern %q to be valid, but got error: %v", pattern, err)
		}
	}

	// 测试无效的模式
	invalidPatterns := []string{
		"/api/users/{}",
		"/api/users/{id",
		"/api/users/id}",
		"/api/users/{{id}}",
		"",
		"api/users", // 不以分隔符开头
	}

	for _, pattern := range invalidPatterns {
		err := matcher.ValidatePattern(pattern)
		if err == nil {
			t.Errorf("Expected pattern %q to be invalid, but got no error", pattern)
		}
	}
}

func TestValidatePatternWithCustomDelimiter(t *testing.T) {
	matcher := NewMatcher(nil, WithDelimiter("."))

	// 测试有效的模式
	validPatterns := []string{
		".config",
		".config.db",
		".config.db.{host}",
		".assets.*",
		".docs.**",
	}

	for _, pattern := range validPatterns {
		err := matcher.ValidatePattern(pattern)
		if err != nil {
			t.Errorf("Expected pattern %q to be valid, but got error: %v", pattern, err)
		}
	}

	// 测试无效的模式
	invalidPatterns := []string{
		".config.db.{}",
		".config.db.{host",
		".config.db.host}",
		".config.db.{{host}}",
		"",
		"config.db", // 不以分隔符开头
	}

	for _, pattern := range invalidPatterns {
		err := matcher.ValidatePattern(pattern)
		if err == nil {
			t.Errorf("Expected pattern %q to be invalid, but got no error", pattern)
		}
	}
}

func TestAddPathWithValidation(t *testing.T) {
	matcher := NewMatcher(nil)

	validPaths := []string{
		"/api/users",
		"/api/users/{id}",
		"/api/users/{id}/profile",
		"/assets/*",
		"/docs/**",
	}

	for _, path := range validPaths {
		err := matcher.AddPathWithValidation(path)
		if err != nil {
			t.Errorf("Expected to add valid path %q successfully, but got error: %v", path, err)
		}
	}

	invalidPaths := []string{
		"/api/users/{}",
		"/api/users/{id",
		"/api/users/id}",
		"/api/users/{{id}}",
		"",
		"api/users",
		"/api/users/{}/profile",
	}

	for _, path := range invalidPaths {
		err := matcher.AddPathWithValidation(path)
		if err == nil {
			t.Errorf("Expected error when adding invalid path %q, but got no error", path)
		}
	}
}