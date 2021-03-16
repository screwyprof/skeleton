package ctxtags

type Option func(*options)

// WithFieldExtractor customizes the function for extracting log fields from protobuf messages, for
// unary and server-streamed methods only.
func WithFieldExtractor(f RequestFieldExtractorFunc) Option {
	return func(o *options) {
		o.requestFieldsFunc = f
	}
}

type options struct {
	requestFieldsFunc RequestFieldExtractorFunc
}

func evaluateOptions(opts ...Option) *options {
	options := &options{}
	for _, o := range opts {
		o(options)
	}
	return options
}
