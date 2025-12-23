package benchmark

import (
	"testing"
	"time"

	xpath "github.com/zhiyunliu/golibs/xpath"
	xpathv2 "github.com/zhiyunliu/golibs/xpath/v2"
)

// BenchmarkV1SimplePath xpath v1 简单路径匹配
func BenchmarkV1SimplePath(b *testing.B) {
	m := xpath.NewMatch([]string{"/api/users", "/api/products", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/users")
	}
}

// BenchmarkV2SimplePath xpath v2 简单路径匹配
func BenchmarkV2SimplePath(b *testing.B) {
	m := xpathv2.NewMatcher([]string{"/api/users", "/api/products", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/users")
	}
}

// BenchmarkV1WildcardPath xpath v1 通配符路径匹配
func BenchmarkV1WildcardPath(b *testing.B) {
	m := xpath.NewMatch([]string{"/api/users/*", "/api/products/**", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/users/123")
	}
}

// BenchmarkV2WildcardPath xpath v2 通配符路径匹配
func BenchmarkV2WildcardPath(b *testing.B) {
	m := xpathv2.NewMatcher([]string{"/api/users/*", "/api/products/**", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/users/123")
	}
}

// BenchmarkV1DeepWildcardPath xpath v1 深层通配符路径匹配
func BenchmarkV1DeepWildcardPath(b *testing.B) {
	m := xpath.NewMatch([]string{"/api/users/*", "/api/products/**", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/products/category/electronics/item/123")
	}
}

// BenchmarkV2DeepWildcardPath xpath v2 深层通配符路径匹配
func BenchmarkV2DeepWildcardPath(b *testing.B) {
	m := xpathv2.NewMatcher([]string{"/api/users/*", "/api/products/**", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/products/category/electronics/item/123")
	}
}

// BenchmarkV1MultipleWildcards xpath v1 多通配符匹配
func BenchmarkV1MultipleWildcards(b *testing.B) {
	m := xpath.NewMatch([]string{"/api/**/users/*", "/api/**/products/**", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/v1/users/profile")
	}
}

// BenchmarkV2MultipleWildcards xpath v2 多通配符匹配
func BenchmarkV2MultipleWildcards(b *testing.B) {
	m := xpathv2.NewMatcher([]string{"/api/**/users/*", "/api/**/products/**", "/static/**"})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("/api/v1/users/profile")
	}
}

// PerformanceComparison runs performance comparison
func PerformanceComparison() {
	// 简单路径匹配对比
	patterns := []string{"/api/users", "/api/products", "/static/**"}
	paths := []string{"/api/users", "/api/products", "/static/css/style.css"}

	// 创建 v1 和 v2 的匹配器
	v1Matcher := xpath.NewMatch(patterns)
	v2Matcher := xpathv2.NewMatcher(patterns)

	// 测试 V1 性能
	start := time.Now()
	for i := 0; i < 100000; i++ {
		for _, path := range paths {
			v1Matcher.Match(path)
		}
	}
	v1Duration := time.Since(start)

	// 测试 V2 性能
	start = time.Now()
	for i := 0; i < 100000; i++ {
		for _, path := range paths {
			v2Matcher.Match(path)
		}
	}
	v2Duration := time.Since(start)

	println("Performance Comparison (100000 iterations):")
	println("V1 (regexp):", v1Duration.String())
	println("V2 (tree-based):", v2Duration.String())
	println("V2 is", float64(v1Duration)/float64(v2Duration), "times faster than V1")
}