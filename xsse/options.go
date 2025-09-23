package xsse

type DecoderOption func(*decoder)

func WithDecoderUnmarshal(callback UnmarshalCallback) DecoderOption {
	return func(d *decoder) {
		d.unmarshalCallback = callback
	}
}
