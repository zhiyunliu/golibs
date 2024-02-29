package xlog

//Logger 基础日志记录接口
type Logger interface {
	Name() string
	SessionID() string
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	Close()
}
