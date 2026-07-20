package v2

import (
	"sort"
	"strings"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/zhiyunliu/golibs/xpath"
)

type Pattern = xpath.Pattern
type Patterns = xpath.Patterns
type StrPattern = xpath.StrPattern

// Match 基于前缀树的匹配器
type Match[T Pattern] struct {
	patterns  []T
	Delimiter string
	cache     *matchCacheWrap[T]
}

type matchCacheWrap[T Pattern] struct {
	enable   bool
	cacheMap cmap.ConcurrentMap[string, T]
}

func (w *matchCacheWrap[T]) Get(key string) (val T, ok bool) {
	if !w.enable {
		return
	}
	val, ok = w.cacheMap.Get(key)
	return
}

func (w *matchCacheWrap[T]) SetIfAbsent(key string, val T) bool {
	if !w.enable {
		return false
	}
	return w.cacheMap.SetIfAbsent(key, val)
}

// NewMatcher 创建一个新的匹配器
func NewMatch(pathList []string, opts ...Option[Pattern]) *Match[Pattern] {
	patterns := make([]Pattern, len(pathList))
	for i := range pathList {
		patterns[i] = StrPattern(pathList[i])
	}
	return NewMatcherPatterns(patterns, opts...)
}

// NewMatcherPatterns 使用Pattern切片创建匹配器
func NewMatcherPatterns[T Pattern](pathList []T, opts ...Option[T]) *Match[T] {
	m := &Match[T]{
		patterns:  pathList,
		Delimiter: "/",
		cache: &matchCacheWrap[T]{
			cacheMap: cmap.New[T](),
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	sort.Sort((SortPatterns[T](m.patterns)))

	return m
}

// Match 尝试匹配给定的路径
func (m *Match[T]) Match(path string) (match bool, pattern T) {
	sep := m.Delimiter
	cacheKey := ""
	if m.CanUseCache() {
		cacheKey = m.buildCacheKey(path, sep)
		if val, ok := m.cache.Get(cacheKey); ok {
			return true, val
		}
	}

	for _, p := range m.patterns {
		if matchPathPattern(path, p.Pattern(), sep) {
			if m.CanUseCache() {
				m.cache.SetIfAbsent(cacheKey, p)
			}
			return true, p
		}
	}
	return false, *new(T)
}

// matchPathPattern 检查路径是否匹配给定模式
func matchPathPattern(path, pattern, sep string) bool {
	// 直接相等的情况
	if path == pattern {
		return true
	}

	pathParts := strings.Split(path, sep)
	patternParts := strings.Split(pattern, sep)

	// 处理以分隔符开头的情况
	if len(pathParts) > 0 && pathParts[0] == "" {
		pathParts = pathParts[1:]
	}
	if len(patternParts) > 0 && patternParts[0] == "" {
		patternParts = patternParts[1:]
	}

	return matchParts(pathParts, patternParts)
}

// matchParts 匹配路径段数组
func matchParts(pathParts, patternParts []string) bool {
	return matchPartsRecursive(pathParts, patternParts, 0, 0)
}

// matchPartsRecursive 递归匹配路径段
func matchPartsRecursive(pathParts, patternParts []string, pathIdx, patternIdx int) bool {
	// 如果模式已经遍历完，只有当路径也遍历完时才匹配
	if patternIdx == len(patternParts) {
		return pathIdx == len(pathParts)
	}

	// 检查当前模式段
	currentPattern := patternParts[patternIdx]

	// 如果是 ** 模式
	if currentPattern == "**" {
		// ** 可以匹配0个或多个路径段
		// 尝试匹配0个路径段（即跳过**）
		if matchPartsRecursive(pathParts, patternParts, pathIdx, patternIdx+1) {
			return true
		}
		// 尝试匹配1个或多个路径段
		for i := pathIdx; i < len(pathParts); i++ {
			if matchPartsRecursive(pathParts, patternParts, i+1, patternIdx+1) {
				return true
			}
		}
		return false
	}

	// 如果路径已经遍历完但模式还有剩余，且下一个模式不是**
	if pathIdx == len(pathParts) {
		return false
	}

	// 检查当前路径段是否匹配当前模式段
	currentPath := pathParts[pathIdx]

	if currentPattern == "*" || matchLiteral(currentPattern, currentPath) {
		// 当前段匹配，继续匹配剩余部分
		return matchPartsRecursive(pathParts, patternParts, pathIdx+1, patternIdx+1)
	}

	return false
}

// matchLiteral 检查字面量模式是否匹配当前路径段，支持 * 通配符
func matchLiteral(pattern, value string) bool {
	// 如果模式中不包含 *，则直接比较
	if !strings.Contains(pattern, "*") {
		return pattern == value
	}

	// 使用简单的通配符匹配
	return simpleMatch(value, pattern)
}

// simpleMatch 实现简单的通配符匹配
func simpleMatch(text, pattern string) bool {
	tIdx, pIdx := 0, 0
	tLen, pLen := len(text), len(pattern)

	var starIdx, matchIdx int = -1, -1

	for tIdx < tLen {
		// 如果字符匹配或模式是 '?'
		if pIdx < pLen && (pattern[pIdx] == '*' || pattern[pIdx] == text[tIdx]) {
			if pattern[pIdx] == '*' {
				starIdx = pIdx
				matchIdx = tIdx
				pIdx++
			} else {
				pIdx++
				tIdx++
			}
		} else if starIdx != -1 { // 如果遇到不匹配但之前有 '*'
			pIdx = starIdx + 1
			matchIdx++
			tIdx = matchIdx
		} else { // 如果不匹配且没有 '*'，返回 false
			return false
		}
	}

	// 检查模式剩余部分是否都是 '*'
	for pIdx < pLen {
		if pattern[pIdx] != '*' {
			return false
		}
		pIdx++
	}

	return true
}

// CanUseCache 检查是否可以使用缓存
func (m *Match[T]) CanUseCache() bool {
	return m.cache.enable
}

func (m *Match[T]) buildCacheKey(path, sep string) string {
	return sep + ":" + path
}

// Count 返回匹配器中模式的数量
func (m *Match[T]) Count() int {
	return len(m.patterns)
}
