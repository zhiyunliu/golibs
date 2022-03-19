package xlog

import (
	"fmt"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

const File string = "file"

//FileAppender 文件FileAppender
type FileAppender struct {
	writers       cmap.ConcurrentMap
	cleanTicker   *time.Ticker
	cleanInterval uint
	closeChan     chan struct{}
	onceLock      sync.Once
}

//NewFileAppender 构建file FileAppender
func NewFileAppender() *FileAppender {
	a := &FileAppender{
		closeChan:     make(chan struct{}),
		writers:       cmap.New(),
		cleanInterval: 60 * 10,
	}
	a.cleanTicker = time.NewTicker(time.Second * time.Duration(a.cleanInterval))
	go a.clean()
	return a
}

func (a *FileAppender) Name() string {
	return File
}

func (a *FileAppender) Write(layout *Layout, event *Event) error {
	filePath := event.Transform(layout.Path, false)
	res := a.writers.Upsert(filePath, nil, func(exists bool, oldval, newval interface{}) interface{} {
		if exists {
			return oldval
		}
		writer, err := newFileWriter(filePath, layout)
		if err != nil {
			err = fmt.Errorf("创建FileWriter.Path=%s.Error:%+v", filePath, err)
			panic(err)
		}
		return writer

	})
	event = event.Format(layout)
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
	for item := range a.writers.IterBuffered() {
		a.writers.RemoveCb(item.Key, func(key string, value interface{}, exists bool) bool {
			if !exists {
				return exists
			}
			w := value.(*fileWriter)
			if time.Since(w.lastWrite) < 5*time.Minute {
				w.Close()
				return true
			}
			return false
		})
	}
}
