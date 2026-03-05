package xlog

import "time"

//Logger 基础日志记录接口
type Logger interface {
	Name() string
	SessionID() string
	Log(level Level, args ...any)
	Logf(level Level, format string, args ...any)
	LogChain(level Level, msg string, opts ...EventOption)
	Close()
}

type EventOption func(evt *Event)

type EventWriter interface {
	Write(evt *Event) error
	GetLastTime() time.Time
	Close()
}
