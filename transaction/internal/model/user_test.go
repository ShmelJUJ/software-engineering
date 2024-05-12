package model_test

import (
	"testing"

	dto "github.com/ShmelJUJ/software-engineering/transaction/internal/generated/models"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestFromCreateTransactionUserDTO(t *testing.T) {
	t.Parallel()

	type args struct {
		transactionUserDTO *dto.CreateTransactionUserRequest
	}

	userID := "test-user-id"
	walletID := "test-wallet-id"

	testcases := []struct {
		name                    string
		args                    args
		expectedTransactionUser *model.TransactionUser
	}{
		{
			name: "Successfully convert CreateTransactionUserDTO to user model",
			args: args{
				transactionUserDTO: &dto.CreateTransactionUserRequest{
					UserID:   (*strfmt.UUID)(&userID),
					WalletID: (*strfmt.UUID)(&walletID),
				},
			},
			expectedTransactionUser: &model.TransactionUser{
				UserID:   userID,
				WalletID: walletID,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualTransactionUser := model.FromCreateTransactionUserDTO(testcase.args.transactionUserDTO)
			assert.Equal(t, testcase.expectedTransactionUser.UserID, actualTransactionUser.UserID)
		})
	}
}

func TestFromAcceptTransactionUserDTO(t *testing.T) {
	t.Parallel()

	type args struct {
		transactionUserDTO *dto.AcceptTransactionUserRequest
	}

	userID := "test-user-id"
	walletID := "test-wallet-id"

	testcases := []struct {
		name                    string
		args                    args
		expectedTransactionUser *model.TransactionUser
	}{
		{
			name: "Successfully convert CreateTransactionUserDTO to user model",
			args: args{
				transactionUserDTO: &dto.AcceptTransactionUserRequest{
					UserID:   (*strfmt.UUID)(&userID),
					WalletID: (*strfmt.UUID)(&walletID),
				},
			},
			expectedTransactionUser: &model.TransactionUser{
				UserID:   userID,
				WalletID: walletID,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualTransactionUser := model.FromAcceptTransactionUserDTO(testcase.args.transactionUserDTO)
			assert.Equal(t, testcase.expectedTransactionUser.UserID, actualTransactionUser.UserID)
		})
	}
}

func TestToGetTransactionUserDTO(t *testing.T) {
	t.Parallel()

	id := "test-id"
	userID := "test-user-id"
	walletID := "test-wallet-id"

	testcases := []struct {
		name                          string
		transactionUser               *model.TransactionUser
		expectedGetTransactionUserDTO *dto.GetTransactionUserResponse
	}{
		{
			name: "Successfully convert CreateTransactionUserDTO to user model",
			transactionUser: &model.TransactionUser{
				ID:       id,
				UserID:   userID,
				WalletID: walletID,
			},
			expectedGetTransactionUserDTO: &dto.GetTransactionUserResponse{
				UserID:   (*strfmt.UUID)(&userID),
				WalletID: (*strfmt.UUID)(&walletID),
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			getTransactionUserDTO := testcase.transactionUser.ToGetTransactionUserDTO()
			assert.Equal(t, testcase.expectedGetTransactionUserDTO, getTransactionUserDTO)
		})
	}
}
