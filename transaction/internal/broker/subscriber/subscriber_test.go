package subscriber

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	mock_subscriber "github.com/ShmelJUJ/software-engineering/pkg/kafka/mocks"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	mock_logger "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/subscriber/dto"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/repository"
	mock_repo "github.com/ShmelJUJ/software-engineering/transaction/internal/repository/mocks"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	testTransactionID = "test-id"
	testReason        = "test-reason"
)

func transactionSubscriberHelper(t *testing.T) (*mock_logger.MockLogger, *mock_subscriber.MockSubscriber, *mock_repo.MockTransactionRepo) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	l := mock_logger.NewMockLogger(mockCtrl)
	sub := mock_subscriber.NewMockSubscriber(mockCtrl)
	repo := mock_repo.NewMockTransactionRepo(mockCtrl)

	return l, sub, repo
}

func TestNewTransactionPublisher(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg             *Config
		log             logger.Logger
		sub             message.Subscriber
		router          *message.Router
		transactionRepo repository.TransactionRepo
	}

	log, sub, repo := transactionSubscriberHelper(t)

	router, err := kafka.NewBrokerRouter()
	assert.NoError(t, err)

	testcases := []struct {
		name                          string
		args                          args
		expectedTransactionSubscriber *TransactionSubscriber
		expectedErr                   error
	}{
		{
			name: "Successfully create new transaction subscriber",
			args: args{
				cfg:             &Config{},
				log:             log,
				sub:             sub,
				router:          router,
				transactionRepo: repo,
			},
			expectedTransactionSubscriber: &TransactionSubscriber{
				cfg:             getDefaultConfig(),
				log:             log,
				sub:             sub,
				router:          router,
				transactionRepo: repo,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualTransactionSubscriber, err := NewTransactionSubscriber(
				testcase.args.cfg,
				testcase.args.log,
				testcase.args.sub,
				testcase.args.router,
				testcase.args.transactionRepo,
			)

			assert.Equal(t, testcase.expectedTransactionSubscriber, actualTransactionSubscriber)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestHandleSucceededTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		msg *message.Message
	}

	ctx := context.Background()

	router, err := kafka.NewBrokerRouter()
	assert.NoError(t, err)

	succeededTransaction := &dto.SucceededTransaction{
		TransactionID: testTransactionID,
	}

	succeededTransactionData, err := json.Marshal(succeededTransaction)
	assert.NoError(t, err)

	someErr := repository.NewChangeTransactionStatusError("test-err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedErr error
	}{
		{
			name: "Successfully handle succeeded transaction",
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), succeededTransactionData),
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Start handle succeeded transaction", map[string]interface{}{
					"transaction_id": succeededTransaction.TransactionID,
				})
				mtr.EXPECT().ChangeTransactionStatus(ctx, succeededTransaction.TransactionID, model.Succeeded).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to handle succeeded transaction",
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), succeededTransactionData),
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Start handle succeeded transaction", map[string]interface{}{
					"transaction_id": succeededTransaction.TransactionID,
				})
				mtr.EXPECT().ChangeTransactionStatus(ctx, succeededTransaction.TransactionID, model.Succeeded).Return(someErr).Times(1)
				ml.EXPECT().Error("failed to change transaction status", map[string]interface{}{
					"error":          someErr,
					"status":         model.Succeeded,
					"transaction_id": succeededTransaction.TransactionID,
				})
			},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			log, sub, repo := transactionSubscriberHelper(t)

			testcase.mock(log, repo)

			transactionSubscriber, err := NewTransactionSubscriber(
				&Config{},
				log,
				sub,
				router,
				repo,
			)
			assert.NoError(t, err)

			err = transactionSubscriber.handleSucceededTransaction(testcase.args.msg)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestHandleFailedTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		msg *message.Message
	}

	ctx := context.Background()

	router, err := kafka.NewBrokerRouter()
	assert.NoError(t, err)

	failedTransaction := &dto.FailedTransaction{
		TransactionID: testTransactionID,
		Reason:        testReason,
	}

	failedTransactionData, err := json.Marshal(failedTransaction)
	assert.NoError(t, err)

	someErr := repository.NewCancelTransactionError("test-err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_repo.MockTransactionRepo)
		expectedErr error
	}{
		{
			name: "Successfully handle failed transaction",
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), failedTransactionData),
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Start handle failed transaction", map[string]interface{}{
					"transaction_id": failedTransaction.TransactionID,
				})
				mtr.EXPECT().CancelTransaction(ctx, failedTransaction.TransactionID, failedTransaction.Reason).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to handle succeeded transaction",
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), failedTransactionData),
			},
			mock: func(ml *mock_logger.MockLogger, mtr *mock_repo.MockTransactionRepo) {
				ml.EXPECT().Debug("Start handle failed transaction", map[string]interface{}{
					"transaction_id": failedTransaction.TransactionID,
				})
				mtr.EXPECT().CancelTransaction(ctx, failedTransaction.TransactionID, failedTransaction.Reason).Return(someErr).Times(1)
				ml.EXPECT().Error("failed to cancel transaction", map[string]interface{}{
					"error":          someErr,
					"reason":         failedTransaction.Reason,
					"transaction_id": failedTransaction.TransactionID,
				})
			},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			log, sub, repo := transactionSubscriberHelper(t)

			testcase.mock(log, repo)

			transactionSubscriber, err := NewTransactionSubscriber(
				&Config{},
				log,
				sub,
				router,
				repo,
			)
			assert.NoError(t, err)

			err = transactionSubscriber.handleFailedTransaction(testcase.args.msg)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}
