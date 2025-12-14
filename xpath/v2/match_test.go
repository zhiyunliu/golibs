package v2

import (
	"testing"
)

func TestBasicMatch(t *testing.T) {
	matcher := NewMatcher(nil)

	// 添加静态路径
	matcher.MustAddPath("/api/users", WithName("用户列表"), WithDesc("获取用户列表"))

	// 测试精确匹配
	result := matcher.Match("/api/users")
	if !result.Matched {
		t.Error("Expected match for /api/users, but got no match")
	}

	if result.Info.Name != "用户列表" {
		t.Errorf("Expected name '用户列表', got '%s'", result.Info.Name)
	}

	// 测试不匹配的情况
	result = matcher.Match("/api/products")
	if result.Matched {
		t.Error("Expected no match for /api/products, but got a match")
	}
}

func TestParamMatch(t *testing.T) {
	matcher := NewMatcher(nil)

	// 添加带参数的路径
	matcher.MustAddPath("/api/users/{id}", WithName("用户详情"))

	// 测试参数匹配
	result := matcher.Match("/api/users/123")
	if !result.Matched {
		t.Error("Expected match for /api/users/123, but got no match")
	}

	if result.Info.Name != "用户详情" {
		t.Errorf("Expected name '用户详情', got '%s'", result.Info.Name)
	}

	if result.Params["id"] != "123" {
		t.Errorf("Expected param id='123', got '%s'", result.Params["id"])
	}
}

func TestWildcardMatch(t *testing.T) {
	matcher := NewMatcher(nil)

	// 添加通配符路径
	matcher.MustAddPath("/api/assets/*", WithName("资源文件"))

	// 测试通配符匹配
	result := matcher.Match("/api/assets/image.jpg")
	if !result.Matched {
		t.Error("Expected match for /api/assets/image.jpg, but got no match")
	}

	if result.Info.Name != "资源文件" {
		t.Errorf("Expected name '资源文件', got '%s'", result.Info.Name)
	}

	// 测试不匹配的情况
	result = matcher.Match("/api/assets/images/logo.png")
	if result.Matched {
		t.Error("Expected no match for /api/assets/images/logo.png with single wildcard, but got a match")
	}
}

func TestCatchAllMatch(t *testing.T) {
	matcher := NewMatcher(nil)

	// 添加多段通配符路径
	matcher.MustAddPath("/api/docs/**", WithName("文档资源"))

	// 测试多段通配符匹配
	result := matcher.Match("/api/docs/guide/intro")
	if !result.Matched {
		t.Error("Expected match for /api/docs/guide/intro, but got no match")
	}

	if result.Info.Name != "文档资源" {
		t.Errorf("Expected name '文档资源', got '%s'", result.Info.Name)
	}

	// 测试根路径匹配
	result = matcher.Match("/api/docs")
	if !result.Matched {
		t.Error("Expected match for /api/docs, but got no match")
	}
}

func TestCustomDelimiter(t *testing.T) {
	// 测试自定义分隔符
	matcher := NewMatcher(nil, WithDelimiter("."))

	// 添加使用点号分隔符的路径
	matcher.MustAddPath(".config.database.host", WithName("数据库主机"))

	// 测试匹配
	result := matcher.Match(".config.database.host")
	if !result.Matched {
		t.Error("Expected match for .config.database.host, but got no match")
	}

	if result.Info.Name != "数据库主机" {
		t.Errorf("Expected name '数据库主机', got '%s'", result.Info.Name)
	}

	// 测试参数匹配
	matcher.MustAddPath(".users.{id}.profile", WithName("用户资料"))
	result = matcher.Match(".users.123.profile")
	if !result.Matched {
		t.Error("Expected match for .users.123.profile, but got no match")
	}

	if result.Params["id"] != "123" {
		t.Errorf("Expected param id='123', got '%s'", result.Params["id"])
	}
}

func TestGroup(t *testing.T) {
	matcher := NewMatcher(nil)

	group := matcher.Group("/api")
	if group.Prefix() != "/api" {
		t.Errorf("Expected prefix '/api', got '%s'", group.Prefix())
	}
}

