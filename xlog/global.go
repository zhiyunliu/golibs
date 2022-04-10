package xlog

var globalPause bool
var logconcurrency int = 1

//Pause 暂停记录
func Pause() {
	globalPause = true
}

//Resume 恢复记录
func Resume() {
	globalPause = false
}

func Concurrency(cnt int) {
	if cnt <= 0 {
		cnt = 1
	}
	logconcurrency = cnt
	adjustmentWriteRoutine()
}
