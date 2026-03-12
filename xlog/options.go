package xlog

//pipe buffer size
var BufferSize = 20000

type options struct {
	sid     string
	name    string
	srvType string
	data    map[string]any
}

func (o *options) reset() {
	o.sid = ""
	o.name = ""
	o.srvType = ""
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
func WithSrvType(srvType string) Option {
	return func(o *options) {
		o.srvType = srvType
	}
}

func WithField(k, v string) Option {
	return func(o *options) {
		if o.data == nil {
			o.data = map[string]any{}
		}
		o.data[k] = v
	}
}

func WithFields(fileds map[string]string) Option {
	return func(o *options) {
		if o.data == nil {
			o.data = map[string]any{}
		}

		for k, v := range fileds {
			o.data[k] = v
		}
	}
}

// EventOption 事件选项
func WithEventData(k string, v any) EventOption {
	return func(evt *Event) {
		if evt.Tags == nil {
			evt.Tags = map[string]any{}
		}
		evt.Tags[k] = v
	}
}
