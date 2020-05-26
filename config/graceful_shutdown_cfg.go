package config

import "time"

type GracefulShutdownCfg struct {
	Enabled bool          `mapstructure:"enabled" json:"enabled"`
	Timeout time.Duration `mapstructure:"timeout" json:"timeout"` // Second
}
