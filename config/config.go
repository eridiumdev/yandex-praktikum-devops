package config

import (
	"errors"

	"github.com/caarlos0/env/v6"
)

// Config sources

const (
	FromEnv = iota
)

func loadConfig(cfg interface{}, source int) error {
	switch source {
	case FromEnv:
		return env.Parse(cfg)
	default:
		return errors.New("invalid config source")
	}
}
