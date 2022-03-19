package xlog

type options struct {
	Sid  string
	Data map[string]string
}

type Option func(*options)

func WithSessionId(sid string) Option {
	return func(o *options) {
		o.Sid = sid
	}
}

func WithField(k, v string) Option {
	return func(o *options) {
		if o.Data == nil {
			o.Data = map[string]string{}
		}
		o.Data[k] = v
	}
}

func WithFields(fileds map[string]string) Option {
	return func(o *options) {
		if o.Data == nil {
			o.Data = map[string]string{}
		}

		for k, v := range fileds {
			o.Data[k] = v
		}
	}
}
