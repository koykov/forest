package forest

type options struct {
	size uint64
}

type Option func(o *options)

func WithSize(size uint64) Option {
	return func(o *options) {
		o.size = size
	}
}

func (o *options) apply(options ...Option) {
	for _, fn := range options {
		fn(o)
	}
}
