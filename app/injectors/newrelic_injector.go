package injectors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/newrelic"
)

func ProvideNewRelic(config *config.GeneralConfig) (newrelic.NewrelicTracer, error) {
	newrelicTracer, err := newrelic.New(&config.NRTracer)
	if err != nil {
		cllogging.GetLogger().WithError(err).Error("failed to init newrelic tracer")
		return newrelicTracer, nil
	}
	return newrelicTracer, nil
}
