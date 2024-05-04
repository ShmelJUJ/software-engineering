package redis

import (
	"errors"
	"fmt"
	"time"

	"dario.cat/mergo"
)

var ErrNilConfig = errors.New("cannot override nil config")

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

// Config represents the redis configuration structure.
type Config struct {
	ConnURL      string
	MaxPoolSize  int
	ConnAttempts int
	ConnTimeout  time.Duration
}

func getDefaultConfig() *Config {
	return &Config{
		MaxPoolSize:  defaultMaxPoolSize,
		ConnAttempts: defaultConnAttempts,
		ConnTimeout:  defaultConnTimeout,
	}
}

func mergeWithDefault(cfg *Config) (*Config, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	defaultCfg := getDefaultConfig()

	if err := mergo.Merge(defaultCfg, cfg, mergo.WithOverride); err != nil {
		return nil, fmt.Errorf("failed to merge configs: %w", err)
	}

	return defaultCfg, nil
}
