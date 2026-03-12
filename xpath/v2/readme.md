# xpath v2

基于树匹配算法的路径匹配器，支持 `*` 和 `**` 的匹配规则。

## 特性

- 使用递归匹配算法进行路径匹配
- 支持 `*` 匹配单个路径段
- 支持 `**` 匹配零个或多个路径段
- 支持自定义分隔符
- 支持通配符模式匹配（如 `aaa*`）
- 支持缓存机制以提高性能

## 用法

```go
package main

import (
    "fmt"
    "github.com/zhiyunliu/golibs/xpath/v2"
)

func main() {
    // 创建匹配器实例
    matcher := v2.NewMatcher([]string{
        "/api/users/*",
        "/api/products/**",
        "/static/**",
        "/*.js",
    })
    
    // 执行匹配
    match, pattern := matcher.Match("/api/users/123")
    fmt.Printf("Match: %v, Pattern: %s\n", match, pattern)
    
    // 使用自定义分隔符
    matcherWithDot := v2.NewMatcher([]string{
        ".**.a.b.js", 
        ".config.*",
    }, v2.WithDelimiter("."))
    
    match, pattern = matcherWithDot.Match(".aa.bb.a.b.js", ".")
    fmt.Printf("Match: %v, Pattern: %s\n", match, pattern)
}
```

## API

### NewMatcher(pathList []string, opts ...Option) *Matcher

创建一个新的匹配器实例，接受路径列表和选项参数。

### NewMatcherPatterns(pathList []Pattern, opts ...Option) *Matcher

使用Pattern接口切片创建匹配器实例。

### Match(path string, spls ...string) (match bool, pattern string)

执行路径匹配，返回是否匹配成功和匹配的模式。

## 选项

### WithCache(enable bool)

启用或禁用缓存功能。

### WithDelimiter(delimiter string)

设置路径分隔符，默认为"/"。

## 匹配规则

- `*` 匹配单个路径段（例如 `/api/*` 匹配 `/api/users` 但不匹配 `/api/users/123`）
- `**` 匹配零个或多个路径段（例如 `/api/**` 匹配 `/api/users`、`/api/users/123`、`/api/users/123/posts` 等）
- 文本匹配支持通配符（例如 `/api/aaa*` 匹配 `/api/aaabbb`）
- 精确匹配（例如 `/api/users` 只匹配 `/api/users`）