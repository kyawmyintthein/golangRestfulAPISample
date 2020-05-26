package jaeger

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/option"
	"github.com/uber/jaeger-client-go"
)

type jaegerLoggerKey struct{}

func WithJaegerLogger(a jaeger.Logger) option.Option {
	return func(o *option.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, jaegerLoggerKey{}, a)
	}
}

type loggerKey struct{}

func WithLogger(a logging.KVLogger) option.Option {
	return func(o *option.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, loggerKey{}, a)
	}
}
