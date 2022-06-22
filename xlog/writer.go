package xlog

import (
	"github.com/zhiyunliu/golibs/xfile"
	"github.com/zhiyunliu/golibs/xstack"
)

const StackSkip = 5

var (
	LogPath = "../conf/logger.json"
)

type Writer func(content ...interface{})

func (l Writer) Write(p []byte) (n int, err error) {
	l(string(p))
	return len(p), nil
}

func getStack() string {
	return xstack.GetStack(StackSkip)
}

//默认appender写入器
var _mainWriter = newlogWriter()

//AddAppender 添加appender
func AddAppender(appender Appender) {
	_mainWriter.Attach(appender)
}

//RemoveAppender 移除Appender
func RemoveAppender(name string) {
	_mainWriter.Detach(name)
}

//RemoveAppender 移除Appender
func Appenders() []string {
	result := make([]string, len(_mainWriter.appenders))
	idx := 0
	for key := range _mainWriter.appenders {
		result[idx] = key
	}
	return result
}

//AddLayout 添加日志输出配置
func AddLayout(l ...*Layout) {
	_mainWriter.Append(l...)
}

func asyncWrite(event *Event) {
	_mainWriter.Log(event)
}

func loadLayout(path string) {
	if !xfile.Exists(path) {
		err := Encode(path)
		if err != nil {
			sysLogger.Errorf("创建日志配置文件失败 %v", err)
			return
		}
	}

	layouts, err := Decode(path)
	if err != nil {
		sysLogger.Errorf("读取配置文件失败 %v", err)
		return
	}
	_globalPause = !layouts.Status
	AddLayout(layouts.Layouts...)
}

//进行日志配置文件初始化
func defaultAppender() error {
	AddAppender(NewFileAppender())
	AddAppender(NewStudoutAppender())
	loadLayout(LogPath)
	return nil
}

var _ = defaultAppender()
