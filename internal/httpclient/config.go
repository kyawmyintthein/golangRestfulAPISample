package httpclient

import "time"

type HttpClientCfg struct {
}

type RetryCfg struct {
	Enabled          bool            `json:"enabled"`
	MaxRetryAttempts uint            `json:"max_retry_attempts"`
	BackOffDurations []time.Duration `json:"back_off_durations"`
}
