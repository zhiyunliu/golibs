package xlog

import (
	"bufio"
	"io"
	"sync"
	"time"

	"github.com/zhiyunliu/golibs/xfile"
)

//writer 文件输出器
type fileWriter struct {
	writer     *bufio.Writer
	lastWrite  time.Time
	layout     *Layout
	interval   time.Duration
	file       io.WriteCloser
	ticker     *time.Ticker
	lock       sync.Mutex
	onceLock   sync.Once
	countChan  chan struct{}
	closeChan  chan struct{}
	Level      Level
	writeCount uint
}

//newwriter 构建基于文件流的日志输出对象,使用带缓冲区的文件写入，缓存区达到4K或每隔3秒写入一次文件。
func newFileWriter(path string, layout *Layout) (fa *fileWriter, err error) {
	fa = &fileWriter{
		layout:    layout,
		interval:  time.Second * 3,
		countChan: make(chan struct{}, 100),
		closeChan: make(chan struct{}),
	}
	fa.file, err = xfile.CreateFile(path)
	if err != nil {
		return
	}
	fa.lock = sync.Mutex{}
	fa.onceLock = sync.Once{}
	fa.Level = layout.Level
	fa.ticker = time.NewTicker(fa.interval)
	fa.writer = bufio.NewWriterSize(fa.file, 4096)
	go fa.timeFlush()
	return
}

//Write 写入日志
func (f *fileWriter) Write(event *Event) {
	if f.Level > event.Level {
		return
	}
	event = event.Format(f.layout)
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.writeCount > 10000 {
		f.countChan <- struct{}{}
		f.writeCount = 0
	}
	f.writer.WriteString(event.Output)
	f.writer.WriteString("\n")
	f.lastWrite = time.Now()
}

//Close 关闭当前appender
func (f *fileWriter) Close() {
	if f == nil {
		return
	}
	f.onceLock.Do(func() {
		f.flush()
		close(f.closeChan)
		f.ticker.Stop()
	})
}

//writeTo 定时写入文件
func (f *fileWriter) timeFlush() {
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
func (f *fileWriter) flush() {
	f.lock.Lock()
	defer f.lock.Unlock()
	if err := f.writer.Flush(); err != nil {
		sysLogger.Error("file.write.err:", err)
	}
}
