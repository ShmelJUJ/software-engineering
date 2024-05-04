package model

import (
	"time"

	dto "github.com/ShmelJUJ/software-engineering/transaction/internal/generated/models"
	"github.com/google/uuid"
)

type TransactionStatus int

const (
	Undefined TransactionStatus = iota
	Created
	Processed
	Canceled
	Failed
	Succeeded
)

func (ts TransactionStatus) String() string {
	switch ts {
	case Created:
		return "created"
	case Processed:
		return "processed"
	case Canceled:
		return "canceled"
	case Failed:
		return "failed"
	case Succeeded:
		return "succeeded"
	default:
		return "undefined"
	}
}

// Represents how the transaction structure is stored in the database.
type Transaction struct {
	ID             string            `db:"transaction_id"`
	SenderID       *string           `db:"sender_id"`
	ReceiverID     string            `db:"receiver_id"`
	Currency       string            `db:"currency"`
	Amount         int64             `db:"amount"`
	Status         TransactionStatus `db:"status"`
	Method         string            `db:"method"`
	CanceledReason string            `db:"canceled_reason"`
	CreatedAt      time.Time         `db:"created_at"`
	UpdatedAt      time.Time         `db:"updated_at"`

	Sender   *TransactionUser `db:"-"`
	Receiver *TransactionUser `db:"-"`
}

// FromCreateTransactionDTO creates a Transaction from a CreateTransactionRequest DTO.
func FromCreateTransactionDTO(transactionDTO *dto.CreateTransactionRequest) *Transaction {
	if transactionDTO == nil || transactionDTO.MoneyInfo == nil ||
		transactionDTO.MoneyInfo.Amount == nil || transactionDTO.MoneyInfo.Currency == nil ||
		transactionDTO.MoneyInfo.Method == nil {
		return nil
	}

	receiver := FromCreateTransactionUserDTO(transactionDTO.Receiver)
	if receiver == nil {
		return nil
	}

	return &Transaction{
		ID:         uuid.NewString(),
		ReceiverID: receiver.ID,
		Currency:   *transactionDTO.MoneyInfo.Currency,
		Amount:     *transactionDTO.MoneyInfo.Amount,
		Status:     Created,
		Method:     *transactionDTO.MoneyInfo.Method,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),

		Receiver: receiver,
	}
}

// FromEditTransactionDTO creates a Transaction from an EditTransactionRequest DTO.
func FromEditTransactionDTO(transactionID string, transactionDTO *dto.EditTransactionRequest) *Transaction {
	if transactionDTO.MoneyInfo == nil {
		return nil
	}

	return &Transaction{
		ID:        transactionID,
		Currency:  transactionDTO.MoneyInfo.Currency,
		Amount:    transactionDTO.MoneyInfo.Amount,
		Method:    transactionDTO.MoneyInfo.Method,
		UpdatedAt: time.Now(),
	}
}

// ToGetTransactionDTO converts a Transaction to a GetTransactionResponse DTO.
func (transaction *Transaction) ToGetTransactionDTO() *dto.GetTransactionResponse {
	transactionStatus := transaction.Status.String()

	transactionResponse := &dto.GetTransactionResponse{
		Amount:   &transaction.Amount,
		Currency: &transaction.Currency,
		Method:   &transaction.Method,
		Status:   &transactionStatus,
		Receiver: transaction.Receiver.ToGetTransactionUserDTO(),
	}

	if transaction.Sender != nil {
		transactionResponse.Sender = transaction.Sender.ToGetTransactionUserDTO()
	}

	return transactionResponse
}
