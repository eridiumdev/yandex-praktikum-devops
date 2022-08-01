package config

type AgentConfig struct {
	CollectInterval DurationSec `env:"POLL_INTERVAL" envDefault:"2"`
	ExportInterval  DurationSec `env:"REPORT_INTERVAL" envDefault:"10"`
	ShutdownTimeout DurationSec `env:"SHUTDOWN_TIMEOUT" envDefault:"3"`

	RandomExporter RandomExporterConfig `envPrefix:"RANDOM_EXPORTER_"`
	HTTPExporter   HTTPExporterConfig
}

type RandomExporterConfig struct {
	Min int `env:"MIN" envDefault:"0"`
	Max int `env:"MAX" envDefault:"9999"`
}

type HTTPExporterConfig struct {
	Address string      `env:"ADDRESS" envDefault:"localhost:8080"`
	Timeout DurationSec `env:"TIMEOUT" envDefault:"3"`
}

func LoadAgentConfig(source int) (*AgentConfig, error) {
	cfg := &AgentConfig{}
	return cfg, loadConfig(cfg, source)
}
