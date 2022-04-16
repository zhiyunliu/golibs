package xlog

import (
	"fmt"
	"sync"
	"sync/atomic"

	"bytes"
)

//Logger 日志对象
type LoggerWrap struct {
	opts    *options
	isPause bool
	idx     int32
}

var (
	loggerEventChan chan *Event
	loggerPool      *sync.Pool
	closeChan       chan struct{}
	hasClosed       = false

	onceLock      sync.Once
	closeChanLock sync.Once
	adjustLock    sync.Mutex
	writeRoutines []chan struct{}
)

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
	closeChan = make(chan struct{})
	loggerEventChan = make(chan *Event, 20000)
	writeRoutines = make([]chan struct{}, 0)
	adjustmentWriteRoutine(1)
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
	wrapper.idx = 0
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
	logger.idx = 0
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
	if logger.isPause || _globalPause {
		return
	}
	logger.Log(LevelDebug, args...)
}

//Debugf 输出debug日志
func (logger *LoggerWrap) Debugf(format string, args ...interface{}) {
	if logger.isPause || _globalPause {
		return
	}
	logger.Logf(LevelDebug, format, args...)
}

//Info 输出info日志
func (logger *LoggerWrap) Info(args ...interface{}) {
	if logger.isPause || _globalPause {
		return
	}
	logger.Log(LevelInfo, args...)
}

//Infof 输出info日志
func (logger *LoggerWrap) Infof(format string, args ...interface{}) {
	if logger.isPause || _globalPause {
		return
	}
	logger.Logf(LevelInfo, format, args...)
}

//Warn 输出info日志
func (logger *LoggerWrap) Warn(args ...interface{}) {
	if logger.isPause || _globalPause {
		return
	}
	logger.Log(LevelWarn, args...)
}

//Warnf 输出info日志
func (logger *LoggerWrap) Warnf(format string, args ...interface{}) {
	if logger.isPause || _globalPause {
		return
	}
	logger.Logf(LevelWarn, format, args...)
}

//Error 输出Error日志
func (logger *LoggerWrap) Error(args ...interface{}) {
	if logger.isPause || _globalPause {
		return
	}
	logger.Log(LevelError, args...)
}

//Errorf 输出Errorf日志
func (logger *LoggerWrap) Errorf(format string, args ...interface{}) {
	if logger.isPause || _globalPause {
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
	atomic.AddInt32(&logger.idx, 1)
	event.Idx = logger.idx
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
	atomic.AddInt32(&logger.idx, 1)
	event.Idx = logger.idx
	loggerEventChan <- event
}

func loopWriteEvent(item chan struct{}) {
	for {
		select {
		case <-item:
			return
		case v, ok := <-loggerEventChan:
			if !ok {
				closeChanLock.Do(func() {
					close(closeChan)
				})
				return
			}
			asyncWrite(v)
			v.Close()
		}
	}
}

func adjustmentWriteRoutine(cnt int) {
	adjustLock.Lock()
	defer adjustLock.Unlock()

	curCnt := len(writeRoutines)
	if cnt == curCnt {
		return
	}

	if cnt > curCnt {
		for i, adc := 0, cnt-curCnt; i < adc; i++ {
			nwr := make(chan struct{})
			writeRoutines = append(writeRoutines, nwr)
			go loopWriteEvent(nwr)
		}
		return
	}

	if cnt < curCnt {
		newRoutines := writeRoutines[0:cnt]
		overRoutines := writeRoutines[cnt:]
		for _, item := range overRoutines {
			close(item)
		}
		writeRoutines = newRoutines
		return
	}
}

func getString(c ...interface{}) string {
	if len(c) == 1 {
		return fmt.Sprintf("%+v", c[0])
	}
	var buf bytes.Buffer
	for i := 0; i < len(c); i++ {
		buf.WriteString(fmt.Sprintf("%+v", c[i]))
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
