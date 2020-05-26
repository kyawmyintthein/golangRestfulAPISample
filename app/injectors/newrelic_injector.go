package injectors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/newrelic"
)

func ProvideNewRelic(config *config.GeneralConfig) (newrelic.NewrelicTracer, error) {
	return newrelic.New(&config.NRTracer)
}
