package v2

import "fmt"

// MatcherError 匹配器错误类型
type MatcherError struct {
	Msg     string
	Pattern string
}

func NewMatcherError(msg, pattern string) *MatcherError {
	return &MatcherError{
		Msg:     msg,
		Pattern: pattern,
	}
}

func (e *MatcherError) Error() string {
	return e.Msg + ": " + e.Pattern
}

// ValidationError 验证错误
type ValidationError struct {
	Pattern string
	Reason  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid pattern %q: %s", e.Pattern, e.Reason)
}