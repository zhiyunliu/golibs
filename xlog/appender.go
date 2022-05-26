package xlog

import (
	"fmt"
	"sync"
)

//Appender 定义appender接口
type Appender interface {
	Name() string
	Write(*Layout, *Event) error
	Close() error
}

type logWriter struct {
	appenders map[string]Appender
	layouts   []*Layout
	lock      sync.RWMutex
}

func newlogWriter() *logWriter {
	return &logWriter{
		appenders: make(map[string]Appender),
		layouts:   make([]*Layout, 0),
	}
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

//Append 添加layout配置
func (a *logWriter) Append(layouts ...*Layout) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, layout := range layouts {
		if layout.Level == LevelOff {
			continue
		}
		if _, ok := a.appenders[layout.Type]; ok {
			layout.Init()
			a.layouts = append(a.layouts, layout)
		}
	}
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
	for _, layout := range a.layouts {
		if layout.Level > event.Level {
			continue
		}
		if apppender, ok := a.appenders[layout.Type]; ok {
			apppender.Write(layout, event.Format(layout))
		}
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
