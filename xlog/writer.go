package xlog

import (
	"github.com/zhiyunliu/golibs/xfile"
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

//默认appender写入器
var mainWriter = newlogWriter()

//AddAppender 添加appender
func AddAppender(appender Appender) {
	mainWriter.Attach(appender)
}

//AddLayout 添加日志输出配置
func AddLayout(l ...*Layout) {
	mainWriter.Append(l...)
}

func asyncWrite(event *Event) {
	mainWriter.Log(event)
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
	globalPause = !layouts.Status
	AddLayout(layouts.Layouts...)
}

var LogPath = "../conf/logger.json"

//进行日志配置文件初始化
func init() {
	AddAppender(NewFileAppender())
	AddAppender(NewStudoutAppender())
	loadLayout(LogPath)
}
