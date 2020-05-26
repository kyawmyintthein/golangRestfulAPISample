package mongo

import (
	"bitbucket.org/libertywireless/circles-framework/cllogging"
	"bitbucket.org/libertywireless/circles-framework/cloption"
	"context"
)

type loggerKey struct{}

func WithLogger(a cllogging.KVLogger) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, loggerKey{}, a)
	}
}
