package config

import "time"

type ServerConfig struct {
	Address         string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"3s"`
}

func LoadServerConfig(source int) (*ServerConfig, error) {
	cfg := &ServerConfig{}
	return cfg, loadConfig(cfg, source)
}
