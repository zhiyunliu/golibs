package xlog

import (
	"sync"

	"github.com/zhiyunliu/golibs/xfile"
)

var globalPause bool

//Pause 暂停记录
func Pause() {
	globalPause = true
}

//Resume 恢复记录
func Resume() {
	globalPause = false
}

var once sync.Once

var LogPath = "../conf/logger.conf"

func init() {
	once.Do(func() {
		if !xfile.Exists(LogPath) {
			err := Encode(LogPath)
			if err != nil {
				sysLogger.Errorf("创建日志配置文件失败 %v", err)
				return
			}
		}

		layouts, err := Decode(LogPath)
		if err != nil {
			sysLogger.Errorf("读取配置文件失败 %v", err)
			return
		}
		AddLayout(layouts.Layouts...)
	})
}
