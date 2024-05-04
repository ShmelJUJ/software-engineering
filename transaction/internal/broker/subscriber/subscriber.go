package subscriber

import (
	"context"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/subscriber/dto"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/repository"
	"github.com/ThreeDotsLabs/watermill/message"
)

// TransactionSubscriber represents a service that subscribes to transaction-related messages
// and handles them based on their type (succeeded or failed).
type TransactionSubscriber struct {
	cfg             *Config
	log             logger.Logger
	sub             message.Subscriber
	router          *message.Router
	transactionRepo repository.TransactionRepo
}

// NewTransactionSubscriber creates a new TransactionSubscriber instance with the provided dependencies.
func NewTransactionSubscriber(
	cfg *Config,
	log logger.Logger,
	sub message.Subscriber,
	router *message.Router,
	transactionRepo repository.TransactionRepo,
) (*TransactionSubscriber, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, NewTransactionSubscriberError("failed to set default config", err)
	}

	return &TransactionSubscriber{
		cfg:             cfg,
		log:             log,
		sub:             sub,
		router:          router,
		transactionRepo: transactionRepo,
	}, nil
}

// RegisterSucceededTransactionHandler registers a handler for succeeded transaction messages.
func (s *TransactionSubscriber) RegisterSucceededTransactionHandler() {
	s.log.Debug("Register succeeded transaction handler", map[string]interface{}{})

	s.router.AddNoPublisherHandler(
		"succeeded_transaction",
		s.cfg.SucceededTransactionTopic,
		s.sub,
		s.handleSucceededTransaction,
	)
}

func (s *TransactionSubscriber) handleSucceededTransaction(msg *message.Message) error {
	ctx := context.Background()

	succeededTransaction := &dto.SucceededTransaction{}
	if err := succeededTransaction.Decode(msg.Payload); err != nil {
		s.log.Error("failed to decode succeeded transaction", map[string]interface{}{
			"error": err,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	s.log.Debug("Start handle succeeded transaction", map[string]interface{}{
		"transaction_id": succeededTransaction.TransactionID,
	})

	if err := s.transactionRepo.ChangeTransactionStatus(ctx, succeededTransaction.TransactionID, model.Succeeded); err != nil {
		s.log.Error("failed to change transaction status", map[string]interface{}{
			"error":          err,
			"status":         model.Succeeded,
			"transaction_id": succeededTransaction.TransactionID,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	return nil
}

// RegisterFailedTransactionHandler registers a handler for failed transaction messages.
func (s *TransactionSubscriber) RegisterFailedTransactionHandler() {
	s.log.Debug("Register failed transaction handler", map[string]interface{}{})

	s.router.AddNoPublisherHandler(
		"failed_transaction",
		s.cfg.FailedTransactionTopic,
		s.sub,
		s.handleFailedTransaction,
	)
}

func (s *TransactionSubscriber) handleFailedTransaction(msg *message.Message) error {
	ctx := context.Background()

	failedTransaction := &dto.FailedTransaction{}
	if err := failedTransaction.Decode(msg.Payload); err != nil {
		s.log.Error("failed to decode failed transaction", map[string]interface{}{
			"error": err,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	s.log.Debug("Start handle failed transaction", map[string]interface{}{
		"transaction_id": failedTransaction.TransactionID,
	})

	if err := s.transactionRepo.CancelTransaction(ctx, failedTransaction.TransactionID, failedTransaction.Reason); err != nil {
		s.log.Error("failed to cancel transaction", map[string]interface{}{
			"error":          err,
			"reason":         failedTransaction.Reason,
			"transaction_id": failedTransaction.TransactionID,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	return nil
}

// Run starts the transaction subscriber's router.
func (s *TransactionSubscriber) Run(ctx context.Context) error {
	s.log.Debug("Run transaction subscriber", map[string]interface{}{})

	return s.router.Run(ctx)
}
