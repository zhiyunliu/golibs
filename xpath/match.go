package xpath

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	cmap "github.com/orcaman/concurrent-map"
)

var specials = `~!@#$%^&*()_+-=<>?:"{}|,./;'[]\`

var tsmp = map[string]string{
	"**": `({0}[{1}\w]+)*`,
	"*":  `({0}[{1}\w]+)`,
}

type sortString []string

func (s sortString) Len() int { return len(s) }

func (s sortString) Swap(i, j int) {
	if i >= len(s) || j >= len(s) {
		return
	}
	s[i], s[j] = s[j], s[i]
}

func (s sortString) Less(i, j int) bool {
	il := len(s[i])
	jl := len(s[j])
	for x := 0; x < jl && x < il; x++ {
		if s[i][x] == s[j][x] {
			continue
		}
		if s[i][x] == []byte("*")[0] {
			return false
		}
		if s[j][x] == []byte("*")[0] {
			return true
		}
		return s[i][x] < s[j][x]
	}
	return s[i] < s[j]
}

//Match 构建模糊匹配缓存查找管理器
type Match struct {
	cache     cmap.ConcurrentMap
	all       []string
	regexpAll []string
}

//NewMatch 构建模糊匹配缓存查找管理器
func NewMatch(all ...string) *Match {
	m := &Match{
		cache: cmap.New(),
		all:   all,
	}
	sort.Sort(sortString(m.all))

	m.regexpAll = make([]string, len(m.all))

	return m
}

//Match Match
func (m *Match) Match(path string, spls ...string) (match bool, pattern string) {
	var err error
	sep := "/"
	if len(spls) > 0 {
		sep = spls[0]
	}
	key := m.buildCacheKey(path, sep)
	if val, ok := m.cache.Get(key); ok {
		return true, val.(string)
	}

	for i, u := range m.all {
		if strings.EqualFold(u, path) {
			m.cache.SetIfAbsent(key, u)
			return true, u
		}
		regp := m.getRegexp(u, i, sep)
		//fmt.Println("regp:", regp)
		match, err = regexp.Match(regp, []byte(path))
		if err != nil {
			match = false
		}
		if match {
			m.cache.SetIfAbsent(key, u)
			return match, u
		}
	}
	return false, ""
}

func (m *Match) buildCacheKey(path, sep string) string {
	return fmt.Sprintf("%s:%s", sep, path)
}

func (m *Match) getRegexp(u string, idx int, sep string) string {
	if m.regexpAll[idx] == "" {
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
					pv = strings.ReplaceAll(nv, `\*`, tsmp["*"])
				}
			}
			sl := sep
			if i <= 0 {
				sl = ""
			}

			pv = strings.Replace(pv, "{0}", sl, -1)
			npts[i] = strings.ReplaceAll(pv, "{1}", curSpecials)
		}
		m.regexpAll[idx] = "^(" + strings.Join(npts, "") + ")$"
	}
	return m.regexpAll[idx]
}

func (m *Match) processSpecial(nv string) string {
	nv = strings.ReplaceAll(nv, `\`, `\\`)
	nv = strings.ReplaceAll(nv, "$", `\$`)
	nv = strings.ReplaceAll(nv, "(", `\(`)
	nv = strings.ReplaceAll(nv, ")", `\)`)
	nv = strings.ReplaceAll(nv, "*", `\*`)
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
