package match

import (
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

//GetValueByIndex 获取数组中的指定元素
func GetValueByIndex(v []string, index int, def ...string) string {
	if len(v) > index {
		return v[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

//ValueMatch 构建模糊匹配缓存查找管理器
type ValueMatch struct {
	cache     cmap.ConcurrentMap
	all       []string
	regexpAll []string
}

//NewMatch 构建模糊匹配缓存查找管理器
func NewMatch(all ...string) *ValueMatch {
	m := &ValueMatch{
		cache: cmap.New(),
		all:   all,
	}
	sort.Sort(sortString(m.all))

	m.regexpAll = make([]string, len(m.all))

	return m
}

//Match Match
func (m *ValueMatch) Match(path string, spl ...string) (match bool, pattern string) {
	var err error
	sep := GetValueByIndex(spl, 0, "/")
	for i, u := range m.all {
		if strings.EqualFold(u, path) {
			m.cache.SetIfAbsent(path, u)
			return true, u
		}
		regp := m.getRegexp(u, i, sep)
		//fmt.Println("regp:", regp)
		match, err = regexp.Match(regp, []byte(path))
		if err != nil {
			match = false
		}
		if match {
			m.cache.SetIfAbsent(path, u)
			return match, u
		}
	}
	return false, ""
}

func (m *ValueMatch) getRegexp(u string, idx int, sep string) string {
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

func (m *ValueMatch) processSpecial(nv string) string {
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
