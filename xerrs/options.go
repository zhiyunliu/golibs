package xerrs

import "github.com/zhiyunliu/golibs/xtypes"

type Option func(*XErr)

func WithCode(code int) Option {
	return func(e *XErr) {
		e.Code = code
	}
}

func WithData(data xtypes.XMap) Option {
	return func(e *XErr) {
		e.Data = data
	}
}

func WithIgnore() Option {
	return func(e *XErr) {
		e.Ignore = true
	}
}
