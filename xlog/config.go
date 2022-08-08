package xlog

var (
	_logPath = "../conf/logger.json"
	_etcPath = "../etc/logger.json"

	_defaultParam = &Param{
		inited:     false,
		ConfigPath: _logPath,
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
	reconfigLogWriter(p)
}
