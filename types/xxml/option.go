package xxml

import "reflect"

//Options Options
type Options struct {
	RootElement    string
	DefaultElement string
	TimeFormat     string
	NumberType     string
	customTypes    map[reflect.Type]Customize
}

//Option Option
type Option func(opt *Options)

//Customize Customize
type Customize interface {
	Marshal(obj interface{}) (string, error)
	Unmarshal(string, obj interface{}) error
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
func WithCustomType(ctype reflect.Type, call Customize) Option {
	return func(opt *Options) {
		opt.customTypes[ctype] = call
	}
}

func (o *Options) isCustomType(rtype reflect.Type) bool {
	_, ok := o.customTypes[rtype]
	return ok
}
