package publisher

import (
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher/dto"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

//go:generate mockgen -package mocks -destination mocks/transaction_publisher_mocks.go github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher TransactionPublisher

const (
	transactionService    = "transaction"
	paymentGatewayService = "payment_gateway"
)

// TransactionPublisher is an interface for publishing processed transactions.
type TransactionPublisher interface {
	PublishProcessedTransaction(transaction *dto.ProcessedTransaction) error
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

// PublishProcessedTransaction publishes a processed transaction.
func (p *transactionPublisher) PublishProcessedTransaction(transaction *dto.ProcessedTransaction) error {
	p.log.Debug("Start publish processed transaction", map[string]interface{}{
		"transaction": transaction,
	})

	monitorDTO := &dto.Process{
		From:    transactionService,
		ToTopic: p.cfg.ProcessedTransactionTopic,
		Payload: transaction,
	}

	payload, err := monitorDTO.Encode()
	if err != nil {
		return NewPublishProcessedTransactionError("failed to encode monitor process dto", err)
	}

	if err = p.pub.Publish(
		p.cfg.ProcessMonitorTopic,
		message.NewMessage(
			watermill.NewUUID(),
			payload,
		),
	); err != nil {
		return NewPublishProcessedTransactionError("failed to publish processed transaction", err)
	}

	return nil
}
