package xlog

import "log"

var (
	_confPath    = "../conf/logger.json"
	_etcPath     = "../etc/logger.json"
	_logfilePath = "../log/%date/%level/%hh.log"

	_defaultParam = &Param{
		inited:     false,
		ConfigPath: _etcPath,
	}
)

type Param struct {
	inited     bool
	ConfigPath string
}

func Config(p *Param) {
	p.inited = true
	_defaultParam.inited = true
	if p.ConfigPath == "" {
		p.ConfigPath = _defaultParam.ConfigPath
	}
	err := reconfigLogWriter(p)
	if err != nil {
		log.Println("reconfigLogWriter.Config:", err)
	}
}
