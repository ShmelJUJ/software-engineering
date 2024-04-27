package middleware

import (
	"errors"
	"fmt"

	"dario.cat/mergo"
)

var ErrNilConfig = errors.New("cannot override nil config")

const (
	defaultIdempotencyName      = "global"
	defaultIdempotencyHeaderKey = "X-Idempotency-Key"
)

// IdempotencyConfig defines the configuration for idempotency behavior.
type IdempotencyConfig struct {
	Name      string
	HeaderKey string
}

// Config represents the main configuration structure.
type Config struct {
	IdempotencyCfg *IdempotencyConfig
}

func getDefaultConfig() *Config {
	return &Config{
		IdempotencyCfg: &IdempotencyConfig{
			Name:      defaultIdempotencyName,
			HeaderKey: defaultIdempotencyHeaderKey,
		},
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
