package config

import (
	"fmt"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/publisher"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/subscriber"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway/algorand"
	"github.com/ilyakaznacheev/cleanenv"
)

type loggerConfig struct {
	Level string `yaml:"level"`
}

type topicDetails struct {
	NumPartitions     int `yaml:"partitions"`
	ReplicationFactor int `yaml:"replication_factor"`
}

type kafkaSubscriberConfig struct {
	Brokers       []string           `yaml:"brokers"`
	TopicDetails  *topicDetails      `yaml:"topic_details"`
	SubscriberCfg *subscriber.Config `yaml:"subscriber"`
}

type kafkaPublisherConfig struct {
	Brokers      []string          `yaml:"brokers"`
	PublisherCfg *publisher.Config `yaml:"publisher"`
}

// Config represents the application's configuration structure.
type Config struct {
	LoggerCfg          *loggerConfig          `yaml:"logger"`
	KafkaPublisherCfg  *kafkaPublisherConfig  `yaml:"kafka_publisher"`
	KafkaSubscriberCfg *kafkaSubscriberConfig `yaml:"kafka_subscriber"`
	AlgorandCfg        *algorand.Config       `yaml:"algorand"`
}

// NewConfig initializes a new Config instance by reading from a YAML file.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return cfg, nil
}
