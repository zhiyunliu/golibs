package xlog

type Level int8

const (
	_skipLevel = 100
)

const (
	LevelAll   Level = 0
	LevelDebug Level = 1
	LevelInfo  Level = 2
	LevelWarn  Level = 3
	LevelError Level = 4
	LevelPanic Level = 5
	LevelFatal Level = 8
	LevelOff   Level = 9
)

var nameMap, levelMap = initMap()

func initMap() (map[Level]string, map[string]Level) {

	var nameMap = map[Level]string{}
	var levelMap = map[string]Level{}

	nameMap[LevelDebug] = "d"
	nameMap[LevelInfo] = "i"
	nameMap[LevelWarn] = "w"
	nameMap[LevelError] = "e"
	nameMap[LevelPanic] = "p"
	nameMap[LevelFatal] = "f"
	nameMap[LevelAll] = "a"
	nameMap[LevelOff] = "o"

	nameMap[_skipLevel+LevelDebug] = "debug"
	nameMap[_skipLevel+LevelInfo] = "info"
	nameMap[_skipLevel+LevelWarn] = "warn"
	nameMap[_skipLevel+LevelError] = "error"
	nameMap[_skipLevel+LevelPanic] = "panic"
	nameMap[_skipLevel+LevelFatal] = "fatal"
	nameMap[_skipLevel+LevelAll] = "all"
	nameMap[_skipLevel+LevelOff] = "off"

	levelMap["d"] = LevelDebug
	levelMap["i"] = LevelInfo
	levelMap["w"] = LevelWarn
	levelMap["e"] = LevelError
	levelMap["p"] = LevelPanic
	levelMap["f"] = LevelFatal
	levelMap["a"] = LevelAll
	levelMap["o"] = LevelOff

	levelMap["debug"] = LevelDebug
	levelMap["info"] = LevelInfo
	levelMap["warn"] = LevelWarn
	levelMap["error"] = LevelError
	levelMap["panic"] = LevelPanic
	levelMap["fatal"] = LevelFatal
	levelMap["all"] = LevelAll
	levelMap["off"] = LevelOff

	return nameMap, levelMap
}

func (l Level) Name() string {
	return nameMap[l]
}
func (l Level) FullName() string {
	return nameMap[_skipLevel+l]
}

func TransLevel(n string) Level {
	return levelMap[n]
}
