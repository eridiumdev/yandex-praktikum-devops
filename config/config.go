package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config sources

const (
	FromEnv = iota
)

// DurationSec is a custom representation for time.Duration
// custom type allows for more flexible unmarshalling, see UnmarshalText()
type DurationSec time.Duration

// UnmarshalText treats provided text as amount of seconds for DurationSec
func (d *DurationSec) UnmarshalText(text []byte) error {
	// Normalize string to "<x>s" format. Trim extra 's' if already in this format
	duration, err := time.ParseDuration(fmt.Sprintf("%ss", strings.TrimRight(string(text), "s")))
	if err != nil {
		return err
	}
	*d = DurationSec(duration)
	return nil
}

func loadConfig(cfg interface{}, source int) error {
	switch source {
	case FromEnv:
		return env.Parse(cfg)
	default:
		return errors.New("invalid config source")
	}
}
