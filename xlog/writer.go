package xlog

import (
	"github.com/zhiyunliu/golibs/xstack"
)

const StackSkip = 5

type Writer func(content ...interface{})

func (l Writer) Write(p []byte) (n int, err error) {
	l(string(p))
	return len(p), nil
}

func getStack() string {
	return xstack.GetStack(StackSkip)
}

// 默认appender写入器
var _mainWriter = newlogWriter()

// AppenderList 获取列表
func AppenderList() []string {
	result := make([]string, len(_mainWriter.appenders))
	idx := 0
	for key := range _mainWriter.appenders {
		result[idx] = key
	}
	return result
}

func asyncWrite(event *Event) {
	_mainWriter.Log(event)
}
