package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcherMatch_Comprehensive(t *testing.T) {
	tests := []struct {
		name          string
		patterns      []string
		path          string
		opts          []Option[Pattern]
		expectMatch   bool
		expectPattern string
		note          string // 用于解释特殊情况
	}{
		{
			name:          "exact match",
			patterns:      []string{"/api/users", "/api/products"},
			path:          "/api/users",
			expectMatch:   true,
			expectPattern: "/api/users",
		},
		{
			name:          "no match",
			patterns:      []string{"/api/users", "/api/products"},
			path:          "/api/orders",
			expectMatch:   false,
			expectPattern: "",
		},
		{
			name:          "single wildcard match",
			patterns:      []string{"/api/*", "/api/products"},
			path:          "/api/users",
			expectMatch:   true,
			expectPattern: "/api/*",
		},
		{
			name:          "double wildcard match zero segments",
			patterns:      []string{"/api/**/users", "/api/products"},
			path:          "/api/users",
			expectMatch:   true,
			expectPattern: "/api/**/users",
		},
		{
			name:          "double wildcard match multiple segments",
			patterns:      []string{"/api/**/profile", "/api/products"},
			path:          "/api/users/123/profile",
			expectMatch:   true,
			expectPattern: "/api/**/profile",
		},
		{
			name:          "literal pattern with wildcard",
			patterns:      []string{"/api/user*", "/api/products"},
			path:          "/api/users",
			expectMatch:   true,
			expectPattern: "/api/user*",
		},
		{
			name:          "literal pattern with multiple wildcards",
			patterns:      []string{"/api/u*rs", "/api/products"},
			path:          "/api/users",
			expectMatch:   true,
			expectPattern: "/api/u*rs",
		},
		{
			name:          "dot separator",
			patterns:      []string{"api.users", "api.products"},
			path:          "api.users",
			opts:          []Option[Pattern]{WithDelimiter[Pattern](".")},
			expectMatch:   true,
			expectPattern: "api.users",
		},
		{
			name:          "dot separator with wildcard",
			patterns:      []string{"api.*", "api.products"},
			path:          "api.users",
			opts:          []Option[Pattern]{WithDelimiter[Pattern](".")},
			expectMatch:   true,
			expectPattern: "api.*",
		},
		{
			name:          "cache enabled - first call",
			patterns:      []string{"/api/users", "/api/products"},
			path:          "/api/users",
			opts:          []Option[Pattern]{WithCache[Pattern](true)},
			expectMatch:   true,
			expectPattern: "/api/users",
		},
		{
			name:          "cache enabled - non matching",
			patterns:      []string{"/api/users", "/api/products"},
			path:          "/api/orders",
			opts:          []Option[Pattern]{WithCache[Pattern](true)},
			expectMatch:   false,
			expectPattern: "",
		},
		{
			name:          "empty patterns",
			patterns:      []string{},
			path:          "/api/users",
			expectMatch:   false,
			expectPattern: "",
		},
		{
			name:          "empty pattern matches empty path",
			patterns:      []string{"", "/api/users"},
			path:          "",
			expectMatch:   true,
			expectPattern: "",
		},
		{
			name:          "path and pattern both without leading delimiter",
			patterns:      []string{"api/users", "api/products"},
			path:          "api/users",
			expectMatch:   true,
			expectPattern: "api/users",
		},
		{
			name:          "path with leading delimiter matches pattern without",
			patterns:      []string{"api/users", "api/products"},
			path:          "/api/users",
			expectMatch:   true, // This works because both get processed the same way
			expectPattern: "api/users",
			note:          "Both path and pattern get processed by removing leading empty segments",
		},
		{
			name:          "pattern with leading delimiter matches path without",
			patterns:      []string{"/api/users", "/api/products"},
			path:          "api/users",
			expectMatch:   true, // This works because both get processed the same way
			expectPattern: "/api/users",
			note:          "Both path and pattern get processed by removing leading empty segments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewMatch(tt.patterns, tt.opts...)
			match, pattern := matcher.Match(tt.path)

			assert.Equal(t, tt.expectMatch, match, "Match result should match expected")
			if tt.expectPattern != "" {
				assert.Equal(t, tt.expectPattern, pattern.Pattern(), "Pattern should match expected")
			} else {
				if !tt.expectMatch {
					assert.Nil(t, pattern, "Pattern should be nil when no match")
				} else {
					// In case of match but expecting empty string pattern
					assert.Equal(t, tt.expectPattern, pattern.Pattern(), "Pattern string should match expected")
				}
			}
		})
	}
}

