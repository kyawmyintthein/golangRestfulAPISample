package logging

type LogCfg struct {
	LogFilePath   string `mapstructure:"filepath" json:"filepath"`
	LogLevel      string `mapstructure:"level" json:"v"`
	JsonLogFormat bool   `mapstructure:"json_format" json:"json_format"`
	LogRotation   bool   `mapstructure:"rotation" json:"rotation"`
	CallerInfo    bool   `mapstructure:"caller_enable" json:"caller_enable"`
}
