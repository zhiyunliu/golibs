package xlog

//Logger 基础日志记录接口
type Logger interface {
	Infof(format string, content ...interface{})
	Info(content ...interface{})

	Errorf(format string, content ...interface{})
	Error(content ...interface{})

	Debugf(format string, content ...interface{})
	Debug(content ...interface{})

	Fatalf(format string, content ...interface{})
	Fatal(content ...interface{})

	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
}
