package xlog

import (
	"sync"
)

var (
	DefaultParam = &ConfigParam{
		ConfigPath: "../etc/logger.json",
		Layout: &Layout{
			LevelName: LevelInfo.FullName(), Path: "../log/%ndate/%level/%hh.log", Content: "[%datetime][%l][%session][%idx] %content",
		},
	}

	_cfglocker = sync.Mutex{}
)

type ConfigParam struct {
	ConfigPath  string  `json:"config_path"`
	Concurrency int     `json:"concurrency"`
	Layout      *Layout `json:"layout"`
}

func Config(opts ...ConfigOption) (err error) {
	var cfgParam = &ConfigParam{}
	for i := range opts {
		opts[i](cfgParam)
	}
	if cfgParam.ConfigPath == "" {
		cfgParam.ConfigPath = DefaultParam.ConfigPath
	}
	if cfgParam.Layout == nil {
		cfgParam.Layout = DefaultParam.Layout
	}
	adjustmentWriteRoutine(cfgParam.Concurrency)

	err = reconfigLog(cfgParam)
	if err != nil {
		return
	}
	return err
}

func reconfigLog(param *ConfigParam) (err error) {
	_cfglocker.Lock()
	defer _cfglocker.Unlock()

	setting, err := loadLayout(param.ConfigPath)
	if err != nil {
		return err
	}
	_globalPause = !setting.Enable

	for apn, layout := range setting.Layout {
		if tmp, ok := _appenderCache.Load(apn); ok {
			_mainWriter.Attach(tmp.(AppenderBuilder).Build(layout))
		}
	}
	return nil
}

// ConfigOption is a function that sets a configuration option.
type ConfigOption func(*ConfigParam)

// WithConfigPath sets the path of the configuration file.
// If not set, the default path "../etc/logger.json" will be used.
func WithConfigPath(path string) ConfigOption {
	return func(p *ConfigParam) {
		p.ConfigPath = path
	}
}

// WithLayout sets the layout of the log.
func WithLayout(layout *Layout) ConfigOption {
	return func(p *ConfigParam) {
		p.Layout = layout
	}
}

// WithConcurrency sets the concurrency of the log.

func WithConcurrency(concurrency int) ConfigOption {
	return func(p *ConfigParam) {
		p.Concurrency = concurrency
	}
}
