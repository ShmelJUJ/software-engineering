package subscriber

import (
	"errors"
	"fmt"

	"dario.cat/mergo"
)

var ErrNilConfig = errors.New("cannot override nil config")

const (
	defaultProcessTopic = "monitor.process"
)

// Config represents the monitor configuration structure.
type Config struct {
	ProcessTopic string
}

func getDefaultConfig() *Config {
	return &Config{
		ProcessTopic: defaultProcessTopic,
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
