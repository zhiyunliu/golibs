package xtypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xtransform"
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

func (m SMap) ScanTo(obj interface{}) error {
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

// Deprecated
func (m SMap) GetWithDefault(name string, def string) string {
	return m.GetOrDefault(name, def)
}

// GetOrDefault gets value from SMap with a default value if the key doesn't exist
func (m SMap) GetOrDefault(name string, def string) string {
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

func (m SMap) Translate(tpl string) string {
	return xtransform.TranslateMap(tpl, m, xtransform.WithAtBraceMode(), xtransform.WithAtMode())
}

func (m SMap) MarshalBinary() (data []byte, err error) {
	tmp := map[string]string(m)
	return json.Marshal(tmp)
}

func (m *SMap) MapScan(obj any) error {
	return mapscan(obj, m)
}

// Value String
func (m SMap) Value() (driver.Value, error) {
	bytes, err := m.MarshalBinary()
	return bytesconv.BytesToString(bytes), err
}

func (m SMap) Len() int {
	return len(m)
}

func (m SMap) IsEmpty() bool {
	return len(m) <= 0
}

func (h SMap) Format(f fmt.State, verb rune) {
	bytes, _ := json.Marshal(h)
	_, _ = f.Write(bytes)
}
