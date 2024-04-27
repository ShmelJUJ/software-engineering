package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type idempotencyConfig struct {
	Name      string `yaml:"name"`
	HeaderKey string `yaml:"header_key"`
}

type middlewareConfig struct {
	IdempotenctCfg *idempotencyConfig `yaml:"idempotency"`
}

type loggerConfig struct {
	Level string `yaml:"level"`
}

type postgresConfig struct {
	Dialect string `yaml:"dialect"`
	URL     string `yaml:"url"`
	PoolMax int    `yaml:"pool_max"`
}

type redisConfig struct {
	URL string `yaml:"url"`
}

type httpConfig struct {
	Port int `yaml:"port"`
}

type publisherConfig struct {
	Brokers                   []string `yaml:"brokers"`
	ProcessedTransactionTopic string   `yaml:"processed_transaction_topic"`
}

type topicDetails struct {
	NumPartitions     int `yaml:"partitions"`
	ReplicationFactor int `yaml:"replication_factor"`
}

type subscriberConfig struct {
	Brokers                   []string      `yaml:"brokers"`
	TopicDetails              *topicDetails `yaml:"topic_details"`
	SucceededTransactionTopic string        `yaml:"succeeded_transaction_topic"`
	FailedTransactionTopic    string        `yaml:"failed_transaction_topic"`
}

// Config represents the overall configuration structure.
type Config struct {
	LoggerCfg     *loggerConfig     `yaml:"logger"`
	HTTPCfg       *httpConfig       `yaml:"http"`
	PostgresCfg   *postgresConfig   `yaml:"postgres"`
	RedisCfg      *redisConfig      `yaml:"redis"`
	MiddlewareCfg *middlewareConfig `yaml:"middleware"`
	PublisherCfg  *publisherConfig  `yaml:"publisher"`
	SubscriberCfg *subscriberConfig `yaml:"subscriber"`
}

// NewConfig initializes a new Config instance by reading from a YAML file.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return cfg, nil
}
