package datetime

const DefaultTimeformat = "2006-01-02 15:04:05"

type Option func(*options)

type options struct {
	Format string
}

// WithFormat
func WithFormat(format string) Option {
	return func(o *options) {
		o.Format = format
	}
}
