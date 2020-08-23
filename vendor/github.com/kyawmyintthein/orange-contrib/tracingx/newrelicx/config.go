package newrelicx

type NewrelicCfg struct {
	Enabled  bool              `json:"enabled" 	mapstructure:"enabled"`
	Name     string            `json:"name" 	mapstructure:"name"`
	License  string            `json:"license"  mapstructure:"license"`
	SkipURLs map[string]string `json:"skip_urls"  mapstructure:"skip_urls"`
}
