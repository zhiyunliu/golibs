package zerolog

import (
	"fmt"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/zhiyunliu/golibs/xlog"
)

const ZeroLog string = "zerolog"

const (
	_clearTimeRange = time.Minute * 1
	_clearInterval  = time.Second * 30
)

// ZeroAppender 文件ZeroAppender
type ZeroAppender struct {
	writers       cmap.ConcurrentMap[string, xlog.EventWriter]
	cleanTicker   *time.Ticker
	cleanInterval time.Duration
	closeChan     chan struct{}
	onceLock      sync.Once
	layout        *xlog.Layout
}

func init() {
	xlog.RegistryBuilder(&zeroApderBuilder{})
}

type zeroApderBuilder struct {
}

func (b *zeroApderBuilder) Name() string {
	return ZeroLog
}

func (b *zeroApderBuilder) Build(layout *xlog.Layout) xlog.Appender {
	a := &ZeroAppender{
		closeChan:     make(chan struct{}),
		writers:       cmap.New[xlog.EventWriter](),
		cleanInterval: _clearInterval,
	}
	a.layout = layout
	a.layout.Init()
	a.cleanTicker = time.NewTicker(a.cleanInterval)
	go a.clean()
	return a
}

func (a *ZeroAppender) Name() string {
	return ZeroLog
}

func (a *ZeroAppender) Layout() *xlog.Layout {
	return a.layout
}

func (a *ZeroAppender) Write(event *xlog.Event) error {
	filePath := event.Transform(a.layout.Path, false)
	var err error
	res := a.writers.Upsert(filePath, nil, func(exists bool, oldval, _ xlog.EventWriter) xlog.EventWriter {
		if exists {
			return oldval
		}
		newval, innerErr := newZeroWriter(filePath, a.layout)
		if innerErr != nil {
			err = innerErr
			return nil
		}
		return newval
	})
	if res == nil {
		return fmt.Errorf("创建ZeroWriter.Path=%s.Error:%+v", filePath, err)
	}

	return res.Write(event)
}

// Close 关闭组件
func (a *ZeroAppender) Close() error {
	a.onceLock.Do(func() {
		close(a.closeChan)
		a.cleanWriters()
	})

	return nil
}

func (a *ZeroAppender) clean() {
	for {
		select {
		case <-a.closeChan:
			return
		case <-a.cleanTicker.C:
			a.cleanWriters()
		}
	}

}

func (a *ZeroAppender) cleanWriters() {
	remvesList := []string{}
	a.writers.IterCb(func(key string, value xlog.EventWriter) {
		lastwrite := value.GetLastTime()
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
		value.Close()
	}
}
