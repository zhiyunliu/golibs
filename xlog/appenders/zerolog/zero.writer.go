package zerolog

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/zhiyunliu/golibs/xfile"
	"github.com/zhiyunliu/golibs/xlog"
	"github.com/zhiyunliu/zerolog"
)

var _ xlog.EventWriter = &zeroWriter{}

type ExtParams struct {
	Console bool `json:"console"`
}

// writer 文件输出器
type zeroWriter struct {
	zerologger *zerolog.Logger
	lastWrite  time.Time
	layout     *xlog.Layout
	file       io.WriteCloser
	onceLock   sync.Once
	countChan  chan struct{}
	closeChan  chan struct{}
	Level      xlog.Level
}

// newZeroWriter 构建基于文件流的日志输出对象,使用带缓冲区的文件写入，缓存区达到4K或每隔3秒写入一次文件。
func newZeroWriter(path string, layout *xlog.Layout) (fa *zeroWriter, err error) {
	fa = &zeroWriter{
		layout:    layout,
		countChan: make(chan struct{}, 100),
		closeChan: make(chan struct{}),
	}
	fa.file, err = xfile.CreateFile(path)
	if err != nil {
		return
	}
	zerolog.FormattedLevels = map[zerolog.Level]string{
		zerolog.DebugLevel: "d",
		zerolog.InfoLevel:  "i",
		zerolog.WarnLevel:  "w",
		zerolog.ErrorLevel: "e",
		zerolog.PanicLevel: "p",
		zerolog.FatalLevel: "f",
		zerolog.Disabled:   "l",
		zerolog.NoLevel:    "n",
		zerolog.TraceLevel: "t",
	}

	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return zerolog.FormattedLevels[l]
	}

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

	var finalOutput io.Writer = fa.file

	if fa.layout != nil && len(fa.layout.ExtParams) > 0 {
		extparam := ExtParams{}
		_ = json.Unmarshal(fa.layout.ExtParams, &extparam)
		if extparam.Console {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}
			finalOutput = io.MultiWriter(finalOutput, consoleWriter)
		}
	}

	// 美化控制台Writer
	log := zerolog.New(finalOutput)
	fa.zerologger = &log

	return fa, nil
}

// Write 写入日志
func (f *zeroWriter) Write(event *xlog.Event) error {
	if f.Level > event.Level {
		return nil
	}
	zevt := f.getZeroEvent(event.Level)
	if event.Level == xlog.LevelPanic {
		zevt.ResetDone(nil)
	}

	if len(event.Tags) > 0 {
		if cip, ok := event.Tags["cip"]; ok && cip != "" {
			zevt.Str("cip", cip)
		}
		if uid, ok := event.Tags["uid"]; ok && uid != "" {
			zevt.Str("uid", uid)
		}
		if span_id, ok := event.Tags["span_id"]; ok && span_id != "" {
			zevt.Str("span_id", span_id)
		}
		if trace_id, ok := event.Tags["trace_id"]; ok && trace_id != "" {
			zevt.Str("trace_id", trace_id)
		}
	}

	zevt.Time("time", event.LogTime).
		Str("sid", event.Session).
		Int32("seq", event.Idx).
		Msg(event.Content)

	f.lastWrite = time.Now()
	return nil
}
func (f *zeroWriter) GetLastTime() time.Time {
	return f.lastWrite
}

// Close 关闭当前appender
func (f *zeroWriter) Close() {
	if f == nil {
		return
	}
	f.onceLock.Do(func() {
		close(f.closeChan)
		if f.file != nil {
			f.file.Close()
		}
	})
}

func (f *zeroWriter) getZeroEvent(level xlog.Level) *zerolog.Event {
	switch level {
	case xlog.LevelDebug:
		return f.zerologger.Debug()
	case xlog.LevelInfo:
		return f.zerologger.Info()
	case xlog.LevelWarn:
		return f.zerologger.Warn()
	case xlog.LevelError:
		return f.zerologger.Error()
	case xlog.LevelPanic:
		return f.zerologger.Panic()
	}
	return f.zerologger.Info()
}
