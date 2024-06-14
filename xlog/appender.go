package xlog

import (
	"log"
	"sync"
)

type AppenderBuilder interface {
	Name() string
	Build(layout *Layout) Appender
}

// Appender 定义appender接口
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

// Attach  添加appender
func (a *logWriter) Attach(appender Appender) {
	name := appender.Name()
	a.lock.Lock()
	defer a.lock.Unlock()
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

// Log 记录日志信息
func (a *logWriter) Log(event *Event) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, adr := range a.appenders {
		adr.Write(event)
	}
}

// Close 关闭日志
func (a *logWriter) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, v := range a.appenders {
		v.Close()
	}
	return nil
}
