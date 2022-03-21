package xlog

type Level int8

const (
	LevelAll   Level = 0
	LevelDebug Level = 1
	LevelInfo  Level = 2
	LevelWarn  Level = 3
	LevelError Level = 4
	LevelFatal Level = 5
	LevelOff   Level = 9
)

var nameMap = map[Level]string{}
var levelMap = map[string]Level{}

func init() {
	nameMap[LevelDebug] = "d"
	nameMap[LevelInfo] = "i"
	nameMap[LevelWarn] = "w"
	nameMap[LevelError] = "e"
	nameMap[LevelFatal] = "f"
	nameMap[LevelAll] = "a"
	nameMap[LevelOff] = "o"

	nameMap[100+LevelDebug] = "debug"
	nameMap[100+LevelInfo] = "info"
	nameMap[100+LevelWarn] = "warn"
	nameMap[100+LevelError] = "error"
	nameMap[100+LevelFatal] = "fatal"
	nameMap[100+LevelAll] = "all"
	nameMap[100+LevelOff] = "off"

	levelMap["d"] = LevelDebug
	levelMap["i"] = LevelInfo
	levelMap["w"] = LevelWarn
	levelMap["e"] = LevelError
	levelMap["f"] = LevelFatal
	levelMap["a"] = LevelAll
	levelMap["o"] = LevelOff

	levelMap["debug"] = LevelDebug
	levelMap["info"] = LevelInfo
	levelMap["warn"] = LevelWarn
	levelMap["error"] = LevelError
	levelMap["fatal"] = LevelFatal
	levelMap["all"] = LevelAll
	levelMap["off"] = LevelOff

}

func (l Level) Name() string {
	return nameMap[l]
}
func (l Level) FullName() string {
	return nameMap[100+l]
}

func TransLevel(n string) Level {
	return levelMap[n]
}
