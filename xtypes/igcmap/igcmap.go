package igcmap

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xtypes"
)

/*
igcmap = ignore case map
忽略大小写的map 结构

*/

type Equalable interface {
	Equal(interface{}) bool
}

type IgcMap struct {
	data    xtypes.XMap
	keyMap  map[string]string
	rwmutex *sync.RWMutex
}

func New(orignal map[string]interface{}) *IgcMap {
	newOrgMap, keyMap := process(orignal)
	m := &IgcMap{
		data:    newOrgMap,
		keyMap:  keyMap,
		rwmutex: &sync.RWMutex{},
	}
	return m
}

// Set 添加新值
func (m *IgcMap) Set(key string, val interface{}) bool {
	lowerkey := strings.ToLower(key)
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()

	rk, ok := m.keyMap[lowerkey]
	if ok && rk != key {
		delete(m.data, rk)
	}
	m.data[key] = val
	m.keyMap[lowerkey] = key
	return ok
}

// Get 获取值
func (m IgcMap) Get(key string) (val interface{}, ok bool) {
	lowerkey := strings.ToLower(key)
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	rk, ok := m.keyMap[lowerkey]
	if !ok {
		return nil, false
	}
	return m.data[rk], true
}

// MergeMap 合并Map
func (m *IgcMap) MergeMap(other map[string]interface{}) {
	for k, v := range other {
		m.Set(k, v)
	}
}

// MergeIgc 合并MergeIgc
func (m *IgcMap) MergeIgc(other *IgcMap) {
	for k, v := range other.data {
		m.Set(k, v)
	}
}

// Iter 迭代每一个字段元素
func (m IgcMap) Iter(callback func(key string, val interface{}) bool) {
	for k, v := range m.data {
		m.rwmutex.RLock()
		if !callback(k, v) {
			m.rwmutex.RUnlock()
			break
		}
		m.rwmutex.RUnlock()
	}
}

// Keys 获取所有的键
func (m IgcMap) Keys() []string {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	keys := m.data.Keys()
	sort.Strings(keys)
	return keys
}

// Del 删除一个元素
func (m *IgcMap) Del(key string) {
	lowerkey := strings.ToLower(key)
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	rk, ok := m.keyMap[lowerkey]
	if !ok {
		return
	}
	delete(m.keyMap, lowerkey)
	delete(m.data, rk)
}

// Orignal 返回Map 的浅拷贝副本
func (m IgcMap) Orignal() map[string]interface{} {
	newVal := map[string]interface{}{}
	for k, v := range m.data {
		newVal[k] = v
	}
	return newVal
}

// Equal 判定两个对象是否相等
func (m *IgcMap) Equal(o *IgcMap) bool {
	if m == nil || o == nil {
		return false
	}

	if len(m.data) != len(o.data) {
		return false
	}

	for k, v := range m.keyMap {
		nk, ok := o.keyMap[k]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(m.data[v], o.data[nk]) {
			return false
		}
	}
	return true
}

func (d IgcMap) String() string {
	dataBytes, _ := json.Marshal(d.data)
	return bytesconv.BytesToString(dataBytes)
}

// Value String
func (d IgcMap) Value() (driver.Value, error) {
	return d.String(), nil
}

func (t *IgcMap) Scan(v interface{}) error {
	switch vt := v.(type) {
	case map[string]any:
		tmp := New(vt)
		*t = *tmp
	case xtypes.XMap:
		tmp := New(vt)
		*t = *tmp
	case []byte:
		mapData := xtypes.XMap{}
		err := json.Unmarshal(vt, &mapData)
		if err != nil {
			return fmt.Errorf("[]byte IgcMap.Scan err:%+v", err)
		}
		tmp := New(mapData)
		*t = *tmp
	case string:
		mapData := xtypes.XMap{}
		err := json.Unmarshal(bytesconv.StringToBytes(vt), &mapData)
		if err != nil {
			return fmt.Errorf("string IgcMap.Scan err:%+v", err)
		}
		tmp := New(mapData)
		*t = *tmp
	default:
		return fmt.Errorf("IgcMap类型处理错误:%+v", v)
	}
	return nil
}

func process(orginal map[string]interface{}) (newMap map[string]interface{}, keyMap map[string]string) {
	keyMap = make(map[string]string)
	newMap = make(map[string]interface{})
	if orginal == nil {
		return
	}
	var (
		ok       bool
		mk       string
		lowerkey string
	)

	xmapVal := xtypes.XMap(orginal)
	keys := xmapVal.Keys()
	sort.Strings(keys)

	for _, k := range keys {
		lowerkey = strings.ToLower(k)
		mk, ok = keyMap[lowerkey]
		if ok {
			delete(newMap, mk)
		}
		newMap[k] = orginal[k]
		keyMap[lowerkey] = k
	}
	return
}
