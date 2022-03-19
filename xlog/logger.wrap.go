package xlog

import (
	"fmt"
	"sync"

	"bytes"
)

//Logger 日志对象
type LoggerWrap struct {
	name        string
	sessions    string
	defaultData map[string]string
	isPause     bool
}

var loggerEventChan chan *Event
var loggerPool *sync.Pool
var closeChan chan struct{}
var onceLock sync.Once
var hasClosed = false

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New("")
		},
	}
	closeChan = make(chan struct{})
	loggerEventChan = make(chan *Event, 20000)
	go loopWriteEvent()

}

//New 根据一个或多个日志名称构建日志对象，该日志对象具有新的session id系统不会缓存该日志组件
func New(name string, opt ...Option) (logger Logger) {
	wrapper := &LoggerWrap{}
	wrapper.name = name
	opts := &options{
		Data: map[string]string{},
	}
	for i := range opt {
		opt[i](opts)
	}
	if opts.Sid == "" {
		wrapper.sessions = CreateSession()
	}
	wrapper.defaultData = opts.Data
	return logger
}

//Name 名字
func (logger *LoggerWrap) Name() string {
	return logger.name
}

//Close 关闭当前日志组件
func (logger *LoggerWrap) Close() {
	loggerPool.Put(logger)
}

//Pause 暂停记录
func (logger *LoggerWrap) Pause() {
	logger.isPause = true
}

//Resume 恢复记录
func (logger *LoggerWrap) Resume() {
	logger.isPause = false
}

//GetSessionID 获取当前日志的session id
func (logger *LoggerWrap) GetSessionID() string {
	return logger.sessions
}

//Debug 输出debug日志
func (logger *LoggerWrap) Debug(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelDebug, content...)
}

//Debugf 输出debug日志
func (logger *LoggerWrap) Debugf(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(format, LevelDebug, content...)
}

//Info 输出info日志
func (logger *LoggerWrap) Info(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelInfo, content...)
}

//Infof 输出info日志
func (logger *LoggerWrap) Infof(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(format, LevelInfo, content...)
}

//Warn 输出info日志
func (logger *LoggerWrap) Warn(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelWarn, content...)
}

//Warnf 输出info日志
func (logger *LoggerWrap) Warnf(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(format, LevelWarn, content...)
}

//Error 输出Error日志
func (logger *LoggerWrap) Error(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelError, content...)
}

//Errorf 输出Errorf日志
func (logger *LoggerWrap) Errorf(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(format, LevelError, content...)
}

func (logger *LoggerWrap) Logf(format string, level Level, content ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			sysLogger.Errorf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	if hasClosed {
		return
	}
	event := NewEvent(logger.name, level, logger.sessions, fmt.Sprintf(format, content...), logger.defaultData)
	loggerEventChan <- event
}
func (logger *LoggerWrap) Log(level Level, content ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			sysLogger.Errorf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	if hasClosed {
		return
	}
	event := NewEvent(logger.name, level, logger.sessions, getString(content...), logger.defaultData)
	loggerEventChan <- event
}

func loopWriteEvent() {
	for v := range loggerEventChan {
		asyncWrite(v)
		v.Close()
	}
	close(closeChan)
}
func getString(c ...interface{}) string {
	if len(c) == 1 {
		return fmt.Sprintf("%v", c[0])
	}
	var buf bytes.Buffer
	for i := 0; i < len(c); i++ {
		buf.WriteString(fmt.Sprint(c[i]))
		if i != len(c)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

//Close 关闭所有日志组件
func Close() {
	onceLock.Do(func() {
		close(loggerEventChan)
		<-closeChan
		mainWriter.Close()
	})
}
