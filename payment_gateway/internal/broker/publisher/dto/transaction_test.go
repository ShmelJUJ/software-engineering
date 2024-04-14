package dto_test

import (
	"testing"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/publisher/dto"
	"github.com/stretchr/testify/assert"
)

func TestSucceededTransactionEncode(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		transaction  *dto.SucceededTransaction
		expectedData []byte
		expectedErr  error
	}{
		{
			name:         "With empty succeeded transaction",
			transaction:  &dto.SucceededTransaction{},
			expectedData: []byte(`{"transaction_id":""}`),
			expectedErr:  nil,
		},
		{
			name: "With full filled succeeded transaction",
			transaction: &dto.SucceededTransaction{
				TransactionID: "123",
			},
			expectedData: []byte(`{"transaction_id":"123"}`),
			expectedErr:  nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualData, err := testcase.transaction.Encode()

			assert.Equal(t, testcase.expectedData, actualData)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestFailedTransactionEncode(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name         string
		transaction  *dto.FailedTransaction
		expectedData []byte
		expectedErr  error
	}{
		{
			name:         "With empty failed transaction",
			transaction:  &dto.FailedTransaction{},
			expectedData: []byte(`{"transaction_id":"","reason":""}`),
			expectedErr:  nil,
		},
		{
			name: "With full filled failed transaction",
			transaction: &dto.FailedTransaction{
				TransactionID: "123",
				Reason:        "test",
			},
			expectedData: []byte(`{"transaction_id":"123","reason":"test"}`),
			expectedErr:  nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualData, err := testcase.transaction.Encode()

			assert.Equal(t, testcase.expectedData, actualData)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}
