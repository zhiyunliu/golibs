package xtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xtransform"
)

type XMap map[string]interface{}

// Keys 从对象中获取数据值，如果不是字符串则返回空
func (m XMap) Keys() []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	return keys
}

// Merge 合并
func (m XMap) Merge(r XMap) {
	for k, v := range r {
		m[k] = v
	}
}

// Get 获取指定元素的值
func (m XMap) Get(name string) (interface{}, bool) {
	v, ok := m[name]
	return v, ok
}

// Scan 以json 标签进行序列化
func (m XMap) Scan(obj interface{}) error {
	bytes, _ := json.Marshal(m)
	return json.Unmarshal(bytes, obj)
}

func (m XMap) Len() int {
	return len(m)
}

func (m XMap) IsEmpty() bool {
	return len(m) == 0
}

func (m XMap) SMap() SMap {
	sm := map[string]string{}
	for k, v := range m {
		sm[k] = fmt.Sprintf("%+v", v)
	}
	return sm
}

func (m XMap) Translate(tpl string) string {
	return xtransform.Translate(tpl, m)
}

func (m XMap) MarshalBinary() (data []byte, err error) {
	tmp := map[string]interface{}(m)
	return json.Marshal(tmp)
}

// Get 获取指定元素的Bool值
func (m XMap) GetBool(name string) bool {
	v, ok := m[name]
	if !ok {
		return false
	}
	return GetBool(v)
}

// Get 获取指定元素的值
func (m XMap) GetString(name string) string {
	v, ok := m[name]
	if !ok {
		return ""
	}
	return GetString(v)
}

func (m XMap) GetInt(key string) (int, error) {
	tmp, ok := m[key]
	if !ok || tmp == nil {
		return 0, nil
	}

	return GetInt(tmp)
}

func (m XMap) GetInt64(key string) (int64, error) {
	tmp, ok := m[key]
	if !ok || tmp == nil {
		return 0, nil
	}

	return GetInt64(tmp)
}

func (m XMap) GetFloat64(key string) (float64, error) {
	tmp, ok := m[key]
	if !ok || tmp == nil {
		return 0, nil
	}
	return GetFloat64(tmp)
}

func (m *XMap) MapScan(obj interface{}) error {
	return mapscan(obj, m)
}

// Value String
func (m XMap) Value() (driver.Value, error) {
	bytes, err := m.MarshalBinary()
	return bytesconv.BytesToString(bytes), err
}
