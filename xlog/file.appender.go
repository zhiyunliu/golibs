package xlog

import (
	"fmt"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

const File string = "file"

const (
	_clearTimeRange = time.Minute * 1
	_clearInterval  = time.Second * 30
)

//FileAppender 文件FileAppender
type FileAppender struct {
	writers       cmap.ConcurrentMap
	cleanTicker   *time.Ticker
	cleanInterval time.Duration
	closeChan     chan struct{}
	onceLock      sync.Once
	layout        *Layout
}

func init() {
	Registry(&fileApderBuilder{})
}

type fileApderBuilder struct {
}

func (b *fileApderBuilder) Name() string {
	return File
}
func (b *fileApderBuilder) DefaultLayout() *Layout {
	return &Layout{LevelName: LevelInfo.Name(), Path: _logfilePath, Content: _defaultLayout}
}
func (b *fileApderBuilder) Build(layout *Layout) Appender {
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

func (a *FileAppender) Layout() *Layout {
	return a.layout
}

func (a *FileAppender) Write(event *Event) error {
	filePath := event.Transform(a.layout.Path, false)
	res := a.writers.Upsert(filePath, nil, func(exists bool, oldval, newval interface{}) interface{} {
		if exists {
			return oldval
		}
		writer, err := newFileWriter(filePath, a.layout)
		if err != nil {
			err = fmt.Errorf("创建FileWriter.Path=%s.Error:%+v", filePath, err)
			panic(err)
		}
		return writer

	})
	res.(*fileWriter).Write(event)
	return nil
}

//Close 关闭组件
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
			value.(*fileWriter).Close()
			remvesList = append(remvesList, key)
			return
		}
	})

	for i := range remvesList {
		a.writers.Remove(remvesList[i])
	}
}
