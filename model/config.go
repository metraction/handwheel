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

// HttpServerConfig holds configuration for the HTTP server (e.g., for health probes)
type HttpServerConfig struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Crane      CraneConfig      `mapstructure:"crane"`
	HttpServer HttpServerConfig `mapstructure:"httpServer"`
	CARootPEM  string           `mapstructure:"ca_root_pem"`
	CAFile     string           `mapstructure:"ca_root_file"`
}
