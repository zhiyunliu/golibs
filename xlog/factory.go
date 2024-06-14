package xlog

import "sync"

var (
	_appenderCache = sync.Map{}
)

func RegistryBuilder(builder AppenderBuilder) {
	_appenderCache.Store(builder.Name(), builder)
}

func GetBuilder(name string) (AppenderBuilder, bool) {
	tmp, ok := _appenderCache.Load(name)
	if !ok {
		return nil, ok
	}
	return tmp.(AppenderBuilder), ok
}
