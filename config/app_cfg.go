package config

type AppCfg struct {
	Environment string `mapstructure:"environment" json:"environment"`
	HttpPort    int    `mapstructure:"http_port" json:"http_port"`
}
