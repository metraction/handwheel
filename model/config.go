package model

type PrometheusConfig struct {
	URL       string `mapstructure:"url"`
	Interval  string `mapstructure:"interval"`
	BatchSize int    `mapstructure:"batch_size"`
}

type Config struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}
