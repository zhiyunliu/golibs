package xlog

type options struct {
	sid  string
	name string
	data map[string]string
}

func (o *options) reset() {
	o.sid = ""
	o.name = ""
	o.data = nil
}

type Option func(*options)

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSid(sid string) Option {
	return func(o *options) {
		o.sid = sid
	}
}

func WithField(k, v string) Option {
	return func(o *options) {
		if o.data == nil {
			o.data = map[string]string{}
		}
		o.data[k] = v
	}
}

func WithFields(fileds map[string]string) Option {
	return func(o *options) {
		if o.data == nil {
			o.data = map[string]string{}
		}

		for k, v := range fileds {
			o.data[k] = v
		}
	}
}
