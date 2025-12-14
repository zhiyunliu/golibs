package v2

import (
	"testing"
)

func TestBasicMatch(t *testing.T) {
	matcher := NewMatcher()

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
	matcher := NewMatcher()

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
	matcher := NewMatcher()

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
	matcher := NewMatcher()

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
	matcher := NewMatcher(WithDelimiter("."))

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
	matcher := NewMatcher()

	// 创建分组
	group := matcher.Group("/api/v1")

	// 在分组中添加路径
	group.MustAddPath("/users", WithName("用户列表"))
	group.MustAddPath("/users/{id}", WithName("用户详情"))

	// 测试分组中的路径匹配
	result := matcher.Match("/api/v1/users")
	if !result.Matched {
		t.Error("Expected match for /api/v1/users, but got no match")
	}

	if result.Info.Name != "用户列表" {
		t.Errorf("Expected name '用户列表', got '%s'", result.Info.Name)
	}

	result = matcher.Match("/api/v1/users/123")
	if !result.Matched {
		t.Error("Expected match for /api/v1/users/123, but got no match")
	}

	if result.Info.Name != "用户详情" {
		t.Errorf("Expected name '用户详情', got '%s'", result.Info.Name)
	}

	if result.Params["id"] != "123" {
		t.Errorf("Expected param id='123', got '%s'", result.Params["id"])
	}
}

func TestMetaOption(t *testing.T) {
	matcher := NewMatcher()

	metaData := map[string]any{
		"version": "v1",
		"auth":    true,
		"timeout": 30,
	}

	matcher.MustAddPath("/api/data", WithMeta(metaData))

	result := matcher.Match("/api/data")
	if !result.Matched {
		t.Error("Expected match for /api/data, but got no match")
	}

	if result.Info.Meta["version"] != "v1" {
		t.Errorf("Expected meta version 'v1', got '%v'", result.Info.Meta["version"])
	}

	if result.Info.Meta["auth"] != true {
		t.Errorf("Expected meta auth true, got '%v'", result.Info.Meta["auth"])
	}

	if result.Info.Meta["timeout"] != 30 {
		t.Errorf("Expected meta timeout 30, got '%v'", result.Info.Meta["timeout"])
	}
}

func TestNewMatcherWithPatterns(t *testing.T) {
	patterns := []string{
		"/api/users",
		"/api/users/{id}",
		"/api/orders/{id}/items",
	}
	
	matcher := NewMatcherWithPatterns(patterns, WithDelimiter("/"))
	
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
