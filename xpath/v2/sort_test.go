package v2

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 实现一个简单的Pattern类型用于测试
type testPattern struct {
	pattern string
}

func (tp testPattern) Pattern() string {
	return tp.pattern
}

func TestSortPatterns_Len(t *testing.T) {
	patterns := SortPatterns[testPattern]{
		{pattern: "/api/users"},
		{pattern: "/api/products"},
		{pattern: "/api/orders"},
	}

	expectedLen := 3
	actualLen := patterns.Len()

	assert.Equal(t, expectedLen, actualLen, "Expected length %d, got %d", expectedLen, actualLen)
}

func TestSortPatterns_Swap(t *testing.T) {
	patterns := SortPatterns[testPattern]{
		{pattern: "/api/users"},
		{pattern: "/api/products"},
	}

	originalFirst := patterns[0].Pattern()
	originalSecond := patterns[1].Pattern()

	// 交换位置
	patterns.Swap(0, 1)

	// 验证交换后的位置
	assert.Equal(t, originalSecond, patterns[0].Pattern(), "First element should be swapped")
	assert.Equal(t, originalFirst, patterns[1].Pattern(), "Second element should be swapped")
}

func TestSortPatterns_SwapOutOfBounds(t *testing.T) {
	// 测试长度为2的数组
	patterns := SortPatterns[testPattern]{
		{pattern: "/api/users"},
		{pattern: "/api/products"},
	}

	originalPatterns := make([]testPattern, len(patterns))
	copy(originalPatterns, patterns)

	// 保存原始值用于比较
	orig0 := patterns[0]
	orig1 := patterns[1]

	// 尝试交换超出边界的索引，应该不会改变数组
	patterns.Swap(0, 5) // 第二个索引超出范围
	// 在这种情况下，Swap会检查到j>=len，直接返回，数组不变
	assert.Equal(t, orig0, patterns[0])
	assert.Equal(t, orig1, patterns[1])

	// 重置
	patterns[0] = orig0
	patterns[1] = orig1

	patterns.Swap(5, 0) // 第一个索引超出范围
	// 在这种情况下，Swap会检查到i>=len，直接返回，数组不变
	assert.Equal(t, orig0, patterns[0])
	assert.Equal(t, orig1, patterns[1])
}

func TestSortPatterns_Less(t *testing.T) {
	testCases := []struct {
		name     string
		first    string
		second   string
		expected bool
	}{
		{"Simple comparison", "/api/a", "/api/b", true},
		{"Reverse comparison", "/api/b", "/api/a", false},
		{"Wildcard vs normal", "/api/*", "/api/a", false},      // * comes after normal chars
		{"Normal vs wildcard", "/api/a", "/api/*", true},       // * comes after normal chars
		{"Double wildcard", "/api/**", "/api/*", false},        // ** comes after *
		{"Double wildcard reverse", "/api/*", "/api/**", true}, // ** comes after *
		{"Same prefix", "/api/v1", "/api/v2", true},            // api/v1 < api/v2
		{"Same prefix reverse", "/api/v2", "/api/v1", false},   // api/v2 !< api/v1
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patterns := SortPatterns[testPattern]{
				{pattern: tc.first},
				{pattern: tc.second},
			}

			result := patterns.Less(0, 1)
			assert.Equal(t, tc.expected, result, "Comparison failed for '%s' < '%s'", tc.first, tc.second)
		})
	}
}

// TestSortPatterns_EdgeCases 测试边界情况
func TestSortPatterns_EdgeCases(t *testing.T) {
	// 测试完全相同的模式
	patterns := SortPatterns[testPattern]{
		{pattern: "/api/test"},
		{pattern: "/api/test"},
	}

	// 一个值不能小于自己
	assert.False(t, patterns.Less(0, 1), "Identical patterns should not be less than each other")
	assert.False(t, patterns.Less(1, 0), "Identical patterns should not be less than each other")

	// 测试空字符串模式
	patterns2 := SortPatterns[testPattern]{
		{pattern: ""},
		{pattern: "/api/test"},
	}

	// 空字符串应该小于任何非空字符串
	assert.True(t, patterns2.Less(0, 1), "Empty pattern should be less than non-empty pattern")
	assert.False(t, patterns2.Less(1, 0), "Non-empty pattern should not be less than empty pattern")
}

func TestSortPatterns_LessLongerVsShorter(t *testing.T) {
	// 根据实际排序逻辑，理解如何处理长短字符串
	patterns := SortPatterns[testPattern]{
		{pattern: "/api/longer/path"},
		{pattern: "/api/short"},
	}

	// 实际上，按字典序比较，"longer/path"[0]='l' < "short"[0]='s'
	// 所以 "/api/longer/path" < "/api/short" 是 true
	// 这意味着 "/api/longer/path" 应该排在 "/api/short" 之前
	result := patterns.Less(0, 1)
	assert.True(t, result, "Longer path (/api/longer/path) should be less than shorter path (/api/short)")

	// 而 "/api/short" < "/api/longer/path" 是false，所以Less(1,0) 应该返回false
	result = patterns.Less(1, 0)
	assert.False(t, result, "Shorter path (/api/short) should not be less than longer path (/api/longer/path)")
}

