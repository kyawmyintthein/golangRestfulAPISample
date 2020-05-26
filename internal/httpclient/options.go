package httpclient

import (
	"bitbucket.org/libertywireless/circles-framework/cljaeger"
	"bitbucket.org/libertywireless/circles-framework/cllogging"
	"bitbucket.org/libertywireless/circles-framework/cloption"
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/newrelic"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/option"
	"time"
)

type httpRequestTimeoutKey struct{}

func WithRequestTimeout(a time.Duration) option.Option {
	return func(o *option.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, httpRequestTimeoutKey{}, a)
	}
}

type newrelicKey struct{}

func WithNewrelic(a newrelic.NewrelicTracer) option.Option {
	return func(o *option.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, newrelicKey{}, a)
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

type jaegerKey struct{}

func WithJaeger(a jaeger.Jaeger) cloption.Option {
	return func(o *option.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, jaegerKey{}, a)
	}
}

type retryConfig struct{}

func WithRetryConfig(a *RetryCfg) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, retryConfig{}, a)
	}
}
