package v2

// Group 路径分组
type Group struct {
	matcher   *Matcher
	prefix    string
	options   []Option
	delimiter string
}

// NewGroup 创建路径分组
func (m *Matcher) Group(prefix string, options ...Option) *Group {
	return &Group{
		matcher:   m,
		prefix:    prefix,
		options:   options,
		delimiter: m.delimiter,
	}
}

// AddPath 分组添加路径
func (g *Group) AddPath(pattern string, options ...Option) error {
	fullPattern := g.matcher.joinPath(g.prefix, pattern)

	// 合并分组选项和路径选项
	allOptions := make([]Option, 0, len(g.options)+len(options))
	allOptions = append(allOptions, g.options...)
	allOptions = append(allOptions, options...)

	return g.matcher.AddPath(fullPattern, allOptions...)
}

// MustAddPath 分组添加路径（panic版本）
func (g *Group) MustAddPath(pattern string, options ...Option) {
	if err := g.AddPath(pattern, options...); err != nil {
		panic(err)
	}
}

// AddPathWithValidation 带验证的分组添加路径
func (g *Group) AddPathWithValidation(pattern string, options ...Option) error {
	fullPattern := g.matcher.joinPath(g.prefix, pattern)

	if err := g.matcher.ValidatePattern(fullPattern); err != nil {
		return err
	}

	// 合并分组选项和路径选项
	allOptions := make([]Option, 0, len(g.options)+len(options))
	allOptions = append(allOptions, g.options...)
	allOptions = append(allOptions, options...)

	return g.matcher.AddPath(fullPattern, allOptions...)
}

// Group 创建子分组
func (g *Group) Group(prefix string, options ...Option) *Group {
	fullPrefix := g.matcher.joinPath(g.prefix, prefix)

	// 合并父分组和子分组的选项
	allOptions := make([]Option, 0, len(g.options)+len(options))
	allOptions = append(allOptions, g.options...)
	allOptions = append(allOptions, options...)

	return &Group{
		matcher:   g.matcher,
		prefix:    fullPrefix,
		options:   allOptions,
		delimiter: g.delimiter,
	}
}
