package xtypes

import (
	"container/list"
	"encoding/json"
	"sync"
)

type SortCompare[K comparable] func(a, b K) bool

// Pair 表示一个键值对
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// SortedMap 是一个有序映射
type SortedMap[K comparable, V any] struct {
	mutex  *sync.Mutex
	list   *list.List
	less   func(a, b K) bool
	lookup map[K]*list.Element
}

// NewSortedMap 创建一个新的有序映射
func NewSortedMap[K comparable, V any](less SortCompare[K]) *SortedMap[K, V] {
	return &SortedMap[K, V]{
		mutex:  &sync.Mutex{},
		list:   list.New(),
		less:   less,
		lookup: make(map[K]*list.Element),
	}
}

// Put 向映射中添加或更新一个键值对
func (m *SortedMap[K, V]) Put(key K, value V) {
	if elem, exists := m.lookup[key]; exists {
		// 键已存在，更新值
		elem.Value = &Pair[K, V]{Key: key, Value: value}
	} else {
		// 键不存在，插入新的键值对
		newPair := &Pair[K, V]{Key: key, Value: value}
		elem := insertSorted(m.list, newPair, m.less)
		m.lookup[key] = elem
	}
}

func (m *SortedMap[K, V]) PutAll(mapval map[K]V) {
	for k, v := range mapval {
		m.Put(k, v)
	}
}

func (m *SortedMap[K, V]) ToMap() map[K]V {
	result := map[K]V{}
	for k, el := range m.lookup {
		result[k] = el.Value.(*Pair[K, V]).Value
	}
	return result
}

// Get 获取键对应的值
func (m *SortedMap[K, V]) Get(key K) (V, bool) {
	if elem, exists := m.lookup[key]; exists {
		return elem.Value.(*Pair[K, V]).Value, true
	}
	var zero V
	return zero, false
}

// Remove 从映射中移除一个键值对
func (m *SortedMap[K, V]) Remove(key K) {
	if elem, exists := m.lookup[key]; exists {
		m.list.Remove(elem)
		delete(m.lookup, key)
	}
}

// Contains 检查映射中是否包含某个键
func (m *SortedMap[K, V]) Contains(key K) bool {
	_, exists := m.lookup[key]
	return exists
}

// Keys 返回映射中的所有键
func (m *SortedMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.list.Len())
	for e := m.list.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*Pair[K, V]).Key)
	}
	return keys
}

// Values 返回映射中的所有值
func (m *SortedMap[K, V]) Values() []V {
	values := make([]V, 0, m.list.Len())
	for e := m.list.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.(*Pair[K, V]).Value)
	}
	return values
}

func (m *SortedMap[K, V]) Len() int {
	return len(m.lookup)
}

func (m *SortedMap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.list = list.New()
	m.lookup = make(map[K]*list.Element)
}

func (m SortedMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToMap())
}

func (m SortedMap[K, V]) MarshalBinary() ([]byte, error) {
	return m.MarshalJSON()
}

// Each 遍历映射中的所有键值对
func (m *SortedMap[K, V]) Each(f func(K, V)) {
	for e := m.list.Front(); e != nil; e = e.Next() {
		pair := e.Value.(*Pair[K, V])
		f(pair.Key, pair.Value)
	}
}

// insertSorted 将键值对插入到已排序的链表中
func insertSorted[K comparable, V any](l *list.List, pair *Pair[K, V], less SortCompare[K]) *list.Element {
	for e := l.Front(); e != nil; e = e.Next() {
		if less(pair.Key, e.Value.(*Pair[K, V]).Key) {
			return l.InsertBefore(pair, e)
		}
	}
	return l.PushBack(pair)
}
