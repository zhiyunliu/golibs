package xlog

import (
	"fmt"

	"github.com/zhiyunliu/golibs/session"
)

//sysLogger 系统日志
var sysLogger = newSysLogger()

type defaultLogger struct {
	appender Appender
}

func newSysLogger() *defaultLogger {
	return &defaultLogger{
		appender: NewStudoutAppender(),
	}
}
func (s *defaultLogger) Error(content ...interface{}) {
	evt := NewEvent("sys", LevelError, session.Create(), fmt.Sprint(content...), "sys", nil)
	s.appender.Write(evt)
}
func (s *defaultLogger) Errorf(f string, content ...interface{}) {
	evt := NewEvent("sys", LevelError, session.Create(), fmt.Sprintf(f, content...), "sys", nil)
	s.appender.Write(evt)
}

func (s *defaultLogger) Panic(content ...interface{}) {
	evt := NewEvent("sys", LevelPanic, session.Create(), fmt.Sprint(content...), "sys", nil)
	s.appender.Write(evt)
}
func (s *defaultLogger) Panicf(f string, content ...interface{}) {
	evt := NewEvent("sys", LevelPanic, session.Create(), fmt.Sprintf(f, content...), "sys", nil)
	s.appender.Write(evt)
}
