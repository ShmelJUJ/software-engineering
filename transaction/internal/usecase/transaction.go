package usecase

import (
	"context"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher/dto"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/repository"
)

//go:generate mockgen -package mocks -destination mocks/transaction_usecase_mocks.go github.com/ShmelJUJ/software-engineering/transaction/internal/usecase TransactionUsecase

// TransactionUsecase defines the interface for transaction-related use cases.
type TransactionUsecase interface {
	GetTransaction(ctx context.Context, transactionID string) (*model.Transaction, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetTransactionStatus(ctx context.Context, transactionID string) (model.TransactionStatus, error)
	CancelTransaction(ctx context.Context, transactionID string, reason string) error
	AcceptTransaction(ctx context.Context, transactionID string, sender *model.TransactionUser) error
	ChangeTransactionStatus(ctx context.Context, transactionID string, status model.TransactionStatus) error
	UpdateTransaction(ctx context.Context, updatedTransaction *model.Transaction) error
}

type transactionUsecase struct {
	transactionRepo      repository.TransactionRepo
	transactionPublisher publisher.TransactionPublisher
	log                  logger.Logger
}

// NewTransactionUsecase creates a new instance of TransactionUsecase.
func NewTransactionUsecase(
	transactionRepo repository.TransactionRepo,
	transactionPublisher publisher.TransactionPublisher,
	log logger.Logger,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo:      transactionRepo,
		transactionPublisher: transactionPublisher,
		log:                  log,
	}
}

// GetTransaction retrieves a transaction by its ID.
func (usecase *transactionUsecase) GetTransaction(ctx context.Context, transactionID string) (*model.Transaction, error) {
	usecase.log.Debug("Get transaction usecase", map[string]interface{}{
		"transaction_id": transactionID,
	})

	return usecase.transactionRepo.GetTransaction(ctx, transactionID)
}

// CreateTransaction creates a new transaction.
func (usecase *transactionUsecase) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	usecase.log.Debug("Create transaction usecase", map[string]interface{}{
		"transaction": transaction,
	})

	return usecase.transactionRepo.CreateTransaction(ctx, transaction)
}

// GetTransactionStatus retrieves the status of a transaction by its ID.
func (usecase *transactionUsecase) GetTransactionStatus(ctx context.Context, transactionID string) (model.TransactionStatus, error) {
	usecase.log.Debug("Get transaction status usecase", map[string]interface{}{
		"transaction_id": transactionID,
	})

	return usecase.transactionRepo.GetTransactionStatus(ctx, transactionID)
}

// CancelTransaction cancels a transaction with a specified reason.
func (usecase *transactionUsecase) CancelTransaction(ctx context.Context, transactionID string, reason string) error {
	usecase.log.Debug("Cancel transaction usecase", map[string]interface{}{
		"transaction_id": transactionID,
		"reason":         reason,
	})

	return usecase.transactionRepo.CancelTransaction(ctx, transactionID, reason)
}

// AcceptTransaction accepts a transaction initiated by a sender.
func (usecase *transactionUsecase) AcceptTransaction(ctx context.Context, transactionID string, sender *model.TransactionUser) error {
	usecase.log.Debug("Accept transaction usecase", map[string]interface{}{
		"transaction_id": transactionID,
	})

	if err := usecase.transactionRepo.AcceptTransaction(ctx, transactionID, sender); err != nil {
		return err
	}

	transaction, err := usecase.transactionRepo.GetTransaction(ctx, transactionID)
	if err != nil {
		return err
	}

	return usecase.transactionPublisher.PublishSucceededTransaction(dto.FromTransactionModel(transaction))
}

// UpdateTransaction updates an existing transaction.
func (usecase *transactionUsecase) UpdateTransaction(ctx context.Context, updatedTransaction *model.Transaction) error {
	usecase.log.Debug("Update transaction usecase", map[string]interface{}{
		"updated_transaction": updatedTransaction,
	})

	return usecase.transactionRepo.UpdateTransaction(ctx, updatedTransaction)
}

// ChangeTransactionStatus changes the status of a transaction.
func (usecase *transactionUsecase) ChangeTransactionStatus(ctx context.Context, transactionID string, status model.TransactionStatus) error {
	usecase.log.Debug("Change transaction status", map[string]interface{}{
		"transaction_id": transactionID,
	})

	return usecase.transactionRepo.ChangeTransactionStatus(ctx, transactionID, status)
}
