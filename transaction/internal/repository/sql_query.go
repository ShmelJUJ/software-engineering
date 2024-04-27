package repository

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
)

const (
	transactionsTable     = "transactions"
	transactionUsersTable = "transaction_users"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func getTransactionQuery(transactionID string) sq.SelectBuilder {
	return psql.
		Select(
			"transaction_id",
			"sender_id",
			"receiver_id",
			"currency",
			"amount",
			"status",
			"method",
			"canceled_reason",
			"created_at",
			"updated_at",
		).
		From(transactionsTable).
		Where(sq.Eq{
			"transaction_id": transactionID,
		})
}

func getTransactionUserQuery(transactionUserID string) sq.SelectBuilder {
	return psql.
		Select(
			"transaction_user_id",
			"user_id",
			"created_at",
			"updated_at",
		).
		From(transactionUsersTable).
		Where(sq.Eq{
			"transaction_user_id": transactionUserID,
		})
}

func createTransactionQuery(transaction *model.Transaction) sq.InsertBuilder {
	return psql.
		Insert(transactionsTable).
		Columns(
			"transaction_id",
			"receiver_id",
			"currency",
			"amount",
			"status",
			"method",
			"canceled_reason",
			"created_at",
			"updated_at",
		).
		Values(
			transaction.ID,
			transaction.ReceiverID,
			transaction.Currency,
			transaction.Amount,
			transaction.Status,
			transaction.Method,
			transaction.CanceledReason,
			transaction.CreatedAt,
			transaction.UpdatedAt,
		)
}

func createTransactionUserQuery(transactionUser *model.TransactionUser) sq.InsertBuilder {
	return psql.
		Insert(transactionUsersTable).
		Columns(
			"transaction_user_id",
			"user_id",
			"created_at",
			"updated_at",
		).
		Values(
			transactionUser.ID,
			transactionUser.UserID,
			transactionUser.CreatedAt,
			transactionUser.UpdatedAt,
		)
}

func getTransactionStatusQuery(transactionID string) sq.SelectBuilder {
	return psql.
		Select(
			"status",
		).
		From(transactionsTable).
		Where(sq.Eq{
			"transaction_id": transactionID,
		})
}

func cancelTransactionQuery(transactionID, reason string) sq.UpdateBuilder {
	return psql.
		Update(transactionsTable).
		Set("status", model.Canceled).
		Set("canceled_reason", reason).
		Set("updated_at", time.Now()).
		Where(sq.Eq{
			"transaction_id": transactionID,
		})
}

func updateTransactionQuery(transaction *model.Transaction) sq.UpdateBuilder {
	query := psql.Update(transactionsTable)

	if transaction.SenderID != nil {
		query = query.Set("sender_id", transaction.SenderID)
	}

	if transaction.ReceiverID != "" {
		query = query.Set("receiver_id", transaction.ReceiverID)
	}

	if transaction.Currency != "" {
		query = query.Set("currency", transaction.Currency)
	}

	if transaction.Amount != 0 {
		query = query.Set("amount", transaction.Amount)
	}

	if transaction.Status != 0 {
		query = query.Set("status", transaction.Status)
	}

	if transaction.CanceledReason != "" {
		query = query.Set("canceled_reason", transaction.CanceledReason)
	}

	query = query.
		Set("updated_at", time.Now()).
		Where(sq.Eq{
			"transaction_id": transaction.ID,
		})

	return query
}

func acceptTransactionQuery(transactionID, senderID string) sq.UpdateBuilder {
	return psql.
		Update(transactionsTable).
		Set("status", model.Processed).
		Set("sender_id", senderID).
		Where(sq.Eq{
			"transaction_id": transactionID,
		})
}

func changeTransactionStatusQuery(transactionID string, status model.TransactionStatus) sq.UpdateBuilder {
	return psql.
		Update(transactionsTable).
		Set("status", status).
		Where(sq.Eq{
			"transaction_id": transactionID,
		})
}
