package dto

import (
	"encoding/json"
	"strconv"

	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
)

// Transaction represents a transaction with details such as ID, value, currency, and payment method.
type Transaction struct {
	TransactionID string `json:"transaction_id"`
	Value         string `json:"value"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
}

type TransactionUser struct {
	UserID   string `json:"user_id"`
	WalletID string `json:"wallet_id"`
}

// ProcessedTransaction represents a processed transaction, including the original transaction details, sender ID, and receiver ID.
type ProcessedTransaction struct {
	Transaction *Transaction     `json:"transaction"`
	Sender      *TransactionUser `json:"sender"`
	Receiver    *TransactionUser `json:"receiver"`
}

// Encode serializes a ProcessedTransaction into a JSON-encoded byte slice.
func (t *ProcessedTransaction) Encode() ([]byte, error) {
	data, err := json.Marshal(&t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FromTransactionModel creates a ProcessedTransaction from a model.Transaction object.
func FromTransactionModel(transaction *model.Transaction) *ProcessedTransaction {
	amount := strconv.Itoa(int(transaction.Amount))

	return &ProcessedTransaction{
		Transaction: &Transaction{
			TransactionID: transaction.ID,
			Value:         amount,
			Currency:      transaction.Currency,
			PaymentMethod: transaction.Method,
		},
		Sender: &TransactionUser{
			UserID:   transaction.Sender.UserID,
			WalletID: transaction.Sender.WalletID,
		},
		Receiver: &TransactionUser{
			UserID:   transaction.Receiver.UserID,
			WalletID: transaction.Receiver.WalletID,
		},
	}
}
