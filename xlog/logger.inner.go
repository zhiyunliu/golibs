package xlog

import "fmt"

//sysLogger 系统日志
var sysLogger = newSysLogger()

type defaultLogger struct {
	appender Appender
	layout   *Layout
}

func newSysLogger() *defaultLogger {
	return &defaultLogger{
		appender: NewStudoutAppender(),
		layout:   &Layout{Layout: "[%datetime.%ms][%l][%session]%content", Level: LevelAll},
	}
}
func (s *defaultLogger) Error(content ...interface{}) {
	evt := NewEvent("sys", LevelError, CreateSession(), fmt.Sprint(content...), nil)
	s.appender.Write(nil, evt)
}
func (s *defaultLogger) Errorf(f string, content ...interface{}) {
	evt := NewEvent("sys", LevelError, CreateSession(), fmt.Sprintf(f, content...), nil)
	s.appender.Write(nil, evt)
}
