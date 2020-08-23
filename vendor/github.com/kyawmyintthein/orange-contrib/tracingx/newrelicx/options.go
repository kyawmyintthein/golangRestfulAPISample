package newrelicx

import (
	"context"

	"github.com/kyawmyintthein/orange-contrib/logx"
	"github.com/kyawmyintthein/orange-contrib/optionx"
)

type loggerKey struct{}

func WithLogger(logger logx.Logger) optionx.Option {
	return func(o *optionx.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, loggerKey{}, logger)
	}
}
