package aes

type Option func(*options)

type options struct {
	IV        []byte
	BlockSize int
}

func WithIV(iv []byte) Option {
	return func(o *options) {
		o.IV = iv
	}
}

func WithBlockSize(size int) Option {
	return func(o *options) {
		o.BlockSize = size
	}
}
