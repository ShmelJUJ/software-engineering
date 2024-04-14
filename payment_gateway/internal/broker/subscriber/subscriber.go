package subscriber

import (
	"context"
	"fmt"
	"sync"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/publisher"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/subscriber/dto"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway/algorand"
	gateway_stub "github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway/stub"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alitto/pond"
)

// TransactionSubscriber represents a subscriber handling transaction-related messages.
type TransactionSubscriber struct {
	paymentWorkers sync.Map
	router         *message.Router
	sub            message.Subscriber
	pub            message.Publisher
	publisherCfg   *publisher.Config
	pool           *pond.WorkerPool
	cfg            *Config
	algorandCfg    *algorand.Config
	log            logger.Logger
}

// NewTransactionSubscriber creates a new instance of TransactionSubscriber.
func NewTransactionSubscriber(
	cfg *Config,
	log logger.Logger,
	router *message.Router,
	sub message.Subscriber,
	pub message.Publisher,
	publisherCfg *publisher.Config,
	algorandCfg *algorand.Config,
) (*TransactionSubscriber, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, NewTransactionSubscriberError("failed to set default config", err)
	}

	return &TransactionSubscriber{
		paymentWorkers: sync.Map{},
		router:         router,
		sub:            sub,
		pub:            pub,
		publisherCfg:   publisherCfg,
		cfg:            cfg,
		algorandCfg:    algorandCfg,
		log:            log,

		pool: pond.New(
			cfg.PoolCfg.NumWorkers,
			cfg.PoolCfg.TasksCapacity,
			pond.IdleTimeout(cfg.PoolCfg.IdleTimeout),
			pond.MinWorkers(cfg.PoolCfg.MinWorkers),
		),
	}, nil
}

// RegisterProcessedTransactionHandler registers a handler for processed transaction messages.
func (s *TransactionSubscriber) RegisterProcessedTransactionHandler() {
	s.log.Debug("Register processed transaction handler", map[string]interface{}{})

	s.router.AddNoPublisherHandler(
		"processed_transaction",
		s.cfg.ProcessedTransactionTopic,
		s.sub,
		s.handleProcessedTransaction,
	)
}

func (s *TransactionSubscriber) handleProcessedTransaction(msg *message.Message) error {
	ctx := context.Background()

	processedTransaction := &dto.ProcessedTransaction{}
	if err := processedTransaction.Decode(msg.Payload); err != nil {
		s.log.Error("failed to decode processed transaction", map[string]interface{}{
			"error": err,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	transactionID := processedTransaction.Transaction.TransactionID

	s.log.Debug("Start handle processed transaction", map[string]interface{}{
		"transaction_id": transactionID,
	})

	paymentGateway, err := s.getPaymentGateway(processedTransaction)
	if err != nil {
		s.log.Error("failed to get payment gateway", map[string]interface{}{
			"transaction_id": transactionID,
			"error":          err,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	s.pool.Submit(func() {
		worker, err := publisher.NewWorker(s.publisherCfg, s.log, paymentGateway, s.pub)
		if err != nil {
			s.log.Error("failed to make new worker", map[string]interface{}{
				"transaction_id": transactionID,
				"error":          err,
			})

			return
		}

		s.paymentWorkers.Store(transactionID, worker)

		if err := worker.Start(ctx); err != nil {
			s.log.Error("failed to start transaction worker", map[string]interface{}{
				"transaction_id": transactionID,
				"error":          err,
			})
		}

		s.paymentWorkers.Delete(transactionID)
	})

	return nil
}

func (s *TransactionSubscriber) getPaymentGateway(processedTransaction *dto.ProcessedTransaction) (gateway.PaymentGateway, error) {
	paymentMethod := processedTransaction.Transaction.PaymentMethod

	if paymentMethod == "algorand" {
		if s.algorandCfg.IsTest {
			return gateway_stub.New(processedTransaction.Transaction.ToTransactionInfo()), nil
		}

		algorandGateway, err := algorand.New(
			s.algorandCfg,
			processedTransaction.Transaction.ToTransactionInfo(),
			&algorand.UserData{},
			&algorand.UserData{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create algorand gateway: %w", err)
		}

		return algorandGateway, nil
	}

	return nil, fmt.Errorf("cannot handle %s payment gateway", paymentMethod)
}

// RegisterCancelledTransactionHandler registers a handler for cancelled transaction messages.
func (s *TransactionSubscriber) RegisterCancelledTransactionHandler() {
	s.log.Debug("Register cancelled transaction handler", map[string]interface{}{})

	s.router.AddNoPublisherHandler(
		"cancelled_transaction",
		s.cfg.CancelledTransactionTopic,
		s.sub,
		s.handleCancelledTransaction,
	)
}

func (s *TransactionSubscriber) handleCancelledTransaction(msg *message.Message) error {
	cancelledTransaction := &dto.CancelledTransaction{}
	if err := cancelledTransaction.Decode(msg.Payload); err != nil {
		s.log.Error("failed to decode cancelled transaction", map[string]interface{}{
			"error": err,
		})

		return nil //nolint:nilerr // it is necessary for a commit to occur and not to hang in a endless loop
	}

	s.log.Debug("Start handle cancelled transaction", map[string]interface{}{
		"transaction_id": cancelledTransaction.TransactionID,
	})

	if value, ok := s.paymentWorkers.Load(cancelledTransaction.TransactionID); ok {
		paymentWorker, ok := value.(publisher.PaymentWorker)
		if !ok {
			s.log.Error("failed to get payment worker from payment workers map", map[string]interface{}{})

			return nil
		}

		if err := paymentWorker.Stop(publisher.CancelledTransaction); err != nil {
			s.log.Error("failed to stop payment worker", map[string]interface{}{
				"error": err,
			})
		}

		s.paymentWorkers.Delete(cancelledTransaction.TransactionID)
	}

	return nil
}

// Run starts the transaction subscriber's router.
func (s *TransactionSubscriber) Run(ctx context.Context) error {
	s.log.Debug("Run transaction subscriber", map[string]interface{}{})

	return s.router.Run(ctx)
}

// Stop stops the transaction subscriber by closing the router and stopping all associated payment workers.
func (s *TransactionSubscriber) Stop() error {
	s.log.Debug("Stop transaction subscriber", map[string]interface{}{})

	if err := s.router.Close(); err != nil {
		return err
	}

	s.paymentWorkers.Range(func(_, value any) bool {
		paymentWorker, ok := value.(publisher.PaymentWorker)
		if !ok {
			return true
		}

		if err := paymentWorker.Stop(publisher.Shutdown); err != nil {
			s.log.Error("failed to stop payment worker", map[string]interface{}{
				"error": err,
			})
		}

		return true
	})

	s.pool.StopAndWait()

	return nil
}
