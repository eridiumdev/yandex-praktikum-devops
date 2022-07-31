package config

import (
	"time"
)

type AgentConfig struct {
	CollectInterval time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ExportInterval  time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"3s"`

	RandomExporter RandomExporterConfig `envPrefix:"RANDOM_EXPORTER_"`
	HTTPExporter   HTTPExporterConfig
}

type RandomExporterConfig struct {
	Min int `env:"MIN" envDefault:"0"`
	Max int `env:"MAX" envDefault:"9999"`
}

type HTTPExporterConfig struct {
	Address string        `env:"ADDRESS" envDefault:"localhost:8080"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"3s"`
}

func LoadAgentConfig(source int) (*AgentConfig, error) {
	cfg := &AgentConfig{}
	return cfg, loadConfig(cfg, source)
}
