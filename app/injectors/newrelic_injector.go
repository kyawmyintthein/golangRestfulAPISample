package injectors

import (
	"context"

	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/orange-contrib/logx"
	"github.com/kyawmyintthein/orange-contrib/tracingx/newrelicx"
)

func ProvideNewRelic(config *config.GeneralConfig, logger logx.Logger) (newrelicx.NewrelicTracer, error) {
	newrelicTracer, err := newrelicx.New(&config.NRTracer, newrelicx.WithLogger(logger))
	if err != nil {
		logx.Error(context.Background(), err, "failed to init newrelic tracer")
		return newrelicTracer, nil
	}
	return newrelicTracer, nil
}
