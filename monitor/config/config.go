package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type loggerConfig struct {
	Level string `yaml:"level"`
}

type publisherConfig struct {
	Brokers []string `yaml:"brokers"`
}

type topicDetails struct {
	NumPartitions     int `yaml:"partitions"`
	ReplicationFactor int `yaml:"replication_factor"`
}

type subscriberConfig struct {
	Brokers      []string      `yaml:"brokers"`
	TopicDetails *topicDetails `yaml:"topic_details"`
	ProcessTopic string        `yaml:"process_topic"`
}

type httpConfig struct {
	Port int `yaml:"port"`
}

// Config represents the overall configuration structure.
type Config struct {
	LoggerCfg     *loggerConfig     `yaml:"logger"`
	PublisherCfg  *publisherConfig  `yaml:"publisher"`
	SubscriberCfg *subscriberConfig `yaml:"subscriber"`
	HTTPCfg       *httpConfig       `yaml:"http"`
}

// NewConfig initializes a new Config instance by reading from a YAML file.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return cfg, nil
}
