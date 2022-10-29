package xpath

type Option func(*Match)

func WithCache(enable bool) Option {
	return func(m *Match) {
		m.cache.enbale = enable
	}
}
