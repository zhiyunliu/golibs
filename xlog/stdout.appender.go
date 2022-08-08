package xlog

import (
	"sync"
	"time"
)

const Stdout string = "stdout"

//StdoutAppender 标准输出器
type StdoutAppender struct {
	stdWriter     *stdWriter
	cleanInterval uint
	closeChan     chan struct{}
	onceLock      sync.Once
	cleanTicker   *time.Ticker
	layout        *Layout
}

func init() {
	Registry(&stdApderBuilder{})
}

type stdApderBuilder struct {
}

func (b *stdApderBuilder) Name() string {
	return Stdout
}
func (b *stdApderBuilder) DefaultLayout() *Layout {
	return &Layout{LevelName: LevelInfo.Name(), Content: _defaultLayout}
}
func (b *stdApderBuilder) Build(layout *Layout) Appender {
	sa := &StdoutAppender{cleanInterval: 20, closeChan: make(chan struct{})}
	sa.cleanTicker = time.NewTicker(time.Duration(sa.cleanInterval) * time.Millisecond)
	sa.layout = layout
	sa.layout.Init()
	sa.stdWriter = newStdWriter(sa.layout)
	return sa
}

func NewStudoutAppender() Appender {
	sa := &StdoutAppender{cleanInterval: 20, closeChan: make(chan struct{})}
	sa.cleanTicker = time.NewTicker(time.Duration(sa.cleanInterval) * time.Millisecond)
	sa.layout = &Layout{LevelName: LevelInfo.FullName(), Content: _defaultLayout}
	sa.layout.Init()
	sa.stdWriter = newStdWriter(sa.layout)
	return sa
}

func (a *StdoutAppender) Name() string {
	return Stdout
}
func (a *StdoutAppender) Layout() *Layout {
	return a.layout
}

//Write 写入日志
func (f *StdoutAppender) Write(event *Event) (err error) {
	f.stdWriter.Write(event)
	return nil
}

//Close 关闭组件
func (a *StdoutAppender) Close() error {
	a.onceLock.Do(func() {
		a.cleanTicker.Stop()
		if a.stdWriter != nil {
			a.stdWriter.Close()
		}
		close(a.closeChan)
	})
	return nil
}
