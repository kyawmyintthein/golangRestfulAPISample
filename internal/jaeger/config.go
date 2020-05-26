package jaeger

type JaegerCfg struct {
	Enabled                   bool    `json:"enabled" mapstructure:"enabled"`
	LogSpans                  bool    `json:"log_spans" mapstructure:"log_spans"`
	LocalServiceName          string  `json:"local_service_name" mapstructure:"local_service_name"`
	SamplerType               string  `json:"sampler_type" mapstructure:"sampler_type"`
	SamplerParam              float64 `json:"sampler_param" mapstructure:"sampler_param"`
	SamplingServerURL         string  `json:"sampling_server_url" mapstructure:"sampling_server_url"`
	LocalAgentPort            string  `json:"local_agent_host_port" mapstructure:"local_agent_host_port"`
	ReporterCollectorEndpoint string  `json:"reporter_collector_endpoint" mapstructure:"reporter_collector_endpoint"`
}
