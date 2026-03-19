package v2

type Option[T Pattern] func(*Match[T])

// WithCache 启用或禁用缓存
func WithCache[T Pattern](enable bool) Option[T] {
	return func(m *Match[T]) {
		m.cache.enable = enable
	}
}

// WithDelimiter 设置分隔符，默认为"/"
func WithDelimiter[T Pattern](delimiter string) Option[T] {
	return func(m *Match[T]) {
		m.Delimiter = delimiter
	}
}