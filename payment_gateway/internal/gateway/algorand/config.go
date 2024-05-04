package algorand

import (
	"errors"
	"fmt"
	"time"

	"dario.cat/mergo"
)

var (
	ErrNilConfig = errors.New("cannot override nil config")

	defaultAlgodToken = "algod-token"
)

const (
	defaultAlgodAddress           = "http://localhost:4001"
	defaultConfirmationWaitRounds = 4
	defaultTimeout                = 3 * time.Second
	defaultRetries                = 10
	defaultIsTest                 = false
)

// Config holds configuration settings for Algorand client.
type Config struct {
	AlgodAddress           string        `yaml:"algod_address"`
	AlgodToken             string        `yaml:"algod_token"`
	ConfirmationWaitRounds uint64        `yaml:"confirmation_wait_rounds"`
	Timeout                time.Duration `yaml:"timeout"`
	Retries                int           `yaml:"retries"`
	IsTest                 bool          `yaml:"is_test"`
}

func getDefaultConfig() *Config {
	return &Config{
		AlgodAddress:           defaultAlgodAddress,
		AlgodToken:             defaultAlgodToken,
		ConfirmationWaitRounds: defaultConfirmationWaitRounds,
		Timeout:                defaultTimeout,
		Retries:                defaultRetries,
		IsTest:                 defaultIsTest,
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
