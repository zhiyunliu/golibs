package xxml

import (
	"reflect"
)

//Options Options
type Options struct {
	RootElement        string
	DefaultElement     string
	TimeFormat         string
	NumberType         string
	Compress           bool
	customMarshalers   map[string]XxmlMarshaler
	customUnmarshalers map[string]XxmlUnmarshaler
}

//Option Option
type Option func(opt *Options)

//Customize Customize
type XxmlMarshaler interface {
	Marshal() (string, error)
}
type XxmlUnmarshaler interface {
	Unmarshal(string) error
}

type Customize interface {
	XxmlMarshaler
	XxmlUnmarshaler
}

//WithTimeformat set time format
func WithTimeformat(format string) Option {
	return func(opt *Options) {
		opt.TimeFormat = format
	}
}

//WithNumberType set num type.
func WithNumberType(numType string) Option {
	return func(opt *Options) {
		opt.NumberType = numType
	}
}

//WithElement set num type.
func WithElement(ele string) Option {
	return func(opt *Options) {
		opt.DefaultElement = ele
	}
}

//WithCustomType WithCustomType
func WithCustomType(call Customize) Option {
	return func(opt *Options) {
		ctype := reflect.TypeOf(call)
		if ctype.Kind() == reflect.Ptr {
			ctype = ctype.Elem()
		}
		opt.customMarshalers[ctype.Name()] = call
		opt.customUnmarshalers[ctype.Name()] = call
	}
}

//WithCompress WithCompress
func WithCompress() Option {
	return func(opt *Options) {
		opt.Compress = true
	}
}

//WithUncompress WithUncompress
func WithUncompress() Option {
	return func(opt *Options) {
		opt.Compress = false
	}
}

func (o *Options) IsCustomMarshaler(rtype reflect.Type) bool {
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	_, ok := o.customMarshalers[rtype.Name()]
	return ok
}
