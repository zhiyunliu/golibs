package v2

type Option func(*Matcher)

// WithCache 启用或禁用缓存
func WithCache(enable bool) Option {
	return func(m *Matcher) {
		m.cache.enable = enable
	}
}

// WithDelimiter 设置分隔符，默认为"/"
func WithDelimiter(delimiter string) Option {
	return func(m *Matcher) {
		m.Delimiter = delimiter
	}
}