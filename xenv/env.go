package xenv

import (
	"os"
	"sync"
)

var Env = sync.Map{}

func Get(key string) string {
	v, ok := Env.Load(key)
	if ok {
		return v.(string)
	}
	return os.Getenv(key)
}

func Set(key string, value string) {
	Env.Store(key, value)
}