func TestMatcherMatch_CacheBehavior(t *testing.T) {
	patterns := []string{"/api/users", "/static/**"}
	matcher := NewMatch(patterns, WithCache[Pattern](true))

	// First call should not use cache
	match, pattern := matcher.Match("/api/users")
	assert.True(t, match)
	assert.Equal(t, "/api/users", pattern.Pattern())

	// Second call should potentially use cache (internal behavior)
	match, pattern = matcher.Match("/api/users")
	assert.True(t, match)
	assert.Equal(t, "/api/users", pattern.Pattern())

	// Test non-matching path
	match, pattern = matcher.Match("/api/orders")
	assert.False(t, match)
	assert.Nil(t, pattern)
}

func TestMatcherMatch_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		expected bool
	}{
		{
			name:     "path with multiple consecutive delimiters",
			patterns: []string{"/api/users"},
			path:     "//api//users", // This won't match as the path gets split differently
			expected: false,
		},
		{
			name:     "pattern with * at beginning",
			patterns: []string{"*users"},
			path:     "apiusers",
			expected: true,
		},
		{
			name:     "pattern with * at end",
			patterns: []string{"api*"},
			path:     "apiusers",
			expected: true,
		},
		{
			name:     "pattern with * in middle",
			patterns: []string{"api*data"},
			path:     "apidata",
			expected: true,
		},
		{
			name:     "pattern with * in middle - no match",
			patterns: []string{"api*data"},
			path:     "apixdatax",
			expected: false,
		},
		{
			name:     "only ** pattern",
			patterns: []string{"**"},
			path:     "/any/path/here",
			expected: true,
		},
		{
			name:     "only * pattern",
			patterns: []string{"*"},
			path:     "single",
			expected: true,
		},
		{
			name:     "only * pattern - no match",
			patterns: []string{"*"},
			path:     "path/with/slashes",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewMatch(tt.patterns)
			match, _ := matcher.Match(tt.path)
			assert.Equal(t, tt.expected, match, "Match result should match expected")
		})
	}
}

// TestMatcherMatch_CodePaths tests specific code paths in the Match function
func TestMatcherMatch_CodePaths(t *testing.T) {
	t.Run("cache path enabled", func(t *testing.T) {
		matcher := NewMatch([]string{"/api/users"}, WithCache[Pattern](true))

		// First call - not in cache
		match, pattern := matcher.Match("/api/users")
		assert.True(t, match)
		assert.Equal(t, "/api/users", pattern.Pattern())

		// Second call - potentially from cache
		match, pattern = matcher.Match("/api/users")
		assert.True(t, match)
		assert.Equal(t, "/api/users", pattern.Pattern())
	})

	t.Run("cache path disabled", func(t *testing.T) {
		matcher := NewMatch([]string{"/api/users"}) // no WithCache option = cache disabled

		match, pattern := matcher.Match("/api/users")
		assert.True(t, match)
		assert.Equal(t, "/api/users", pattern.Pattern())
	})

	t.Run("no match case", func(t *testing.T) {
		matcher := NewMatch([]string{"/api/users", "/api/orders"})

		match, pattern := matcher.Match("/api/products")
		assert.False(t, match)
		assert.Nil(t, pattern)
	})
}
