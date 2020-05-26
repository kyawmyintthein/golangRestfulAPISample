package sql

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/option"
)

type loggerKey struct{}

func WithLogger(a logging.KVLogger) option.Option {
	return func(o *option.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, loggerKey{}, a)
	}
}
