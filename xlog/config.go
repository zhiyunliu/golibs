package xlog

import "log"

var (
	DefaultParam = &Param{
		inited:     false,
		ConfigPath: "../etc/logger.json",
		Layout: Layout{
			LevelName: LevelInfo.FullName(), Path: "../log/%ndate/%level/%hh.log", Content: "[%datetime][%l][%session][%idx] %content",
		},
	}
)

type Param struct {
	inited     bool
	ConfigPath string
	Layout     Layout
}

func Config(p *Param) {
	p.inited = true
	if p.ConfigPath == "" {
		p.ConfigPath = DefaultParam.ConfigPath
	}
	err := reconfigLogWriter(p)
	if err != nil {
		log.Println("reconfigLogWriter.Config:", err)
	}
}
