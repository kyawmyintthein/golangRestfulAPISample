package newrelic

type NewrelicCfg struct {
	Enabled  bool              `json:"enabled" 	mapstructure:"enabled"`
	Name     string            `json:"name" 	mapstructure:"name"`
	License  string            `json:"license"  mapstructure:"license"`
	SkipUrls map[string]string `json:"skip_urls"  mapstructure:"skip_urls"`
}
