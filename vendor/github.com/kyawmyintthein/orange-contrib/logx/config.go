package logx

type LogCfg struct {
	LogFilePath string `mapstructure:"log_file" json:"log_file"`
	LogLevel    string `mapstructure:"log_level" json:"log_level"`
	LogFormat   string `mapstructure:"log_format" json:"log_format"`
	LogRotation bool   `mapstructure:"log_rotation" json:"log_rotation"`
}
