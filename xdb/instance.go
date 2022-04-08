package xdb

import (
	"fmt"
	"sync"
)

var cache sync.Map

func init() {
	cache = sync.Map{}
}

func getDbCacheKey(name string) string {
	return fmt.Sprintf("db_config_%s", name)
}

func SetDbConfig(name string, cfg *Config) {
	cache.Store(getDbCacheKey(name), cfg)
}

func GetDbConfig(name string) *Config {
	val, ok := cache.Load(getDbCacheKey(name))
	if !ok {
		panic(fmt.Errorf("不存在DB=%s的配置", name))
	}
	return val.(*Config)
}

func GetDB(name string) IDB {
	key := fmt.Sprintf("db_instance_%s", name)

	obj, ok := cache.Load(key)
	if !ok {
		dbcfg := GetDbConfig(name)
		instance, err := NewDB(dbcfg.Proto, dbcfg.ConnString, dbcfg.MaxOpen, dbcfg.MaxIdle, dbcfg.LifeTime)
		if err != nil {
			panic(fmt.Errorf("创建数据库失败:%w,name=%s", err, name))
		}
		cache.Store(key, instance)
		obj = instance
	}
	return obj.(IDB)
}
