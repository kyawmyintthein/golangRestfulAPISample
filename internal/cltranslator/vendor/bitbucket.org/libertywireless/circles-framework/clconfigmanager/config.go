package clconfigmanager

import (
	"time"
)

type ETCDCfg struct {
	Enabled                   bool          `mapstructure:"enabled" json:"enabled"`
	Endpoints                 []string      `mapstructure:"endpoints" json:"endpoints"`
	Username                  string        `mapstructure:"username" json:"username"`
	Password                  string        `mapstructure:"password" json:"password"`
	DialTimeout               time.Duration `mapstructure:"dial_timeout_sec" json:"dial_timeout_sec"`
	RequestTimeout            time.Duration `mapstructure:"req_timeout_sec" json:"req_timeout_sec"`
	WatcherPath               string        `mapstructure:"watcher_path" json:"watcher_path"`
	SkipErrorOnETCDConnFailed bool          `mapstructure:"skip_error_connection_failed" json:"skip_error_connection_failed"`
	ConfigurationLabels       []string      `mapstructure:"configuration_labels" json:"configuration_labels"`
}

type FileType string

const (
	Yaml string = "yaml"
	Yml  string = "yml"
	JSON string = "json"
)

type CCMCfg struct {
	ETCD                      ETCDCfg `mapstructure:"etcd" json:"etcd"`
	EnabledConfigLocalization bool    `mapstructure:"enabled_config_localization" json:"enabled_config_localization"`
	Config                    Config  `mapstructure:"config" json:"config"`
}

type Config struct {
	AppConfig     AppFileCfg `mapstructure:"app" json:"app"`
	CirclesConfig FileCfg    `mapstructure:"circles_config" json:"circles_config"`
	CirclesLocale FileCfg    `mapstructure:"circles_locale" json:"circles_locale"`
}

type AppFileCfg struct {
	BaseDir      string   `mapstructure:"base_dir" json:"base_dir"`
	AppFilePaths []string `mapstructure:"app_files" json:"app_files"`
}

type FileCfg struct {
	BaseDir           string   `mapstructure:"base_dir" json:"base_dir"`
	FileType          FileType `mapstructure:"file_type" json:"file_type"`
	AppFilePaths      []string `mapstructure:"app_files" json:"app_files"`
	SharedFilePaths   []string `mapstructure:"shared_files" json:"shared_files"`
	AppDirectories    []string `mapstructure:"app_directories" json:"app_directories"`
	SharedDirectories []string `mapstructure:"shared_directories" json:"shared_directories"`
}
