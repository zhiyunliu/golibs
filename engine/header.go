package engine

import (
	"encoding/json"
	"fmt"
	"net/textproto"
	"strings"
)

type Header map[string]string

func (h Header) Get(key string) string {
	//兼容旧版本的header可能直接输入导致的大小写不一致问题
	if v, ok := h[key]; ok {
		return v
	}
	key = textproto.CanonicalMIMEHeaderKey(key)
	if v, ok := h[key]; ok {
		return v
	}
	key = strings.ToLower(key)
	if v, ok := h[key]; ok {
		return v
	}
	return ""
}

func (h Header) Set(key, value string) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	h[key] = value
}

func (h Header) Del(key string) {
	delete(h, key)
}

func (h Header) Keys() []string {
	keys := make([]string, len(h))
	idx := 0
	for k := range h {
		keys[idx] = k
		idx++
	}
	return keys
}

func (h Header) Len() int {
	return len(h)
}
func (h Header) IsEmpty() bool {
	return len(h) == 0
}
func (h Header) Values() map[string]string {
	return h
}
func (h Header) Format(f fmt.State, verb rune) {
	bytes, _ := json.Marshal(h)
	_, _ = f.Write(bytes)
}
