package base_repository

import (
	"bitbucket.org/libertywireless/circles-framework/clnewrelic"
	"bitbucket.org/libertywireless/circles-framework/cloption"
	"context"
)

type newRelicTracerKey struct{}

func WithNRTracer(a clnewrelic.NewrelicTracer) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, newRelicTracerKey{}, a)
	}
}