func TestMetaOption(t *testing.T) {
	matcher := NewMatcher(nil)
	matcher.MustAddPath("/api/test", WithName("测试接口"))

	result := matcher.Match("/api/test")
	if !result.Matched {
		t.Error("Expected match for /api/test, but got no match")
	}
}

func TestNewMatcherWithPatterns(t *testing.T) {
	patterns := []string{
		"/api/users",
		"/api/users/{id}",
		"/api/orders/{id}/items",
	}
	
	matcher := NewMatcher(patterns, WithDelimiter("/"))
	
	// 测试匹配
	result := matcher.Match("/api/users")
	if !result.Matched {
		t.Error("Expected match for /api/users, but got no match")
	}
	
	result = matcher.Match("/api/users/123")
	if !result.Matched {
		t.Error("Expected match for /api/users/123, but got no match")
	}
	
	if result.Params["id"] != "123" {
		t.Errorf("Expected param id='123', got '%s'", result.Params["id"])
	}
	
	result = matcher.Match("/api/orders/456/items")
	if !result.Matched {
		t.Error("Expected match for /api/orders/456/items, but got no match")
	}
	
	if result.Params["id"] != "456" {
		t.Errorf("Expected param id='456', got '%s'", result.Params["id"])
	}
}

func TestWithName(t *testing.T) {
	matcher := NewMatcher(nil)
	matcher.MustAddPath("/api/test", WithName("测试名称"))

	result := matcher.Match("/api/test")
	if !result.Matched {
		t.Error("Expected match for /api/test, but got no match")
	}

	if result.Info.Name != "测试名称" {
		t.Errorf("Expected name '测试名称', got '%s'", result.Info.Name)
	}
}

func TestWithDesc(t *testing.T) {
	matcher := NewMatcher(nil)
	matcher.MustAddPath("/api/test", WithDesc("测试描述"))

	result := matcher.Match("/api/test")
	if !result.Matched {
		t.Error("Expected match for /api/test, but got no match")
	}

	if result.Info.Desc != "测试描述" {
		t.Errorf("Expected desc '测试描述', got '%s'", result.Info.Desc)
	}
}

func TestWithMeta(t *testing.T) {
	meta := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	matcher := NewMatcher(nil)
	matcher.MustAddPath("/api/test", WithMeta(meta))

	result := matcher.Match("/api/test")
	if !result.Matched {
		t.Error("Expected match for /api/test, but got no match")
	}

	if result.Info.Meta["key1"] != "value1" {
		t.Errorf("Expected meta key1 'value1', got '%v'", result.Info.Meta["key1"])
	}

	if result.Info.Meta["key2"] != "value2" {
		t.Errorf("Expected meta key2 'value2', got '%v'", result.Info.Meta["key2"])
	}
}

func TestWithMetaMerge(t *testing.T) {
	meta1 := map[string]any{
		"key1": "value1",
	}

	meta2 := map[string]any{
		"key2": "value2",
	}

	matcher := NewMatcher(nil)
	// 我们不能在NewMatcher中使用WithMeta，因为它是Option类型而不是MatcherOption类型
	matcher.MustAddPath("/api/test", WithMeta(meta1), WithMeta(meta2))

	result := matcher.Match("/api/test")
	if !result.Matched {
		t.Error("Expected match for /api/test, but got no match")
	}

	if result.Info.Meta["key1"] != "value1" {
		t.Errorf("Expected meta key1 'value1', got '%v'", result.Info.Meta["key1"])
	}

	if result.Info.Meta["key2"] != "value2" {
		t.Errorf("Expected meta key2 'value2', got '%v'", result.Info.Meta["key2"])
	}
}

func TestWithDelimiter(t *testing.T) {
	matcher := NewMatcher(nil, WithDelimiter(":"))

	if matcher.delimiter != ":" {
		t.Errorf("Expected delimiter ':', got '%s'", matcher.delimiter)
	}
}
