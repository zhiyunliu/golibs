package xnet

import "github.com/zhiyunliu/golibs/xlog"

type Logger interface {
	Logf(level xlog.Level, format string, args ...any)
}
