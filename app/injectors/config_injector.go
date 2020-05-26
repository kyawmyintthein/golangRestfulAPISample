package injectors

import (
	"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/spf13/viper"
	"os"
)

func ProvideConfig(filepaths ...string) (*config.GeneralConfig, error) {
	if len(filepaths) == 0 {
		panic(fmt.Errorf("Empty config file"))
	}

	viper.SetConfigFile(filepaths[0])
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	for _, filepath := range filepaths[1:] {
		func(filepath string) {
			f, err := os.Open(filepath)
			if err != nil {
				panic(fmt.Errorf("Fatal error read config file: %s \n", err))
			}
			defer f.Close()
			err = viper.MergeConfig(f)
			if err != nil {
				panic(fmt.Errorf("Fatal error mergeing config file: %s \n", err))
			}
		}(filepath)
	}

	var config config.GeneralConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return &config, fmt.Errorf("Fatal error marshal config file: %s \n", err)
	}
	return &config, nil
}