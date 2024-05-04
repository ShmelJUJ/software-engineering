package publisher

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/publisher/dto"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"gopkg.in/tomb.v2"
)

//go:generate mockgen -package mocks -destination mocks/worker_mocks.go github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/publisher PaymentWorker

const (
	paymentGatewayService = "payment_gateway"
)

// StopReason represents the reason for stopping a PaymentWorker.
type StopReason int

const (
	Undefined StopReason = iota
	CancelledTransaction
	Shutdown
)

var (
	errCancelledTransation = errors.New("transaction cancelled")
	errShutdownWorker      = errors.New("worker shutdown")
)

// PaymentWorker defines the behavior of a worker responsible for processing payments.
type PaymentWorker interface {
	Start(context.Context) error
	Stop(StopReason) error
}

type paymentWorker struct {
	cfg          *Config
	gateway      gateway.PaymentGateway
	pub          message.Publisher
	log          logger.Logger
	tomb         tomb.Tomb
	cancelled    atomic.Bool
	cancel, stop chan struct{}
}

// NewWorker creates a new paymentWorker instance.
func NewWorker(
	cfg *Config,
	log logger.Logger,
	g gateway.PaymentGateway,
	pub message.Publisher,
) (PaymentWorker, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, NewWorkerError("failed to merge with default config", err)
	}

	return &paymentWorker{
		cfg:       cfg,
		log:       log,
		gateway:   g,
		pub:       pub,
		tomb:      tomb.Tomb{},
		cancelled: atomic.Bool{},
		cancel:    make(chan struct{}, 1),
		stop:      make(chan struct{}, 1),
	}, nil
}

// Start begins the payment processing pipeline.
func (worker *paymentWorker) Start(ctx context.Context) error {
	worker.log.Debug("Start payment worker", map[string]interface{}{
		"transaction_id": worker.gateway.TransactionID(),
	})

	paymentID, err := worker.gateway.CreatePayment(ctx)
	if err != nil {
		return NewStartError("failed to create payment", err)
	}

	return worker.proccessPayment(ctx, paymentID)
}

func (worker *paymentWorker) proccessPayment(ctx context.Context, paymentID string) error {
	worker.log.Debug("Payment worker start payment processing", map[string]interface{}{
		"transaction_id": worker.gateway.TransactionID(),
		"payment_id":     paymentID,
	})

	ticker := time.NewTicker(worker.gateway.Timeout())
	defer ticker.Stop()

	paymentProccessingTimeout := time.After(worker.cfg.PaymentProccessingTime)

	var retries atomic.Int32

	for {
		select {
		case <-worker.tomb.Dying():
			if errors.Is(worker.tomb.Err(), errCancelledTransation) {
				worker.log.Debug("Transaction was cancelled by the user", map[string]interface{}{
					"transaction_id": worker.gateway.TransactionID(),
				})
			} else if errors.Is(worker.tomb.Err(), errShutdownWorker) {
				worker.log.Debug("Payment worker shutdown", map[string]interface{}{
					"transaction_id": worker.gateway.TransactionID(),
					"payment_id":     paymentID,
				})

				return worker.handleFailedTransaction("Payment gateway service is shutting down")
			}

			return worker.tomb.Err()

		case <-paymentProccessingTimeout:
			return worker.handleFailedTransaction("Maximum transaction processing time has expired")

		case <-ticker.C:
			worker.tomb.Go(func() error {
				if int(retries.Load()) == worker.gateway.Retries() {
					return worker.handleFailedTransaction("Number of transaction verification retries has expired")
				}

				paymentStatus, err := worker.gateway.CheckStatus(ctx, paymentID)
				if err != nil {
					return worker.handleFailedTransaction(fmt.Sprintf("Failed to check payment status: %s", err.Error()))
				}

				if worker.cancelled.Load() {
					return nil
				}

				switch paymentStatus {
				case gateway.Succeeded:
					return worker.handleSucceededTransaction()

				case gateway.Cancelled:
					return worker.handleFailedTransaction("Payment gateway cancelled the transaction")
				}

				worker.log.Debug("Payment worker waiting for a change in the transaction status", map[string]interface{}{
					"transaction_id": worker.gateway.TransactionID(),
					"payment_id":     paymentID,
					"retries":        retries.Load(),
				})

				retries.Add(1)

				return nil
			})
		}
	}
}

func (worker *paymentWorker) handleFailedTransaction(reason string) error {
	worker.log.Debug("Payment worker handle failed transaction", map[string]interface{}{
		"transaction_id": worker.gateway.TransactionID(),
		"reason":         reason,
	})

	failedTransaction := &dto.FailedTransaction{
		TransactionID: worker.gateway.TransactionID(),
		Reason:        reason,
	}

	monitorDTO := &dto.Process{
		From:    paymentGatewayService,
		ToTopic: worker.cfg.FailedTransactionTopic,
		Payload: failedTransaction,
	}

	payload, err := monitorDTO.Encode()
	if err != nil {
		return NewProccessPaymentError("failed to encode transaction", err)
	}

	if err := worker.pub.Publish(
		worker.cfg.MonitorProcessTopic,
		message.NewMessage(
			watermill.NewUUID(),
			payload,
		),
	); err != nil {
		return NewProccessPaymentError("failed to publish message", err)
	}

	worker.tomb.Kill(nil)
	worker.cancelled.Store(true)

	return worker.tomb.Wait()
}

func (worker *paymentWorker) handleSucceededTransaction() error {
	worker.log.Debug("Payment worker handle succeeded transaction", map[string]interface{}{
		"transaction_id": worker.gateway.TransactionID(),
	})

	succeededTransaction := &dto.SucceededTransaction{
		TransactionID: worker.gateway.TransactionID(),
	}

	monitorDTO := &dto.Process{
		From:    paymentGatewayService,
		ToTopic: worker.cfg.SucceededTransactionTopic,
		Payload: succeededTransaction,
	}

	payload, err := monitorDTO.Encode()
	if err != nil {
		return NewProccessPaymentError("failed to encode monitor dto", err)
	}

	if err := worker.pub.Publish(
		worker.cfg.MonitorProcessTopic,
		message.NewMessage(
			watermill.NewUUID(),
			payload,
		),
	); err != nil {
		return NewProccessPaymentError("failed to publish message", err)
	}

	worker.tomb.Kill(nil)
	worker.cancelled.Store(true)

	return worker.tomb.Wait()
}

// Stop stops the payment processing based on the specified reason.
func (worker *paymentWorker) Stop(reason StopReason) error {
	worker.log.Debug("Stop payment worker", map[string]interface{}{
		"transaction_id": worker.gateway.TransactionID(),
		"reason":         reason,
	})

	worker.cancelled.Store(true)

	switch reason {
	case CancelledTransaction:
		worker.tomb.Kill(errCancelledTransation)

	case Shutdown:
		worker.tomb.Kill(errShutdownWorker)
	}

	return worker.tomb.Wait()
}
