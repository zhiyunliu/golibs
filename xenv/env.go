package xenv

import (
	"os"
	"sync"
)

var Env = sync.Map{}

// Get env value
func Get(key string) string {
	v, ok := Env.Load(key)
	if ok {
		return v.(string)
	}
	return os.Getenv(key)
}

// GetOrDefault gets env value with default value if not exists
func GetOrDefault(key string, defaultVal string) string {
	v, ok := Env.Load(key)
	if ok {
		return v.(string)
	}

	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// Set env value
func Set(key string, value string) {
	Env.Store(key, value)
}
