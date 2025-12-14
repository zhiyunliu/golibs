package v2

import (
	"testing"
)

func TestGroupBasics(t *testing.T) {
	matcher := NewMatcher()

	// 创建分组
	apiGroup := matcher.Group("/api")
	
	// 在分组中添加路径
	apiGroup.MustAddPath("/users", WithName("用户列表"))
	apiGroup.MustAddPath("/users/{id}", WithName("用户详情"))

	// 测试匹配
	result := matcher.Match("/api/users")
	if !result.Matched {
		t.Error("Expected match for /api/users, but got no match")
	}

	if result.Info.Name != "用户列表" {
		t.Errorf("Expected name '用户列表', got '%s'", result.Info.Name)
	}

	result = matcher.Match("/api/users/123")
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

func TestNestedGroups(t *testing.T) {
	matcher := NewMatcher()

	// 创建嵌套分组
	apiGroup := matcher.Group("/api")
	v1Group := apiGroup.Group("/v1")
	usersGroup := v1Group.Group("/users")

	// 在最内层分组中添加路径
	usersGroup.MustAddPath("/{id}", WithName("用户详情"))
	usersGroup.MustAddPath("/{id}/profile", WithName("用户资料"))

	// 测试匹配
	result := matcher.Match("/api/v1/users/123")
	if !result.Matched {
		t.Error("Expected match for /api/v1/users/123, but got no match")
	}

	if result.Info.Name != "用户详情" {
		t.Errorf("Expected name '用户详情', got '%s'", result.Info.Name)
	}

	if result.Params["id"] != "123" {
		t.Errorf("Expected param id='123', got '%s'", result.Params["id"])
	}

	result = matcher.Match("/api/v1/users/123/profile")
	if !result.Matched {
		t.Error("Expected match for /api/v1/users/123/profile, but got no match")
	}

	if result.Info.Name != "用户资料" {
		t.Errorf("Expected name '用户资料', got '%s'", result.Info.Name)
	}

	if result.Params["id"] != "123" {
		t.Errorf("Expected param id='123', got '%s'", result.Params["id"])
	}
}

func TestGroupWithOptions(t *testing.T) {
	matcher := NewMatcher()

	// 创建带选项的分组
	versionMeta := map[string]any{
		"version": "v1",
		"prefix":  "/api/v1",
	}

	apiGroup := matcher.Group("/api", WithMeta(versionMeta))

	// 在分组中添加路径
	apiGroup.MustAddPath("/users", WithName("用户列表"))

	// 测试匹配和元数据
	result := matcher.Match("/api/users")
	if !result.Matched {
		t.Error("Expected match for /api/users, but got no match")
	}

	if result.Info.Name != "用户列表" {
		t.Errorf("Expected name '用户列表', got '%s'", result.Info.Name)
	}

	// 检查来自分组的元数据是否合并到了路径信息中
	if result.Info.Meta["version"] != "v1" {
		t.Errorf("Expected meta version 'v1', got '%v'", result.Info.Meta["version"])
	}

	if result.Info.Meta["prefix"] != "/api/v1" {
		t.Errorf("Expected meta prefix '/api/v1', got '%v'", result.Info.Meta["prefix"])
	}
}

func TestGroupWithValidation(t *testing.T) {
	matcher := NewMatcher()

	apiGroup := matcher.Group("/api")

	// 测试带验证的路径添加
	err := apiGroup.AddPathWithValidation("/users/{id}")
	if err != nil {
		t.Errorf("Expected to add path with validation successfully, but got error: %v", err)
	}

	// 测试无效路径
	err = apiGroup.AddPathWithValidation("/users/{}/profile")
	if err == nil {
		t.Error("Expected error when adding invalid path with validation, but got no error")
	}
}
