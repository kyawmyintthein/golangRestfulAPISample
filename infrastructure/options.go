package infrastructure

import (
	"context"

	"github.com/kyawmyintthein/orange-contrib/optionx"
	"github.com/kyawmyintthein/orange-contrib/tracingx/newrelicx"
)

type newrelicTracerKey struct{}

func WithNewrelicTracer(nr newrelicx.NewrelicTracer) optionx.Option {
	return func(o *optionx.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, newrelicTracerKey{}, nr)
	}
}
