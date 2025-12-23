package main

import (
	"fmt"
	"time"

	xpath "github.com/zhiyunliu/golibs/xpath"
	xpathv2 "github.com/zhiyunliu/golibs/xpath/v2"
)

func main() {
	PerformanceComparison()
}

// PerformanceComparison runs performance comparison
func PerformanceComparison() {
	// 简单路径匹配对比
	patterns := []string{"/api/users", "/api/products", "/static/**"}
	paths := []string{"/api/users", "/api/products", "/static/css/style.css"}

	// 创建 v1 和 v2 的匹配器
	v1Matcher := xpath.NewMatch(patterns)
	v2Matcher := xpathv2.NewMatch(patterns)

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

	fmt.Printf("Performance Comparison (100000 iterations):\n")
	fmt.Printf("V1 (regexp): %v\n", v1Duration)
	fmt.Printf("V2 (tree-based): %v\n", v2Duration)
	fmt.Printf("V2 is %.2f times faster than V1\n", float64(v1Duration)/float64(v2Duration))
}
