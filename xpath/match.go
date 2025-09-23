package xpath

import (
	"regexp"
	"sort"
	"strings"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/zhiyunliu/golibs/bytesconv"
)

type Pattern interface {
	Pattern() string
}

type Patterns []Pattern

type StrPattern string

func (p StrPattern) Pattern() string {
	return string(p)
}

var specials = `~!@#$%^&*()_+-=<>?:"{}|,./;'[]\`

var tsmp = map[string]string{
	"**": `({0}[{1}\w]+)*`,
	"*":  `({0}[{1}\w]+)`,
}
var partsMp = map[string]string{
	"*": `[{1}\w]+`,
}

type sortString []Pattern

func (s sortString) Len() int { return len(s) }

func (s sortString) Swap(i, j int) {
	if i >= len(s) || j >= len(s) {
		return
	}
	s[i], s[j] = s[j], s[i]
}

func (s sortString) Less(i, j int) bool {
	iv := s[i].Pattern()
	jv := s[j].Pattern()
	il := len(iv)
	jl := len(jv)

	for x := 0; x < jl && x < il; x++ {
		if iv[x] == jv[x] {
			continue
		}
		if iv[x] == byte('*') {
			return false
		}
		if jv[x] == byte('*') {
			return true
		}
		return iv[x] < jv[x]
	}
	return iv < jv
}

// Match 构建模糊匹配缓存查找管理器
type Match struct {
	mutex     sync.Mutex
	cache     *matchCacheWrap
	all       []Pattern
	regexpAll []*regexp.Regexp
}
type matchCacheWrap struct {
	enbale   bool
	cacheMap cmap.ConcurrentMap[string, Pattern]
}

func (w *matchCacheWrap) Get(key string) (val Pattern, ok bool) {
	if !w.enbale {
		return
	}
	val, ok = w.cacheMap.Get(key)
	return
}
func (w *matchCacheWrap) SetIfAbsent(key string, val Pattern) bool {
	if !w.enbale {
		return false
	}
	return w.cacheMap.SetIfAbsent(key, val)
}

func NewMatch(pathList []string, opts ...Option) *Match {
	patterns := make([]Pattern, len(pathList))
	for i := range pathList {
		patterns[i] = StrPattern(pathList[i])
	}
	return NewMatchPatterns(patterns, opts...)
}

// NewMatch 构建模糊匹配缓存查找管理器
func NewMatchPatterns(pathList []Pattern, opts ...Option) *Match {
	m := &Match{
		cache: &matchCacheWrap{
			cacheMap: cmap.New[Pattern](),
		},
		all: pathList,
	}
	for i := range opts {
		opts[i](m)
	}
	sort.Sort(sortString(m.all))

	m.regexpAll = make([]*regexp.Regexp, len(m.all))

	return m
}

func (m *Match) CanUseCache() bool {
	return m.cache.enbale
}

func (m *Match) Match(path string, spls ...string) (match bool, pattern string) {
	match, tmpPattern := m.MatchPattern(path, spls...)
	if tmpPattern != nil {
		pattern = tmpPattern.Pattern()
	}
	return match, pattern
}

// M

// Match Match
func (m *Match) MatchPattern(path string, spls ...string) (match bool, pattern Pattern) {
	sep := "/"
	if len(spls) > 0 {
		sep = spls[0]
	}

	cacheKey := ""
	if m.CanUseCache() {
		cacheKey = m.buildCacheKey(path, sep)
		if val, ok := m.cache.Get(cacheKey); ok {
			return true, val
		}
	}

	for i, u := range m.all {
		if strings.EqualFold(u.Pattern(), path) {
			if m.CanUseCache() {
				m.cache.SetIfAbsent(cacheKey, u)
			}
			return true, u
		}
		regp := m.getRegexp(u.Pattern(), i, sep)
		match = regp.Match(bytesconv.StringToBytes(path))
		if match {
			if m.CanUseCache() {
				m.cache.SetIfAbsent(cacheKey, u)
			}
			return match, u
		}
	}
	return false, nil
}

func (m *Match) buildCacheKey(path, sep string) string {
	return sep + ":" + path
}

func (m *Match) getRegexp(u string, idx int, sep string) *regexp.Regexp {
	if m.regexpAll[idx] == nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		if m.regexpAll[idx] != nil {
			return m.regexpAll[idx]
		}
		parties := strings.Split(u, sep)
		npts := make([]string, len(parties))
		curSpecials := m.processSpecial(strings.ReplaceAll(specials, sep, ""))
		sep = m.processSpecial(sep)

		for i := range parties {
			if parties[i] == "" {
				continue
			}
			pv, ok := tsmp[parties[i]]
			if !ok {
				nv := m.processSpecial(parties[i])
				if !strings.Contains(nv, "*") {
					pv = nv
					if i > 0 {
						pv = sep + nv
					}
				} else {
					pv = sep + strings.ReplaceAll(nv, `*`, partsMp["*"])
				}
			}
			sl := sep
			if i <= 0 {
				sl = ""
			}

			pv = strings.Replace(pv, "{0}", sl, -1)
			npts[i] = strings.ReplaceAll(pv, "{1}", curSpecials)
		}
		m.regexpAll[idx] = regexp.MustCompile("^(" + strings.Join(npts, "") + ")$")
	}
	return m.regexpAll[idx]
}

func (m *Match) processSpecial(nv string) string {
	nv = strings.ReplaceAll(nv, `\`, `\\`)
	nv = strings.ReplaceAll(nv, "$", `\$`)
	nv = strings.ReplaceAll(nv, "(", `\(`)
	nv = strings.ReplaceAll(nv, ")", `\)`)
	nv = strings.ReplaceAll(nv, "+", `\+`)
	nv = strings.ReplaceAll(nv, ".", `\.`)
	nv = strings.ReplaceAll(nv, "[", `\[`)
	nv = strings.ReplaceAll(nv, "]", `\]`)
	nv = strings.ReplaceAll(nv, "?", `\?`)
	nv = strings.ReplaceAll(nv, "^", `\^`)
	nv = strings.ReplaceAll(nv, "{", `\{`)
	nv = strings.ReplaceAll(nv, "|", `\|`)
	nv = strings.ReplaceAll(nv, "-", `\-`)
	return nv
}
