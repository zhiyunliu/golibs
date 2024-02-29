package xlog

import (
	"encoding/json"
	"sync"

	"github.com/zhiyunliu/golibs/bytesconv"
)

type Formater func(e *Event, isJson bool) string

var (
	rwlocker   = sync.RWMutex{}
	partFormat = map[string]Formater{}
)

func Transform(key string, evt *Event, isJson bool) string {
	rwlocker.RLock()
	if call, ok := partFormat[key]; ok {
		rwlocker.RUnlock()
		return call(evt, isJson)
	}
	rwlocker.RUnlock()
	return ""
}

func RegistryFormater(key string, fmter Formater) {
	rwlocker.Lock()
	partFormat[key] = fmter
	rwlocker.Unlock()
}

func init() {
	RegistryFormater("tags", func(e *Event, isJson bool) string {
		bytes, _ := json.Marshal(e.Tags)
		return bytesconv.BytesToString(bytes)
	})
}
