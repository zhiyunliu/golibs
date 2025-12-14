# xpath/v2

路径匹配器，支持参数匹配、通配符匹配和多段匹配。

## 使用示例

```go
import "github.com/zhiyunliu/golibs/xpath/v2"

// 创建匹配器
matcher := NewMatcher(nil)

// 添加路径规则
matcher.MustAddPath("/api/users", WithName("用户列表"))
matcher.MustAddPath("/api/users/{id}", WithName("用户详情"))
matcher.MustAddPath("/api/assets/*", WithName("资源文件"))
matcher.MustAddPath("/api/docs/**", WithName("文档资源"))

// 匹配路径
result := matcher.Match("/api/users/123")
if result.Matched {
    fmt.Println("匹配成功:", result.Info.Name) // 用户详情
    fmt.Println("参数:", result.Params["id"]) // 123
}
```

## 自定义分隔符

```go
// 使用自定义分隔符
matcher := NewMatcher(nil, WithDelimiter("."))
matcher.MustAddPath(".config.database.host", WithName("数据库主机"))
result := matcher.Match(".config.database.host")
```