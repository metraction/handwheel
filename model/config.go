package model

type PrometheusConfig struct {
	URL      string `mapstructure:"url"`
	Interval string `mapstructure:"interval"`
}

type Config struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}
