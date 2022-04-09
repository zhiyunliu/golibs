package xtypes

type XMaps []XMap

func (ms *XMaps) Append(i ...XMap) XMaps {
	*ms = append(*ms, i...)
	return *ms
}

func (ms XMaps) IsEmpty() bool {
	return ms == nil || len(ms) == 0
}

func (ms XMaps) Len() int {
	return len(ms)
}

func (ms XMaps) Get(idx int) XMap {
	if idx < 0 || len(ms) <= idx {
		return map[string]interface{}{}
	}
	if len(ms) > idx {
		return ms[idx]
	}
	return map[string]interface{}{}
}
