package xxml

type MapKeyMarshaler interface {
	MarshalKey() string
}

type MapKeyUnmarshaler interface {
	UnmarshalKey(val string) error
}
