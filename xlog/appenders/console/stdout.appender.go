package console

import (
	"sync"

	"github.com/zhiyunliu/golibs/xlog"
)

const Stdout string = "stdout"

// StdoutAppender 标准输出器
type StdoutAppender struct {
	stdWriter *stdWriter
	onceLock  sync.Once
	layout    *xlog.Layout
}

func init() {
	xlog.RegistryBuilder(&stdApderBuilder{})
}

type stdApderBuilder struct {
}

func (b *stdApderBuilder) Name() string {
	return Stdout
}

func (b *stdApderBuilder) Build(layout *xlog.Layout) xlog.Appender {
	sa := &StdoutAppender{}
	sa.layout = layout
	sa.layout.Init()
	sa.stdWriter = newStdWriter(sa.layout)
	return sa
}

func (a *StdoutAppender) Name() string {
	return Stdout
}
func (a *StdoutAppender) Layout() *xlog.Layout {
	return a.layout
}

// Write 写入日志
func (f *StdoutAppender) Write(event *xlog.Event) (err error) {
	f.stdWriter.Write(event)
	return nil
}

// Close 关闭组件
func (a *StdoutAppender) Close() error {
	a.onceLock.Do(func() {
		if a.stdWriter != nil {
			a.stdWriter.Close()
		}
	})
	return nil
}
