package config

type SwaggerCfg struct {
	Host     string `mapstructure:"host"`
	Version  string `mapstructure:"version"`
	BasePath string `mapstructure:"base_path"`
}
