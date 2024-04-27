package repository

import (
	"context"
	"fmt"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/pkg/postgres"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -package mocks -destination mocks/transaction_repository_mocks.go github.com/ShmelJUJ/software-engineering/transaction/internal/repository TransactionRepo

// TransactionRepo defines the interface for transaction-related operations.
type TransactionRepo interface {
	GetTransaction(ctx context.Context, transactionID string) (*model.Transaction, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetTransactionStatus(ctx context.Context, transactionID string) (model.TransactionStatus, error)
	CancelTransaction(ctx context.Context, transactionID string, reason string) error
	AcceptTransaction(ctx context.Context, transactionID string, sender *model.TransactionUser) error
	ChangeTransactionStatus(ctx context.Context, transactionID string, status model.TransactionStatus) error
	UpdateTransaction(ctx context.Context, updatedTransaction *model.Transaction) error
}

type transactionRepo struct {
	pg  *postgres.Postgres
	log logger.Logger
}

// NewTransactionRepo creates a new instance of TransactionRepo.
func NewTransactionRepo(
	pg *postgres.Postgres,
	log logger.Logger,
) TransactionRepo {
	return &transactionRepo{
		pg:  pg,
		log: log,
	}
}

// GetTransaction retrieves a transaction from the database.
func (repo *transactionRepo) GetTransaction(ctx context.Context, transactionID string) (*model.Transaction, error) {
	query := getTransactionQuery(transactionID)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, NewGetTransactionError("failed to get transaction sql query", err)
	}

	var transaction *model.Transaction

	if err = repo.pg.TrManager.Do(ctx, func(ctx context.Context) error {
		transactionConn := repo.pg.GetTransactionConn(ctx)

		rows, err := transactionConn.Query(ctx, sqlQuery, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		collectedTransaction, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Transaction])
		if err != nil {
			return err
		}

		transaction = &collectedTransaction

		receiver, err := repo.getTransactionUserInTx(ctx, transaction.ReceiverID)
		if err != nil {
			return err
		}

		transaction.Receiver = receiver

		if transaction.SenderID == nil {
			sender, err := repo.getTransactionUserInTx(ctx, *transaction.SenderID)
			if err != nil {
				return err
			}

			transaction.Sender = sender
		}

		return nil
	}); err != nil {
		return nil, NewGetTransactionError("failed to get transaction", err)
	}

	return transaction, nil
}

func (repo *transactionRepo) getTransactionUserInTx(ctx context.Context, transactionUserID string) (*model.TransactionUser, error) {
	query := getTransactionUserQuery(transactionUserID)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, NewGetTransactionUserError("failed to get transaction sql query", err)
	}

	transactionConn := repo.pg.GetTransactionConn(ctx)

	rows, err := transactionConn.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, NewGetTransactionUserError("failed to query get transaction user sql query", err)
	}
	defer rows.Close()

	collectedTransactionUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.TransactionUser])
	if err != nil {
		return nil, NewGetTransactionUserError("failed to get transaction user structure from row", err)
	}

	return &collectedTransactionUser, nil
}

// CreateTransaction creates a new transaction in the database.
func (repo *transactionRepo) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	query := createTransactionQuery(transaction)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return NewCreateTransactionError("failed to get create transaction sql query", err)
	}

	if err = repo.pg.TrManager.Do(ctx, func(ctx context.Context) error {
		transactionConn := repo.pg.GetTransactionConn(ctx)

		if err := repo.createTransactionUserInTx(ctx, transaction.Receiver); err != nil {
			return err
		}

		if _, err := transactionConn.Exec(ctx, sqlQuery, args...); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return NewCreateTransactionError("failed to create transaction", err)
	}

	return nil
}

func (repo *transactionRepo) createTransactionUserInTx(ctx context.Context, transactionUser *model.TransactionUser) error {
	query := createTransactionUserQuery(transactionUser)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return NewCreateTransactionUserError("failed to get create transaction user sql query", err)
	}

	transactionConn := repo.pg.GetTransactionConn(ctx)

	if _, err = transactionConn.Exec(ctx, sqlQuery, args...); err != nil {
		return NewCreateTransactionUserError("failed to Exec create transaction user sql query", err)
	}

	return nil
}

// GetTransactionStatus retrieves the status of a transaction from the database.
func (repo *transactionRepo) GetTransactionStatus(ctx context.Context, transactionID string) (model.TransactionStatus, error) {
	query := getTransactionStatusQuery(transactionID)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return model.Undefined, NewGetTransactionStatusError("failed to get transaction status sql query", err)
	}

	row := repo.pg.Pool.QueryRow(ctx, sqlQuery, args...)

	var transactionStatus model.TransactionStatus
	if err = row.Scan(&transactionStatus); err != nil {
		return model.Undefined, NewGetTransactionStatusError("failed to get transaction status type from row", err)
	}

	return transactionStatus, nil
}

// CancelTransaction cancels a transaction with a specified reason.
func (repo *transactionRepo) CancelTransaction(ctx context.Context, transactionID, reason string) error {
	query := cancelTransactionQuery(transactionID, reason)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return NewCancelTransactionError("failed to get cancel transaction sql query", err)
	}

	if _, err = repo.pg.Pool.Exec(ctx, sqlQuery, args...); err != nil {
		return NewCancelTransactionError("failed to Exec cancel transaction sql query", err)
	}

	return nil
}

// AcceptTransaction accepts a transaction with a specified sender.
func (repo *transactionRepo) AcceptTransaction(ctx context.Context, transactionID string, sender *model.TransactionUser) error {
	query := acceptTransactionQuery(transactionID, sender.ID)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return NewAcceptTransactionError("failed to get accept transaction sql query", err)
	}

	if err := repo.pg.TrManager.Do(ctx, func(ctx context.Context) error {
		transactionConn := repo.pg.GetTransactionConn(ctx)

		if err = repo.createTransactionUserInTx(ctx, sender); err != nil {
			return fmt.Errorf("failed to create transaction user in tx: %w", err)
		}

		if _, err = transactionConn.Exec(ctx, sqlQuery, args...); err != nil {
			return fmt.Errorf("failed to Exec accept transaction sql query: %w", err)
		}

		return nil
	}); err != nil {
		return NewAcceptTransactionError("failed to accept transaction", err)
	}

	return nil
}

// UpdateTransaction updates an existing transaction in the database.
func (repo *transactionRepo) UpdateTransaction(ctx context.Context, updatedTransaction *model.Transaction) error {
	query := updateTransactionQuery(updatedTransaction)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return NewUpdateTransactionError("failed to get update transaction sql query", err)
	}

	if _, err = repo.pg.Pool.Exec(ctx, sqlQuery, args...); err != nil {
		return NewUpdateTransactionError("failed to Exec update transaction sql query", err)
	}

	return nil
}

// ChangeTransactionStatus changes the status of a transaction in the database.
func (repo *transactionRepo) ChangeTransactionStatus(ctx context.Context, transactionID string, status model.TransactionStatus) error {
	query := changeTransactionStatusQuery(transactionID, status)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return NewChangeTransactionStatusError("failed to get change transaction status sql query", err)
	}

	if _, err = repo.pg.Pool.Exec(ctx, sqlQuery, args...); err != nil {
		return NewChangeTransactionStatusError("failed to Exec change transaction status sql query", err)
	}

	return nil
}
