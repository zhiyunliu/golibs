package v2

type Option func(*Match)

// WithCache 启用或禁用缓存
func WithCache(enable bool) Option {
	return func(m *Match) {
		m.cache.enable = enable
	}
}

// WithDelimiter 设置分隔符，默认为"/"
func WithDelimiter(delimiter string) Option {
	return func(m *Match) {
		m.Delimiter = delimiter
	}
}
