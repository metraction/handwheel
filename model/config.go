package model

type PrometheusConfig struct {
	URL      string `mapstructure:"url"`
	Interval string `mapstructure:"interval"`
}

type CraneConfig struct {
	Images_whitelist []string `mapstructure:"images_whitelist"`
	RegistryUsername string   `mapstructure:"registry_username"`
	RegistryPassword string   `mapstructure:"registry_password"`
}

type Config struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Crane      CraneConfig      `mapstructure:"crane"`
}
