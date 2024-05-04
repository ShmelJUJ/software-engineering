package dto_test

import (
	"testing"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/subscriber/dto"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway"
	"github.com/stretchr/testify/assert"
)

const (
	transactionID = "test-transaction-id"
	value         = "test-value"
	currency      = "test-currency"
	paymentMethod = "test-payment-method"
)

func TestToTransactionInfo(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name                       string
		transaction                *dto.Transaction
		expectedGatewayTransaction *gateway.TransactionInfo
	}{
		{
			name: "Successfully convert to gateway transaction info",
			transaction: &dto.Transaction{
				TransactionID: transactionID,
				Value:         value,
				Currency:      currency,
				PaymentMethod: paymentMethod,
			},
			expectedGatewayTransaction: &gateway.TransactionInfo{
				TransactionID: transactionID,
				Value:         value,
				Currency:      currency,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualGatewayTransaction := testcase.transaction.ToTransactionInfo()

			assert.Equal(t, testcase.expectedGatewayTransaction, actualGatewayTransaction)
		})
	}
}

func TestProcessedTransactionDecode(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}

	testcases := []struct {
		name                string
		args                args
		transaction         *dto.ProcessedTransaction
		expectedTransaction *dto.ProcessedTransaction
		expectedErr         error
	}{
		{
			name: "Successfully decode succeeded transaction",
			args: args{
				data: []byte(`{"transaction":{"transaction_id":"123","value":"v"},"sender_id":"456","receiver_id":"789"}`),
			},
			transaction: &dto.ProcessedTransaction{},
			expectedTransaction: &dto.ProcessedTransaction{
				Transaction: &dto.Transaction{
					TransactionID: "123",
					Value:         "v",
				},
				SenderID:   "456",
				ReceiverID: "789",
			},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			err := testcase.transaction.Decode(testcase.args.data)

			assert.Equal(t, testcase.expectedTransaction, testcase.transaction)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestCancelledTransactionDecode(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}

	testcases := []struct {
		name                string
		args                args
		transaction         *dto.CancelledTransaction
		expectedTransaction *dto.CancelledTransaction
		expectedErr         error
	}{
		{
			name: "Successfully decode cancelled transaction",
			args: args{
				data: []byte(`{"transaction_id":"123"}`),
			},
			transaction: &dto.CancelledTransaction{},
			expectedTransaction: &dto.CancelledTransaction{
				TransactionID: "123",
			},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			err := testcase.transaction.Decode(testcase.args.data)

			assert.Equal(t, testcase.expectedTransaction, testcase.transaction)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}
