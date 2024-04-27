package publisher

import (
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher/dto"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

//go:generate mockgen -package mocks -destination mocks/transaction_publisher_mocks.go github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher TransactionPublisher

// TransactionPublisher is an interface for publishing succeeded transactions.
type TransactionPublisher interface {
	PublishSucceededTransaction(transaction *dto.ProcessedTransaction) error
}

type transactionPublisher struct {
	cfg *Config
	log logger.Logger
	pub message.Publisher
}

// NewTransactionPublisher creates a new TransactionPublisher instance.
func NewTransactionPublisher(
	cfg *Config,
	log logger.Logger,
	pub message.Publisher,
) (TransactionPublisher, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, NewTransactionPublisherError("failed to set default config", err)
	}

	return &transactionPublisher{
		cfg: cfg,
		log: log,
		pub: pub,
	}, nil
}

// PublishSucceededTransaction publishes a succeeded transaction.
func (p *transactionPublisher) PublishSucceededTransaction(transaction *dto.ProcessedTransaction) error {
	p.log.Debug("Start publish succeeded transaction", map[string]interface{}{
		"transaction": transaction,
	})

	payload, err := transaction.Encode()
	if err != nil {
		return NewPublishSucceededTransactionError("failed to encode succeeded transaction", err)
	}

	if err = p.pub.Publish(
		p.cfg.ProcessedTransactionTopic,
		message.NewMessage(
			watermill.NewUUID(),
			payload,
		),
	); err != nil {
		return NewPublishSucceededTransactionError("failed to publish succeded transaction", err)
	}

	return nil
}
