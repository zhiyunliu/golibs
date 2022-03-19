package xlog

import (
	"sync"
	"time"
)

const Stdout string = "stdout"

//StdoutAppender 标准输出器
type StdoutAppender struct {
	writerLock    sync.Mutex
	stdWriter     *stdWriter
	cleanInterval uint
	closeChan     chan struct{}
	onceLock      sync.Once
	cleanTicker   *time.Ticker
}

//NewStudoutAppender 构建基于文件流的日志输出对象
func NewStudoutAppender() (sa *StdoutAppender) {
	sa = &StdoutAppender{cleanInterval: 200, closeChan: make(chan struct{})}
	sa.cleanTicker = time.NewTicker(time.Duration(sa.cleanInterval) * time.Millisecond)
	return
}

func (a *StdoutAppender) Name() string {
	return Stdout
}

//Write 写入日志
func (f *StdoutAppender) Write(layout *Layout, event *Event) (err error) {
	if f.stdWriter == nil {
		f.writerLock.Lock()
		if f.stdWriter == nil {
			f.stdWriter, err = newStdWriter(layout)
		}
		f.writerLock.Unlock()
	}
	f.stdWriter.Write(event)
	return nil
}

//Close 关闭组件
func (a *StdoutAppender) Close() error {
	a.onceLock.Do(func() {
		a.cleanTicker.Stop()
		a.stdWriter.Close()
		close(a.closeChan)
	})
	return nil
}
