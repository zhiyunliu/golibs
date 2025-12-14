package v2

import (
	"testing"
)

func TestMatcherError(t *testing.T) {
	err := NewMatcherError("测试错误", "/invalid/pattern")

	expectedMsg := "测试错误: /invalid/pattern"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Msg != "测试错误" {
		t.Errorf("Expected Msg '测试错误', got '%s'", err.Msg)
	}

	if err.Pattern != "/invalid/pattern" {
		t.Errorf("Expected Pattern '/invalid/pattern', got '%s'", err.Pattern)
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Pattern: "/invalid/pattern",
		Reason:  "测试原因",
	}

	expectedMsg := `invalid pattern "/invalid/pattern": 测试原因`
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	if err.Pattern != "/invalid/pattern" {
		t.Errorf("Expected Pattern '/invalid/pattern', got '%s'", err.Pattern)
	}

	if err.Reason != "测试原因" {
		t.Errorf("Expected Reason '测试原因', got '%s'", err.Reason)
	}
}