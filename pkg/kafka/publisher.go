package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

//go:generate mockgen -package mocks -destination mocks/publisher_mocks.go github.com/ThreeDotsLabs/watermill/message Publisher

const (
	defaultConnAttemts = 10
	defaultConnTimeout = time.Second
)

var defaultPublisherConfig = kafka.PublisherConfig{
	Marshaler: kafka.DefaultMarshaler{},
}

// PublisherOption defines the functional option pattern for configuring the Kafka publisher.
type PublisherOption func(*kafka.PublisherConfig) error

// WithMarshaler sets the message marshaller for the Kafka publisher configuration.
func WithMarshaler(marhaler kafka.Marshaler) PublisherOption {
	return func(pc *kafka.PublisherConfig) error {
		pc.Marshaler = marhaler
		return nil
	}
}

// WithPublisherSaramaConfig sets the Sarama configuration for the Kafka publisher.
func WithPublisherSaramaConfig(config *sarama.Config) PublisherOption {
	return func(pc *kafka.PublisherConfig) error {
		pc.OverwriteSaramaConfig = config
		return nil
	}
}

// WithPublisherOTELEnabled sets the OpenTelemetry enabled flag for the Kafka publisher.
func WithPublisherOTELEnabled(otelEnabled bool) PublisherOption {
	return func(pc *kafka.PublisherConfig) error {
		pc.OTELEnabled = otelEnabled
		return nil
	}
}

// WithPublisherTracer sets the tracer for the Kafka publisher configuration.
func WithPublisherTracer(tracer kafka.SaramaTracer) PublisherOption {
	return func(pc *kafka.PublisherConfig) error {
		pc.Tracer = tracer
		return nil
	}
}

// NewPublisher creates a new Kafka publisher with the provided broker addresses and options.
func NewPublisher(brokers []string, opts ...PublisherOption) (message.Publisher, error) {
	publisherConfig := defaultPublisherConfig
	publisherConfig.Brokers = brokers

	for _, opt := range opts {
		if err := opt(&publisherConfig); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	connAttempts := defaultConnAttemts
	connTimeout := defaultConnTimeout

	log, err := logger.NewLogrusLogger("info")
	if err != nil {
		return nil, fmt.Errorf("failed to create new logger: %w", err)
	}

	var publisher message.Publisher

	for connAttempts > 0 {
		publisher, err = kafka.NewPublisher(publisherConfig, defaultLogger)
		if err == nil {
			break
		}

		log.Info("Publisher is trying to connect...", map[string]interface{}{
			"attempts_left": connAttempts,
		})

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create a new publisher: %w", err)
	}

	return publisher, nil
}
