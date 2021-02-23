package pathmatch

import (
	"regexp"
	"sort"
	"strings"

	cmap "github.com/orcaman/concurrent-map"
)

var tsmp = map[string]string{
	"/**": `(/+\w+[\.\w]+)*`,
	".**": `(\.+\w+[\.\w]+)*`,
	"**":  `(\w+[\.\w]+)*`,
	"*":   `(\w+[\.\w]*)`,
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

//PathMatch 构建模糊匹配缓存查找管理器
//使用/** 和 /* 进行模式匹配
type PathMatch struct {
	cache     cmap.ConcurrentMap
	all       []string
	regexpAll []string
}

//NewMatch 构建模糊匹配缓存查找管理器
func NewMatch(all ...string) *PathMatch {
	m := &PathMatch{
		cache: cmap.New(),
		all:   all,
	}
	sort.Sort(sortString(m.all))

	m.regexpAll = make([]string, len(m.all))

	return m
}

//Match Match
func (m *PathMatch) Match(path string, spl ...string) (match bool, pattern string) {
	var err error
	sep := "/"
	if len(spl) > 0 {
		sep = spl[0]
	}
	for i, u := range m.all {
		if strings.EqualFold(u, path) {
			m.cache.SetIfAbsent(path, u)
			return true, u
		}
		regp := m.getRegexp(u, i, sep)

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

func (m *PathMatch) getRegexp(u string, idx int, sep string) string {
	if m.regexpAll[idx] == "" {
		nv := u
		nv = strings.ReplaceAll(nv, ".", `\.`)
		nv = strings.ReplaceAll(nv, "+", `\+`)
		nv = strings.ReplaceAll(nv, "$", `\$`)
		nv = strings.ReplaceAll(nv, "^", `\^`)

		parties := strings.Split(nv, sep)
		npts := make([]string, len(parties))
		for i := range parties {
			if parties[i] == "" {
				continue
			}
			pv, ok := tsmp[sep+parties[i]]
			if !ok {
				pv = sep + strings.ReplaceAll(parties[i], "*", tsmp["*"])
			}
			npts[i] = pv
		}
		m.regexpAll[idx] = "^(" + strings.Join(npts, "") + ")$"
	}
	return m.regexpAll[idx]
}
