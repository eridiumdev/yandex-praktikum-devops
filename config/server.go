package config

type ServerConfig struct {
	Address          string      `env:"ADDRESS" envDefault:"localhost:8080"`
	BackuperFilePath string      `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	ShutdownTimeout  DurationSec `env:"SHUTDOWN_TIMEOUT" envDefault:"3"`
	Backup           BackupConfig
}

type BackupConfig struct {
	Interval  DurationSec `env:"STORE_INTERVAL" envDefault:"300"`
	DoRestore bool        `env:"RESTORE" envDefault:"true"`
}

func LoadServerConfig(source int) (*ServerConfig, error) {
	cfg := &ServerConfig{}
	return cfg, loadConfig(cfg, source)
}
