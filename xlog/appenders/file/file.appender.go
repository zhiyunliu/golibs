package xlog

import (
	"fmt"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/zhiyunliu/golibs/xlog"
)

const File string = "file"

const (
	_clearTimeRange = time.Minute * 1
	_clearInterval  = time.Second * 30
)

// FileAppender 文件FileAppender
type FileAppender struct {
	writers       cmap.ConcurrentMap
	cleanTicker   *time.Ticker
	cleanInterval time.Duration
	closeChan     chan struct{}
	onceLock      sync.Once
	layout        *xlog.Layout
}

func init() {
	xlog.RegistryBuilder(&fileApderBuilder{})
}

type fileApderBuilder struct {
}

func (b *fileApderBuilder) Name() string {
	return File
}

func (b *fileApderBuilder) Build(layout *xlog.Layout) xlog.Appender {
	a := &FileAppender{
		closeChan:     make(chan struct{}),
		writers:       cmap.New(),
		cleanInterval: _clearInterval,
	}
	a.layout = layout
	a.layout.Init()
	a.cleanTicker = time.NewTicker(a.cleanInterval)
	go a.clean()
	return a
}

func (a *FileAppender) Name() string {
	return File
}

func (a *FileAppender) Layout() *xlog.Layout {
	return a.layout
}

func (a *FileAppender) Write(event *xlog.Event) error {
	filePath := event.Transform(a.layout.Path, false)
	res := a.writers.Upsert(filePath, nil, func(exists bool, oldval, newval interface{}) interface{} {
		if exists {
			return oldval
		}
		writer, err := newFileWriter(filePath, a.layout)
		if err != nil {
			return fmt.Errorf("创建FileWriter.Path=%s.Error:%+v", filePath, err)
		}
		return writer
	})
	if err, ok := res.(error); ok {
		return err
	}

	res.(*fileWriter).Write(event)
	return nil
}

// Close 关闭组件
func (a *FileAppender) Close() error {
	a.onceLock.Do(func() {
		close(a.closeChan)
		a.cleanWriters()
	})

	return nil
}

func (a *FileAppender) clean() {
	for {
		select {
		case <-a.closeChan:
			return
		case <-a.cleanTicker.C:
			a.cleanWriters()
		}
	}

}

func (a *FileAppender) cleanWriters() {
	remvesList := []string{}
	a.writers.IterCb(func(key string, value interface{}) {
		lastwrite := value.(*fileWriter).lastWrite
		if time.Since(lastwrite) >= _clearTimeRange {
			remvesList = append(remvesList, key)
			return
		}
	})

	for i := range remvesList {
		value, ok := a.writers.Get(remvesList[i])
		if !ok {
			continue
		}
		a.writers.Remove(remvesList[i])
		value.(*fileWriter).Close()
	}
}
