package xlog

import (
	"fmt"
	"sync"

	"bytes"
)

//Logger 日志对象
type LoggerWrap struct {
	opts    *options
	isPause bool
}

var (
	loggerEventChan chan *Event
	loggerPool      *sync.Pool
	closeChan       chan struct{}
	onceLock        sync.Once
	hasClosed       = false
)

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
	closeChan = make(chan struct{})
	loggerEventChan = make(chan *Event, 20000)
	go loopWriteEvent()

}

//New 根据一个或多个日志名称构建日志对象，该日志对象具有新的session id系统不会缓存该日志组件
func New(opt ...Option) (logger Logger) {
	wrapper := &LoggerWrap{}
	opts := &options{
		data: map[string]string{},
	}
	for i := range opt {
		opt[i](opts)
	}
	wrapper.opts = opts
	return wrapper
}

//Name 名字
func (logger *LoggerWrap) Name() string {
	return logger.opts.name
}

//Close 关闭当前日志组件
func (logger *LoggerWrap) Close() {
	logger.opts.reset()
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
	return logger.opts.sid
}

//Debug 输出debug日志
func (logger *LoggerWrap) Debug(args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelDebug, args...)
}

//Debugf 输出debug日志
func (logger *LoggerWrap) Debugf(format string, args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(LevelDebug, format, args...)
}

//Info 输出info日志
func (logger *LoggerWrap) Info(args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelInfo, args...)
}

//Infof 输出info日志
func (logger *LoggerWrap) Infof(format string, args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(LevelInfo, format, args...)
}

//Warn 输出info日志
func (logger *LoggerWrap) Warn(args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelWarn, args...)
}

//Warnf 输出info日志
func (logger *LoggerWrap) Warnf(format string, args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(LevelWarn, format, args...)
}

//Error 输出Error日志
func (logger *LoggerWrap) Error(args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Log(LevelError, args...)
}

//Errorf 输出Errorf日志
func (logger *LoggerWrap) Errorf(format string, args ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.Logf(LevelError, format, args...)
}

func (logger *LoggerWrap) Logf(level Level, format string, args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			sysLogger.Errorf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	if hasClosed {
		return
	}
	event := NewEvent(logger.opts.name, level, logger.opts.sid, fmt.Sprintf(format, args...), logger.opts.data)
	loggerEventChan <- event
}
func (logger *LoggerWrap) Log(level Level, args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			sysLogger.Errorf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	if hasClosed {
		return
	}
	event := NewEvent(logger.opts.name, level, logger.opts.sid, getString(args...), logger.opts.data)
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

func GetLogger(opts ...Option) Logger {
	log := loggerPool.Get().(*LoggerWrap)
	for i := range opts {
		opts[i](log.opts)
	}
	return log
}

//Close 关闭所有日志组件
func Close() {
	onceLock.Do(func() {
		close(loggerEventChan)
		<-closeChan
		mainWriter.Close()
	})
}
