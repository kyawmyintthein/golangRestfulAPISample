package injectors

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/orange-contrib/logx"
)

func ProvideLogger(config *config.GeneralConfig) logx.Logger {
	return logx.Init(&config.Log)
}
