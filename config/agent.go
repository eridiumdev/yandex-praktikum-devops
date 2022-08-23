package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	Logger LoggerConfig

	CollectInterval time.Duration `env:"POLL_INTERVAL"`
	ExportInterval  time.Duration `env:"REPORT_INTERVAL"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"3s"`

	RandomExporter RandomExporterConfig `envPrefix:"RANDOM_EXPORTER_"`
	HTTPExporter   HTTPExporterConfig

	HashKey string `env:"KEY"`
}

type RandomExporterConfig struct {
	Min int `env:"MIN" envDefault:"0"`
	Max int `env:"MAX" envDefault:"9999"`
}

type HTTPExporterConfig struct {
	Address string        `env:"ADDRESS"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"3s"`
}

func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}

	// Parse flag-settable fields
	flag.DurationVar(&cfg.CollectInterval, "p", 2*time.Second, "metrics collect/poll interval")
	flag.DurationVar(&cfg.ExportInterval, "r", 10*time.Second, "metrics export/report interval")
	flag.StringVar(&cfg.HTTPExporter.Address, "a", "localhost:8080", "HTTP exporter target address")
	flag.StringVar(&cfg.HashKey, "k", "", "Hash key for signing metrics data")

	parseLoggerConfigFlags(&cfg.Logger)

	flag.Parse()

	// Parse env-settable fields, override if already set
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
