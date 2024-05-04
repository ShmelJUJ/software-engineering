package subscriber

import (
	"errors"
	"fmt"

	"dario.cat/mergo"
)

var ErrNilConfig = errors.New("cannot override nil config")

const (
	defaultFailedTransactionTopic    = "transaction.failed"
	defaultSucceededTransactionTopic = "transaction.succeeded"
)

// Config represents the subscriber configuration structure.
type Config struct {
	FailedTransactionTopic    string
	SucceededTransactionTopic string
}

func getDefaultConfig() *Config {
	return &Config{
		FailedTransactionTopic:    defaultFailedTransactionTopic,
		SucceededTransactionTopic: defaultSucceededTransactionTopic,
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
