package xlog

import (
	"fmt"
	"sync"
)

var (
	_appenderCache = sync.Map{}
	_defaultLayout = "[%time][%l][%session][%idx] %content"
)

func Registry(builder AppenderBuilder) {
	_appenderCache.Store(builder.Name(), builder)
}

type AppenderBuilder interface {
	Name() string
	DefaultLayout() *Layout
	Build(layout *Layout) Appender
}

//Appender 定义appender接口
type Appender interface {
	Name() string
	Layout() *Layout
	Write(*Event) error
	Close() error
}

type logWriter struct {
	appenders map[string]Appender
	lock      sync.RWMutex
}

func newlogWriter() *logWriter {
	return &logWriter{
		appenders: make(map[string]Appender),
	}
}
func (a *logWriter) RebuildAppender(newMap map[string]Appender) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.appenders = newMap
}

//Attach  添加appender
func (a *logWriter) Attach(appender Appender) {
	name := appender.Name()
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, ok := a.appenders[name]; ok {
		panic(fmt.Errorf("重复注册appender:%s", name))
	}
	a.appenders[name] = appender
}
func (a *logWriter) Detach(name string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if appender, ok := a.appenders[name]; ok {
		appender.Close()
	}
	delete(a.appenders, name)
}

//Log 记录日志信息
func (a *logWriter) Log(event *Event) {
	defer func() {
		if err := recover(); err != nil {
			sysLogger.Panicf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, adr := range a.appenders {
		adr.Write(event)
	}
}

//Close 关闭日志
func (a *logWriter) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, v := range a.appenders {
		v.Close()
	}
	return nil
}
