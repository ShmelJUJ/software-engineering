package publisher

import (
	"context"
	"testing"
	"time"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway"
	gateway_mocks "github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway/mocks"
	kafka_mocks "github.com/ShmelJUJ/software-engineering/pkg/kafka/mocks"
	logger_mocks "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	transactionID  = "test-transaction-id"
	paymentID      = "test-payment-id"
	retries        = 5
	timeout        = 50 * time.Millisecond
	failedTopic    = "failed"
	succeededTopic = "succeeded"
	monitorTopic   = "monitor.process"
)

func publisherHelper(t *testing.T) (*logger_mocks.MockLogger, *gateway_mocks.MockPaymentGateway, *kafka_mocks.MockPublisher) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	log := logger_mocks.NewMockLogger(mockCtrl)
	g := gateway_mocks.NewMockPaymentGateway(mockCtrl)
	p := kafka_mocks.NewMockPublisher(mockCtrl)

	return log, g, p
}

func TestStartWorker(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type args struct {
		ctx context.Context
		cfg *Config
	}

	testcases := []struct {
		name        string
		args        args
		mock        func(*logger_mocks.MockLogger, *gateway_mocks.MockPaymentGateway, *kafka_mocks.MockPublisher)
		expectedErr error
	}{
		{
			name: "Create payment error",
			args: args{
				ctx: ctx,
				cfg: &Config{},
			},
			mock: func(ml *logger_mocks.MockLogger, mpg *gateway_mocks.MockPaymentGateway, _ *kafka_mocks.MockPublisher) {
				mpg.EXPECT().TransactionID().Return(transactionID).Times(1)
				ml.EXPECT().Debug("Start payment worker", map[string]interface{}{
					"transaction_id": transactionID,
				})
				mpg.EXPECT().CreatePayment(ctx).Return(paymentID, &gateway.CreatePaymentError{})
			},
			expectedErr: &StartError{
				msg: "failed to create payment",
				err: &gateway.CreatePaymentError{},
			},
		},
		// {
		// 	name: "Transaction processing time has expired",
		// 	args: args{
		// 		ctx: ctx,
		// 		cfg: &Config{
		// 			PaymentProccessingTime: 500 * time.Millisecond,
		// 			FailedTransactionTopic: failedTopic,
		// 		},
		// 	},
		// 	mock: func(ml *logger_mocks.MockLogger, mpg *gateway_mocks.MockPaymentGateway, mp *kafka_mocks.MockPublisher) {
		// 		mpg.EXPECT().TransactionID().Return(transactionID).Times(4)
		// 		ml.EXPECT().Debug("Start payment worker", map[string]interface{}{
		// 			"transaction_id": transactionID,
		// 		}).Times(1)
		// 		mpg.EXPECT().CreatePayment(ctx).Return(paymentID, nil).Times(1)
		// 		ml.EXPECT().Debug("Payment worker start payment processing", map[string]interface{}{
		// 			"transaction_id": transactionID,
		// 			"payment_id":     paymentID,
		// 		}).Times(1)
		// 		mpg.EXPECT().Timeout().Return(450 * time.Millisecond).Times(1)
		// 		mpg.EXPECT().Retries().Return(1).Times(1)
		// 		mpg.EXPECT().CheckStatus(ctx, paymentID).Return(gateway.Pending, nil).Times(1)
		// 		ml.EXPECT().Debug("Payment worker waiting for a change in the transaction status", map[string]interface{}{
		// 			"transaction_id": transactionID,
		// 			"payment_id":     paymentID,
		// 			"retries":        0,
		// 		})
		// 		ml.EXPECT().Debug("Payment worker handle failed transaction", map[string]interface{}{
		// 			"transaction_id": transactionID,
		// 			"reason":         "Maximum transaction processing time has expired",
		// 		}).Times(1)
		// 		mp.EXPECT().Publish(failedTopic, gomock.Any()).Return(nil).Times(1)
		// 	},
		// 	expectedErr: nil,
		// },
		{
			name: "Cancelled transaction by payment gateway",
			args: args{
				ctx: ctx,
				cfg: &Config{
					PaymentProccessingTime: time.Minute,
					FailedTransactionTopic: failedTopic,
				},
			},
			mock: func(ml *logger_mocks.MockLogger, mpg *gateway_mocks.MockPaymentGateway, mp *kafka_mocks.MockPublisher) {
				mpg.EXPECT().TransactionID().Return(transactionID).Times(4)
				ml.EXPECT().Debug("Start payment worker", map[string]interface{}{
					"transaction_id": transactionID,
				}).Times(1)
				mpg.EXPECT().CreatePayment(ctx).Return(paymentID, nil).Times(1)
				ml.EXPECT().Debug("Payment worker start payment processing", map[string]interface{}{
					"transaction_id": transactionID,
					"payment_id":     paymentID,
				}).Times(1)
				mpg.EXPECT().Timeout().Return(timeout).Times(1)
				mpg.EXPECT().Retries().Return(retries).Times(1)
				mpg.EXPECT().CheckStatus(ctx, paymentID).Return(gateway.Cancelled, nil).Times(1)
				ml.EXPECT().Debug("Payment worker handle failed transaction", map[string]interface{}{
					"transaction_id": transactionID,
					"reason":         "Payment gateway cancelled the transaction",
				}).Times(1)
				mp.EXPECT().Publish(monitorTopic, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Succeeded transaction",
			args: args{
				ctx: ctx,
				cfg: &Config{
					PaymentProccessingTime:    time.Minute,
					SucceededTransactionTopic: succeededTopic,
				},
			},
			mock: func(ml *logger_mocks.MockLogger, mpg *gateway_mocks.MockPaymentGateway, mp *kafka_mocks.MockPublisher) {
				mpg.EXPECT().TransactionID().Return(transactionID).Times(4)
				ml.EXPECT().Debug("Start payment worker", map[string]interface{}{
					"transaction_id": transactionID,
				}).Times(1)
				mpg.EXPECT().CreatePayment(ctx).Return(paymentID, nil).Times(1)
				ml.EXPECT().Debug("Payment worker start payment processing", map[string]interface{}{
					"transaction_id": transactionID,
					"payment_id":     paymentID,
				}).Times(1)
				mpg.EXPECT().Timeout().Return(timeout).Times(1)
				mpg.EXPECT().Retries().Return(retries).Times(1)
				mpg.EXPECT().CheckStatus(ctx, paymentID).Return(gateway.Succeeded, nil).Times(1)
				ml.EXPECT().Debug("Payment worker handle succeeded transaction", map[string]interface{}{
					"transaction_id": transactionID,
				}).Times(1)
				mp.EXPECT().Publish(monitorTopic, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			mockLog, mockGateway, mockPublisher := publisherHelper(t)
			testcase.mock(mockLog, mockGateway, mockPublisher)

			worker, err := NewWorker(testcase.args.cfg, mockLog, mockGateway, mockPublisher)
			assert.NoError(t, err)

			err = worker.Start(testcase.args.ctx)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}
