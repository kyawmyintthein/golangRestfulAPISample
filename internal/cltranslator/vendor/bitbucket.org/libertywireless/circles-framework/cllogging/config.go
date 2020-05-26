package cllogging

type LogCfg struct {
	LogFilePath   string `mapstructure:"log_file" json:"log_file"`
	LogLevel      string `mapstructure:"log_level" json:"log_level"`
	JsonLogFormat bool   `mapstructure:"json_log_format" json:"json_log_format"`
	LogRotation   bool   `mapstructure:"log_rotation" json:"log_rotation"`
	CallerInfo    bool   `mapstructure:"caller_enable" json:"caller_enable"`
}
