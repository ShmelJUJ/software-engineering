package subscriber

import (
	"errors"
	"fmt"
	"time"

	"dario.cat/mergo"
)

var ErrNilConfig = errors.New("cannot override nil config")

const (
	defaultIdleTimeout   = time.Minute
	defaultMinWorkers    = 0
	defaultNumWorkers    = 100
	defaultTasksCapacity = 1000

	defaultProcessedTransactionTopic = "transaction.processed"
	defaultCancelledTransactionTopic = "transaction.cancelled"
)

// PoolConfig holds configuration settings for worker pool.
type PoolConfig struct {
	IdleTimeout   time.Duration `yaml:"idle_timeout"`
	MinWorkers    int           `yaml:"min_workers"`
	NumWorkers    int           `yaml:"num_workers"`
	TasksCapacity int           `yaml:"tasks_capacity"`
}

// Config represents the transaction subscriber configuration.
type Config struct {
	PoolCfg                   *PoolConfig `yaml:"pool"`
	ProcessedTransactionTopic string      `yaml:"processed_transaction_topic"`
	CancelledTransactionTopic string      `yaml:"cancelled_transaction_topic"`
}

func getDefaultConfig() *Config {
	return &Config{
		PoolCfg: &PoolConfig{
			IdleTimeout:   defaultIdleTimeout,
			MinWorkers:    defaultMinWorkers,
			NumWorkers:    defaultNumWorkers,
			TasksCapacity: defaultTasksCapacity,
		},
		ProcessedTransactionTopic: defaultProcessedTransactionTopic,
		CancelledTransactionTopic: defaultCancelledTransactionTopic,
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
