package benchmark

import (
	"runtime"
	"testing"

	xpath "github.com/zhiyunliu/golibs/xpath"
	xpathv2 "github.com/zhiyunliu/golibs/xpath/v2"
)

// MemStats 用于记录内存使用情况
type MemStats struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

// getMemStats 获取当前内存统计信息
func getMemStats() MemStats {
	var m runtime.MemStats
	runtime.GC() // 先执行垃圾回收
	runtime.ReadMemStats(&m)
	return MemStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}
}

// BenchmarkV1MemoryUsage 测试xpath v1内存使用情况
func BenchmarkV1MemoryUsage(b *testing.B) {
	var m1, m2 MemStats

	m1 = getMemStats()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建Matcher实例
		matcher := xpath.NewMatch([]string{
			"/api/users/*",
			"/api/products/**",
			"/static/**",
			"/api/**/profile",
			"/admin/*",
			"/public/**",
			"/files/*/upload",
			"/data/**/export",
		})

		// 执行一些匹配操作
		_, _ = matcher.Match("/api/users/123")
		_, _ = matcher.Match("/api/products/electronics/phone")
		_, _ = matcher.Match("/static/css/style.css")
		_, _ = matcher.Match("/api/v1/users/profile")
	}

	b.StopTimer()
	m2 = getMemStats()

	b.ReportMetric(float64(m2.Alloc-m1.Alloc)/float64(b.N), "B/op")
	b.ReportMetric(float64(m2.NumGC-m1.NumGC)/float64(b.N), "GC/op")
}

// BenchmarkV2MemoryUsage 测试xpath v2内存使用情况
func BenchmarkV2MemoryUsage(b *testing.B) {
	var m1, m2 MemStats

	m1 = getMemStats()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建Matcher实例
		matcher := xpathv2.NewMatch([]string{
			"/api/users/*",
			"/api/products/**",
			"/static/**",
			"/api/**/profile",
			"/admin/*",
			"/public/**",
			"/files/*/upload",
			"/data/**/export",
		})

		// 执行一些匹配操作
		_, _ = matcher.Match("/api/users/123")
		_, _ = matcher.Match("/api/products/electronics/phone")
		_, _ = matcher.Match("/static/css/style.css")
		_, _ = matcher.Match("/api/v1/users/profile")
	}

	b.StopTimer()
	m2 = getMemStats()

	b.ReportMetric(float64(m2.Alloc-m1.Alloc)/float64(b.N), "B/op")
	b.ReportMetric(float64(m2.NumGC-m1.NumGC)/float64(b.N), "GC/op")
}

// BenchmarkV1MemoryUsageWithCache 测试xpath v1带缓存的内存使用情况
func BenchmarkV1MemoryUsageWithCache(b *testing.B) {
	// xpath v1 doesn't have cache option
	b.Skip("xpath v1 doesn't support cache")
}

// BenchmarkV2MemoryUsageWithCache 测试xpath v2带缓存的内存使用情况
func BenchmarkV2MemoryUsageWithCache(b *testing.B) {
	var m1, m2 MemStats

	m1 = getMemStats()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建Matcher实例，启用缓存
		matcher := xpathv2.NewMatch([]string{
			"/api/users/*",
			"/api/products/**",
			"/static/**",
			"/api/**/profile",
			"/admin/*",
			"/public/**",
			"/files/*/upload",
			"/data/**/export",
		}, xpathv2.WithCache[xpathv2.Pattern](true))

		// 执行一些匹配操作
		_, _ = matcher.Match("/api/users/123")
		_, _ = matcher.Match("/api/products/electronics/phone")
		_, _ = matcher.Match("/static/css/style.css")
		_, _ = matcher.Match("/api/v1/users/profile")
	}

	b.StopTimer()
	m2 = getMemStats()

	b.ReportMetric(float64(m2.Alloc-m1.Alloc)/float64(b.N), "B/op")
	b.ReportMetric(float64(m2.NumGC-m1.NumGC)/float64(b.N), "GC/op")
}
