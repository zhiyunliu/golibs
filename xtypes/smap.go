package xtypes

import (
	"bytes"
	"encoding/json"
)

type SMap map[string]string

func (m SMap) Keys() []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	return keys
}

func (m SMap) Scan(obj interface{}) error {
	bytes, _ := json.Marshal(m)
	return json.Unmarshal(bytes, obj)
}

func (m SMap) Read(p []byte) (n int, err error) {
	dataBytes, _ := json.Marshal(m)
	return bytes.NewReader(dataBytes).Read(p)
}

func (m SMap) Get(name string) string {
	if v, ok := m[name]; ok {
		return v
	}
	return ""
}

func (m SMap) GetWithDefault(name string, def string) string {
	if v, ok := m[name]; ok {
		return v
	}
	return def
}

func (m SMap) Del(key string) {
	delete(m, key)
}

func (m SMap) Set(key, val string) {
	m[key] = val
}

func (m SMap) Values() map[string]string {
	return m
}
