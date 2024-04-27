package publisher

import (
	"errors"
	"fmt"

	"dario.cat/mergo"
)

var ErrNilConfig = errors.New("cannot override nil config")

const (
	defaultProcessedTransactionTopic = "transaction.processed"
)

// Config represents the publisher configuration structure.
type Config struct {
	ProcessedTransactionTopic string
}

func getDefaultConfig() *Config {
	return &Config{
		ProcessedTransactionTopic: defaultProcessedTransactionTopic,
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
