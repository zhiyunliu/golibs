package xerrs

import "github.com/zhiyunliu/golibs/xtypes"

type Option func(*xErr)

func WithCode(code int) Option {
	return func(e *xErr) {
		e.Code = code
	}
}

func WithData(data xtypes.XMap) Option {
	return func(e *xErr) {
		e.Data = data
	}
}

func WithIgnore() Option {
	return func(e *xErr) {
		e.Ignore = true
	}
}
