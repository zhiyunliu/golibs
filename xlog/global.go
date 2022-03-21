package xlog

var globalPause bool

//Pause 暂停记录
func Pause() {
	globalPause = true
}

//Resume 恢复记录
func Resume() {
	globalPause = false
}
