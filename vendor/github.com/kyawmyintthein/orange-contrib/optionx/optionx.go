package optionx

import (
	"context"
)

// Options struct to keep all parameters in context.Context
type Options struct {
	Context context.Context
}

// Option optional parameter as function
type Option func(o *Options)

// NewOptions convert all Option function to Options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
