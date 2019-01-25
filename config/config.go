package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type (

	MongodbConfig struct{
		Database  string		`mapstructure:"database"`
		Host  string			`mapstructure:"host"`
	}

	SwaggerConfig struct{
		Host string 	`mapstructure:"host"`
		Version string 	`mapstructure:"version"`
		BasePath string `mapstructure:"base_path"`
	}

	LogConfig struct {
		LogFilePath   string  `mapstructure:"log_file"`
		LogLevel      string  `mapstructure:"log_level"`
		JsonLogFormat bool    `mapstructure:"json_log_format"`
		LogRotation   bool    `mapstructure:"log_rotation"`
	}

	GeneralConfig struct{
		Log  	 LogConfig  				`mapstructure:"log"`
		Mongodb  MongodbConfig     			`mapstructure:"mongodb"`
		Swagger  SwaggerConfig              `mapstructure:"swagger"`
	}
)

func Loadconfig(filepaths ...string)  *GeneralConfig {
	if len(filepaths) == 0{
		panic(fmt.Errorf("Empty config file"))
	}

	viper.SetConfigFile(filepaths[0])
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	for _, filepath := range filepaths[1:] {
		func(filepath string){
			f, err := os.Open(filepath)
			if err != nil{
				panic(fmt.Errorf("Fatal error read config file: %s \n", err))
			}
			defer f.Close()
			err = viper.MergeConfig(f)
			if err != nil {
				panic(fmt.Errorf("Fatal error mergeing config file: %s \n", err))
			}
		}(filepath)
	}

	var config GeneralConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error marshal config file: %s \n", err))
	}
	return &config
}
