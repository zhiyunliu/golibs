package xlog

import (
	"log"
	"sync"

	"github.com/zhiyunliu/golibs/xstack"
)

const StackSkip = 5

var (
	_cfglocker = sync.Mutex{}
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

//AppenderList 获取列表
func AppenderList() []string {
	result := make([]string, len(_mainWriter.appenders))
	idx := 0
	for key := range _mainWriter.appenders {
		result[idx] = key
	}
	return result
}

func asyncWrite(event *Event) {
	if !_defaultParam.inited {
		err := reconfigLogWriter(_defaultParam)
		if err != nil {
			log.Println("reconfigLogWriter.asyncWrite:", err)
		}
	}
	_mainWriter.Log(event)
}

func reconfigLogWriter(param *Param) error {
	_cfglocker.Lock()
	defer _cfglocker.Unlock()
	if param.inited {
		return nil
	}
	param.inited = true

	layoutSetting, err := loadLayout(param.ConfigPath, _etcPath)
	if err != nil {
		return err
	}
	newAppenderMap := make(map[string]Appender)
	for apn, layout := range layoutSetting.Layout {
		if tmp, ok := _appenderCache.Load(apn); ok {
			newAppenderMap[apn] = tmp.(AppenderBuilder).Build(layout)
		}
	}

	_mainWriter.RebuildAppender(newAppenderMap)
	return nil
}
