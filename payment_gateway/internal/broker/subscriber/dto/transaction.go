package dto

import (
	"encoding/json"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway"
)

// Transaction represents a basic transaction with specific fields.
type Transaction struct {
	TransactionID string `json:"transaction_id"`
	Value         string `json:"value"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
}

// ToTransactionInfo converts a Transaction to a gateway.TransactionInfo object.
func (t *Transaction) ToTransactionInfo() *gateway.TransactionInfo {
	return &gateway.TransactionInfo{
		TransactionID: t.TransactionID,
		Value:         t.Value,
		Currency:      t.Currency,
	}
}

type TransactionUser struct {
	UserID   string `json:"user_id"`
	WalletID string `json:"wallet_id"`
}

// ProcessedTransaction represents a transaction that has been fully processed.
type ProcessedTransaction struct {
	Transaction *Transaction     `json:"transaction"`
	Sender      *TransactionUser `json:"sender"`
	Receiver    *TransactionUser `json:"receiver"`
}

// Decode populates a ProcessedTransaction object from JSON data.
func (t *ProcessedTransaction) Decode(data []byte) error {
	return json.Unmarshal(data, &t)
}

// CancelledTransaction represents a cancelled transaction with its ID.
type CancelledTransaction struct {
	TransactionID string `json:"transaction_id"`
}

// Decode populates a CancelledTransaction object from JSON data.
func (t *CancelledTransaction) Decode(data []byte) error {
	return json.Unmarshal(data, &t)
}