func TestSortPatterns_Sorting(t *testing.T) {
	// 测试整个排序过程
	unsortedPatterns := []testPattern{
		{pattern: "/api/**"},
		{pattern: "/api/users/*"},
		{pattern: "/api/users/profile"},
		{pattern: "/api/products"},
		{pattern: "/api/*"},
	}

	// 对排序前后的模式进行检查
	initialOrder := make([]string, len(unsortedPatterns))
	for i, p := range unsortedPatterns {
		initialOrder[i] = p.Pattern()
	}

	sort.Sort(SortPatterns[testPattern](unsortedPatterns))

	actualOrder := make([]string, len(unsortedPatterns))
	for i, p := range unsortedPatterns {
		actualOrder[i] = p.Pattern()
	}

	// 根据实际排序逻辑修正期望值
	// 根据Less函数逻辑，* 应该排在普通字符后面，** 排在 * 后面
	expectedOrder := []string{
		"/api/products",      // 没有通配符
		"/api/users/profile", // 没有通配符
		"/api/users/*",       // 有 * 通配符
		"/api/*",             // 有 * 通配符
		"/api/**",            // 有 ** 通配符
	}

	assert.Equal(t, expectedOrder, actualOrder, "Sorted order doesn't match expected")
}

func TestSortPatterns_LessWithSamePrefix(t *testing.T) {
	// 测试具有相同前缀的模式排序
	patterns := SortPatterns[testPattern]{
		{pattern: "/api/v1/users"},
		{pattern: "/api/v1/*"},
		{pattern: "/api/v1/**"},
	}

	// 验证排序顺序: /api/v1/** 应该排在最后，/api/v1/* 在中间，/api/v1/users 最前
	assert.True(t, patterns.Less(0, 1), "Specific path should come before wildcard")               // /api/v1/users < /api/v1/*
	assert.True(t, patterns.Less(1, 2), "Single wildcard should come before double wildcard")      // /api/v1/* < /api/v1/**
	assert.False(t, patterns.Less(2, 1), "Double wildcard should not come before single wildcard") // /api/v1/** !< /api/v1/*
	assert.False(t, patterns.Less(2, 0), "Double wildcard should not come before specific path")   // /api/v1/** !< /api/v1/users
}

func TestSortPatterns_EmptyAndSingleElement(t *testing.T) {
	// 测试空数组
	emptyPatterns := SortPatterns[testPattern]{}
	assert.Equal(t, 0, emptyPatterns.Len(), "Empty array should have length 0")

	// 测试单元素数组
	singlePattern := SortPatterns[testPattern]{{pattern: "/api/test"}}
	assert.Equal(t, 1, singlePattern.Len(), "Single element array should have length 1")

	// 对单元素数组进行交换，应该是无操作
	original := singlePattern[0]
	singlePattern.Swap(0, 0)
	assert.Equal(t, original, singlePattern[0], "Swapping element with itself should not change anything")
}

func TestSortPatterns_ComplexSorting(t *testing.T) {
	// 更复杂的排序测试
	unsortedPatterns := []testPattern{
		{pattern: "/app/**"},
		{pattern: "/app/user/profile"},
		{pattern: "/app/*"},
		{pattern: "/app/user/*"},
		{pattern: "/app/user/settings"},
		{pattern: "/app/data"},
	}

	sort.Sort(SortPatterns[testPattern](unsortedPatterns))

	actualOrder := make([]string, len(unsortedPatterns))
	for i, p := range unsortedPatterns {
		actualOrder[i] = p.Pattern()
	}

	// 预期顺序应该是按照字符串比较规则，特殊字符在字母之后
	// 具体路径排在通配符前面，单星号在双星号前面
	expectedOrder := []string{
		"/app/data",          // 具体路径
		"/app/user/profile",  // 具体路径
		"/app/user/settings", // 具体路径
		"/app/user/*",        // 单星号
		"/app/*",             // 单星号
		"/app/**",            // 双星号
	}

	assert.Equal(t, expectedOrder, actualOrder, "Complex sorted order doesn't match expected")
}

// 专门用于调试排序逻辑的测试
func TestDebugSortLogic(t *testing.T) {
	patterns := SortPatterns[testPattern]{
		{pattern: "/api/longer/path"},
		{pattern: "/api/short"},
	}

	// 检查实际比较结果
	result01 := patterns.Less(0, 1) // "/api/longer/path" < "/api/short"
	result10 := patterns.Less(1, 0) // "/api/short" < "/api/longer/path"

	t.Logf("/api/longer/path < /api/short = %t", result01)
	t.Logf("/api/short < /api/longer/path = %t", result10)

	// 根据逻辑，"l" < "s"，所以 "/api/longer/path" < "/api/short" 应该是 true
	// 因此 "/api/short" < "/api/longer/path" 应该是 false
	assert.True(t, result01, "/api/longer/path should be less than /api/short")
	assert.False(t, result10, "/api/short should not be less than /api/longer/path")
}

type sortKeyTestPattern struct {
	pattern string
}

func (tp sortKeyTestPattern) Pattern() string {
	return tp.pattern
}

func (tp sortKeyTestPattern) SortKey() string {
	return tp.pattern
}

// 专门用于调试排序逻辑的测试
func TestDebugSortLogicSortKey(t *testing.T) {
	patterns := SortPatterns[sortKeyTestPattern]{
		{pattern: "/api/longer/path"},
		{pattern: "/api/short"},
	}

	// 检查实际比较结果
	result01 := patterns.Less(0, 1) // "/api/longer/path" < "/api/short"
	result10 := patterns.Less(1, 0) // "/api/short" < "/api/longer/path"

	t.Logf("/api/longer/path < /api/short = %t", result01)
	t.Logf("/api/short < /api/longer/path = %t", result10)

	// 根据逻辑，"l" < "s"，所以 "/api/longer/path" < "/api/short" 应该是 true
	// 因此 "/api/short" < "/api/longer/path" 应该是 false
	assert.True(t, result01, "/api/longer/path should be less than /api/short")
	assert.False(t, result10, "/api/short should not be less than /api/longer/path")
}
