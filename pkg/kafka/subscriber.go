package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

//go:generate mockgen -package mocks -destination mocks/subscriber_mocks.go github.com/ThreeDotsLabs/watermill/message Subscriber

const (
	defaultNumPartitions     = 1
	defaultReplicationFactor = 2
)

var defaultSubscriberConfig = kafka.SubscriberConfig{
	Unmarshaler:   kafka.DefaultMarshaler{},
	ConsumerGroup: watermill.NewUUID(),
	InitializeTopicDetails: &sarama.TopicDetail{
		NumPartitions:     defaultNumPartitions,
		ReplicationFactor: defaultReplicationFactor,
	},
}

// SubscriberOption represents a functional option for configuring a Kafka subscriber.
type SubscriberOption func(*kafka.SubscriberConfig) error

// WithUnmarshaler sets the unmarshaler for deserializing Kafka messages.
func WithUnmarshaler(unmarshaler kafka.Unmarshaler) SubscriberOption {
	return func(sc *kafka.SubscriberConfig) error {
		sc.Unmarshaler = unmarshaler
		return nil
	}
}

// WithSubscriberSaramaConfig sets the Sarama configuration for the Kafka subscriber.
func WithSubscriberSaramaConfig(config *sarama.Config) SubscriberOption {
	return func(sc *kafka.SubscriberConfig) error {
		sc.OverwriteSaramaConfig = config
		return nil
	}
}

// WithSubscriberConsumerGroup sets the consumer group for the Kafka subscriber.
func WithSubscriberConsumerGroup(consumerGroup string) SubscriberOption {
	return func(sc *kafka.SubscriberConfig) error {
		sc.ConsumerGroup = consumerGroup
		return nil
	}
}

// WithSubscriberInitTopic sets the initial topic details for the Kafka subscriber.
func WithSubscriberInitTopic(initTopic *sarama.TopicDetail) SubscriberOption {
	return func(sc *kafka.SubscriberConfig) error {
		sc.InitializeTopicDetails = initTopic
		return nil
	}
}

// WithSubscriberTracer sets the tracer for tracing Kafka messages.
func WithSubscriberTracer(tracer *kafka.OTELSaramaTracer) SubscriberOption {
	return func(sc *kafka.SubscriberConfig) error {
		sc.Tracer = tracer
		return nil
	}
}

// NewSubscriber creates a new Kafka subscriber with the specified brokers and options.
func NewSubscriber(brokers []string, opts ...SubscriberOption) (message.Subscriber, error) {
	subscriberConfig := defaultSubscriberConfig
	subscriberConfig.Brokers = brokers

	for _, opt := range opts {
		if err := opt(&subscriberConfig); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	connAttempts := defaultConnAttemts
	connTimeout := defaultConnTimeout

	log, err := logger.NewLogrusLogger("info")
	if err != nil {
		return nil, fmt.Errorf("failed to create new logger: %w", err)
	}

	var subscriber message.Subscriber

	for connAttempts > 0 {
		subscriber, err = kafka.NewSubscriber(subscriberConfig, defaultLogger)
		if err == nil {
			break
		}

		log.Info("Subscriber is trying to connect...", map[string]interface{}{
			"attempts_left": connAttempts,
		})

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create a new subscriber: %w", err)
	}

	return subscriber, nil
}
