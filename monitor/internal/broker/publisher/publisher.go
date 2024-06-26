package publisher

import (
	"encoding/json"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

//go:generate mockgen -package mocks -destination mocks/monitor_publisher_mocks.go github.com/ShmelJUJ/software-engineering/monitor/internal/broker/publisher MonitorPublisher

// MonitorPublisher defines the interface for publishing process messages to a topic.
type MonitorPublisher interface {
	PublishProcess(toTopic string, payload any) error
}

type monitorPublisher struct {
	log logger.Logger
	pub message.Publisher
}

// NewMonitorPublisher creates a new instance of MonitorPublisher.
func NewMonitorPublisher(
	log logger.Logger,
	pub message.Publisher,
) MonitorPublisher {
	return &monitorPublisher{
		log: log,
		pub: pub,
	}
}

// PublishProcess publishes a process message to the specified topic.
func (p *monitorPublisher) PublishProcess(toTopic string, payload any) error {
	p.log.Debug("Publish process", map[string]interface{}{
		"to_topic": toTopic,
		"payload":  payload,
	})

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return NewPublishProcessError("failed to marhal payload", err)
	}

	if err := p.pub.Publish(
		toTopic,
		message.NewMessage(
			watermill.NewUUID(),
			payloadData,
		),
	); err != nil {
		return NewPublishProcessError("failed to publish process message", err)
	}

	return nil
}
