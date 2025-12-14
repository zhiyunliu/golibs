package v2

import "strings"

// ValidatePattern 验证路径模式
func ValidatePattern(pattern string) error {
	if pattern == "" {
		return &ValidationError{Pattern: pattern, Reason: "empty pattern"}
	}

	if pattern[0] != '/' {
		return &ValidationError{Pattern: pattern, Reason: "must start with '/'"}
	}

	segments := splitPath(pattern)
	for i, seg := range segments {
		if seg == "" {
			return &ValidationError{Pattern: pattern, Reason: "empty segment"}
		}

		// 检查**只能出现在最后
		if seg == "**" && i != len(segments)-1 {
			return &ValidationError{
				Pattern: pattern,
				Reason:  "catch-all (**) must be the last segment",
			}
		}

		// 检查参数格式
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			paramName := seg[1 : len(seg)-1]
			if paramName == "" {
				return &ValidationError{Pattern: pattern, Reason: "empty parameter name"}
			}
			// 检查参数名有效性
			if !isValidParamName(paramName) {
				return &ValidationError{
					Pattern: pattern,
					Reason:  "invalid parameter name: " + paramName,
				}
			}
		}
	}

	return nil
}

// isValidParamName 检查参数名是否有效
func isValidParamName(name string) bool {
	if name == "" {
		return false
	}

	// 参数名只能包含字母、数字和下划线
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_') {
			return false
		}
	}
	return true
}

// AddPathWithValidation 带验证的添加路径
func (m *Matcher) AddPathWithValidation(pattern string, options ...Option) error {
	if err := ValidatePattern(pattern); err != nil {
		return err
	}
	return m.AddPath(pattern, options...)
}