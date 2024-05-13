package model

import (
	"time"

	dto "github.com/ShmelJUJ/software-engineering/transaction/internal/generated/models"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

// Represents how the transaction user structure is stored in the database.
type TransactionUser struct {
	ID        string    `db:"transaction_user_id"`
	UserID    string    `db:"user_id"`
	WalletID  string    `db:"wallet_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// FromCreateTransactionUserDTO converts a CreateTransactionUserRequest DTO to a TransactionUser entity.
func FromCreateTransactionUserDTO(transactionUserDTO *dto.CreateTransactionUserRequest) *TransactionUser {
	if transactionUserDTO == nil || transactionUserDTO.UserID == nil || transactionUserDTO.WalletID == nil {
		return nil
	}

	return &TransactionUser{
		ID:        uuid.NewString(),
		UserID:    transactionUserDTO.UserID.String(),
		WalletID:  transactionUserDTO.WalletID.String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// FromAcceptTransactionUserDTO converts an AcceptTransactionUserRequest DTO to a TransactionUser entity.
func FromAcceptTransactionUserDTO(transactionUserDTO *dto.AcceptTransactionUserRequest) *TransactionUser {
	if transactionUserDTO == nil || transactionUserDTO.UserID == nil || transactionUserDTO.WalletID == nil {
		return nil
	}

	return &TransactionUser{
		ID:        uuid.NewString(),
		UserID:    transactionUserDTO.UserID.String(),
		WalletID:  transactionUserDTO.WalletID.String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ToGetTransactionUserDTO converts a TransactionUser entity to a GetTransactionUserResponse DTO.
func (transactionUser *TransactionUser) ToGetTransactionUserDTO() *dto.GetTransactionUserResponse {
	if transactionUser.UserID == "" {
		return nil
	}

	return &dto.GetTransactionUserResponse{
		UserID:   (*strfmt.UUID)(&transactionUser.UserID),
		WalletID: (*strfmt.UUID)(&transactionUser.WalletID),
	}
}
