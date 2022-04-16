package xlog

var _globalPause bool

//Pause 暂停记录
func Pause() {
	_globalPause = true
}

//Resume 恢复记录
func Resume() {
	_globalPause = false
}

func Concurrency(cnt int) {
	if cnt <= 0 {
		cnt = 1
	}
	adjustmentWriteRoutine(cnt)
}
