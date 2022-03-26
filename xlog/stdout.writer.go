package xlog

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zhiyunliu/golibs/xlog/log"
)

//stdWriter 控制台输出器
type stdWriter struct {
	writer     *bytes.Buffer
	output     *log.Logger
	lastWrite  time.Time
	layout     *Layout
	interval   time.Duration
	ticker     *time.Ticker
	lock       sync.Mutex
	onceLock   sync.Once
	countChan  chan struct{}
	closeChan  chan struct{}
	Level      Level
	writeCount uint
}

//newwriter 构建基于文件流的日志输出对象,使用带缓冲区的文件写入，缓存区达到4K或每隔3秒写入一次文件。
func newStdWriter(layout *Layout) (fa *stdWriter, err error) {
	fa = &stdWriter{
		layout:    layout,
		interval:  time.Second * 3,
		countChan: make(chan struct{}, 100),
		closeChan: make(chan struct{}),
	}
	fa.writer = bytes.NewBufferString("")
	fa.Level = layout.Level
	fa.ticker = time.NewTicker(fa.interval)

	fa.output = log.New(fa.writer, "", log.Llongcolor)
	fa.output.SetOutputLevel(log.Ldebug)

	go fa.timeFlush()
	return
}

//Write 写入日志
func (f *stdWriter) Write(event *Event) {
	if f.Level > event.Level {
		return
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.writeCount > 10000 {
		f.countChan <- struct{}{}
		f.writeCount = 0
	}

	switch event.Level {
	case LevelDebug:
		f.output.Debug(event.Output)
	case LevelInfo:
		f.output.Info(event.Output)
	case LevelWarn:
		f.output.Warn(event.Output)
	case LevelError:
		f.output.Error(event.Output)
	case LevelFatal:
		f.output.Output("", log.Lfatal, 1, fmt.Sprintln(event.Output))
	}
	f.lastWrite = time.Now()
}

//Close 关闭当前appender
func (f *stdWriter) Close() {
	f.onceLock.Do(func() {
		f.flush()
		close(f.closeChan)
		f.ticker.Stop()
	})
}

//writeTo 定时写入文件
func (f *stdWriter) timeFlush() {
	for {
		select {
		case <-f.closeChan:
			return
		case <-f.ticker.C:
			f.flush()
		case <-f.countChan:
			f.flush()
		}
	}
}
func (f *stdWriter) flush() {
	f.lock.Lock()
	f.writer.WriteTo(os.Stdout)
	f.writer.Reset()
	f.lock.Unlock()
}
