package xlog

//Logger 基础日志记录接口
type Logger interface {
	Name() string

	Log(level Level, content ...interface{})
	Logf(level Level, format string, content ...interface{})
}
