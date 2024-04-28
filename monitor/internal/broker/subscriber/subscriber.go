package subscriber

import (
	"context"
	"errors"

	"github.com/ShmelJUJ/software-engineering/monitor/internal/broker/publisher"
	"github.com/ShmelJUJ/software-engineering/monitor/internal/broker/subscriber/dto"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ThreeDotsLabs/watermill/message"
)

var ErrFailedVerification = errors.New("failed verification")

const (
	fromTransaction    = "transaction"
	fromPaymentGateway = "payment_gateway"

	toProcessedTransactionTopic = "transaction.processed"
	toSucceededTransactionTopic = "transaction.succeeded"
	toFailedTransactionTopic    = "transaction.failed"
)

// MonitorSubscriber represents a subscriber for monitoring.
type MonitorSubscriber struct {
	cfg        *Config
	log        logger.Logger
	sub        message.Subscriber
	router     *message.Router
	monitorPub publisher.MonitorPublisher
}

// NewMonitorSubscriber creates a new MonitorSubscriber instance.
func NewMonitorSubscriber(
	cfg *Config,
	log logger.Logger,
	sub message.Subscriber,
	router *message.Router,
	monitorPub publisher.MonitorPublisher,
) (*MonitorSubscriber, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, NewMonitorSubscriberError("failed to set default config", err)
	}

	return &MonitorSubscriber{
		cfg:        cfg,
		log:        log,
		sub:        sub,
		router:     router,
		monitorPub: monitorPub,
	}, nil
}

// RegisterProcessHandler registers the process handler for the monitor subscriber.
func (s *MonitorSubscriber) RegisterProcessHandler() {
	s.log.Debug("Register process handler", map[string]interface{}{})

	s.router.AddNoPublisherHandler(
		"process",
		s.cfg.ProcessTopic,
		s.sub,
		s.handleProcess,
	)
}

func (s *MonitorSubscriber) handleProcess(msg *message.Message) error {
	processDTO := &dto.Process{}
	if err := processDTO.Decode(msg.Payload); err != nil {
		s.log.Error("failed to decode process dto", map[string]interface{}{
			"error": err,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	s.log.Debug("Start handle process", map[string]interface{}{
		"from":     processDTO.From,
		"to_topic": processDTO.ToTopic,
		"payload":  processDTO.Payload,
	})

	if !verify(processDTO.From, processDTO.ToTopic) {
		s.log.Error("failed verification", map[string]interface{}{
			"error":    ErrFailedVerification,
			"from":     processDTO.From,
			"to_topic": processDTO.ToTopic,
		})

		return nil
	}

	if err := s.monitorPub.PublishProcess(processDTO.ToTopic, processDTO.Payload); err != nil {
		s.log.Error("failed to publish process", map[string]interface{}{
			"error":    err,
			"from":     processDTO.From,
			"to_topic": processDTO.ToTopic,
			"payload":  processDTO.Payload,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	return nil
}

func verify(from, toTopic string) bool {
	if from == fromTransaction && toTopic == toProcessedTransactionTopic {
		return true
	} else if from == fromPaymentGateway && (toTopic == toFailedTransactionTopic || toTopic == toSucceededTransactionTopic) {
		return true
	}

	return false
}

// Run starts the monitor subscriber's router.
func (s *MonitorSubscriber) Run(ctx context.Context) error {
	s.log.Debug("Run monitor subscriber", map[string]interface{}{})

	return s.router.Run(ctx)
}
