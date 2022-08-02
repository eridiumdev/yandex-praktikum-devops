package config

import "flag"

type LoggerConfig struct {
	Level string
	Mode  string
}

func parseLoggerConfigFlags(cfg *LoggerConfig) {
	flag.StringVar(&cfg.Level, "log-level", "info",
		"logger verbosity level, crit | error | info | debug")
	flag.StringVar(&cfg.Mode, "log-mode", "dev",
		"logger mode, possible values:\n- dev, for colored unstructured output\n- prod, for colorless JSON output\n")
}
