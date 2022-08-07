package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Logger           LoggerConfig
	Address          string        `env:"ADDRESS"`
	FileBackuperPath string        `env:"STORE_FILE"`
	ShutdownTimeout  time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"3s"`
	Backup           BackupConfig
	HashKey          string `env:"KEY"`
}

type BackupConfig struct {
	Interval  time.Duration `env:"STORE_INTERVAL"`
	DoRestore bool          `env:"RESTORE"`
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}

	// Parse flag-settable fields
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&cfg.FileBackuperPath, "f", "/tmp/devops-metrics-db.json", "backup file path")
	flag.BoolVar(&cfg.Backup.DoRestore, "r", true, "restore from backup file on server start")
	flag.DurationVar(&cfg.Backup.Interval, "i", 300*time.Second, "backup/store interval")
	flag.StringVar(&cfg.HashKey, "k", "", "Hash key for verifying incoming requests' hash-sums")

	parseLoggerConfigFlags(&cfg.Logger)

	flag.Parse()

	// Parse env-settable fields, override if already set
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
