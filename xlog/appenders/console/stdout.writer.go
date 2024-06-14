package console

import (
	"bytes"
	"os"
	"sync"
	"time"

	"github.com/zhiyunliu/golibs/xlog"
	"github.com/zhiyunliu/golibs/xlog/appenders/console/log"
)

// stdWriter 控制台输出器
type stdWriter struct {
	writer    *bytes.Buffer
	output    *log.Logger
	lastWrite time.Time
	layout    *xlog.Layout
	mutex     sync.Mutex
	onceLock  sync.Once
	Level     xlog.Level
}

// newwriter 构建基于文件流的日志输出对象,使用带缓冲区的文件写入，缓存区达到4K或每隔3秒写入一次文件。
func newStdWriter(layout *xlog.Layout) (fa *stdWriter) {
	fa = &stdWriter{
		layout: layout,
	}
	fa.onceLock = sync.Once{}
	fa.mutex = sync.Mutex{}
	fa.writer = bytes.NewBufferString("")
	fa.Level = layout.Level
	//fa.ticker = time.NewTicker(fa.interval)

	fa.output = log.New(fa.writer, "", log.Llongcolor)
	fa.output.SetOutputLevel(log.Ldebug)

	//go fa.timeFlush()
	return
}

// Write 写入日志
func (f *stdWriter) Write(event *xlog.Event) {
	if f.Level > event.Level {
		return
	}
	//保持串行化
	f.mutex.Lock()
	defer f.mutex.Unlock()

	event = event.Format(f.layout)
	switch event.Level {
	case xlog.LevelDebug:
		f.output.Debug(event.Output)
	case xlog.LevelInfo:
		f.output.Info(event.Output)
	case xlog.LevelWarn:
		f.output.Warn(event.Output)
	case xlog.LevelError:
		f.output.Error(event.Output)
	case xlog.LevelPanic:
		f.output.Panic(event.Output)
	case xlog.LevelFatal:
		f.output.Output("", log.Lfatal, 1, event.Output)
	}
	f.writer.WriteTo(os.Stdout)
	f.writer.Reset()
	f.lastWrite = time.Now()
}

// Close 关闭当前appender
func (f *stdWriter) Close() {
	//empty func
}
