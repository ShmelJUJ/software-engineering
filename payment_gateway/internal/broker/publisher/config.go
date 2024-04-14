package publisher

import (
	"errors"
	"fmt"
	"time"

	"dario.cat/mergo"
)

var ErrNilConfig = errors.New("cannot override nil config")

const (
	defaultPaymentProccessingTime    = 30 * time.Second
	defaultSucceededTransactionTopic = "transaction.succeeded"
	defaultFailedTransactionTopic    = "transaction.failed"
)

// Config represents publisher configuration parameters.
type Config struct {
	PaymentProccessingTime    time.Duration `yaml:"payment_processing_time"`
	SucceededTransactionTopic string        `yaml:"succeeded_transaction_topic"`
	FailedTransactionTopic    string        `yaml:"failed_transaction_topic"`
}

func getDefaultConfig() *Config {
	return &Config{
		PaymentProccessingTime:    defaultPaymentProccessingTime,
		SucceededTransactionTopic: defaultSucceededTransactionTopic,
		FailedTransactionTopic:    defaultFailedTransactionTopic,
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
