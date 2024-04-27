package publisher

import (
	"testing"

	mock_publisher "github.com/ShmelJUJ/software-engineering/pkg/kafka/mocks"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	mock_logger "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher/dto"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func transactionPublisherHelper(t *testing.T) (*mock_logger.MockLogger, *mock_publisher.MockPublisher) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	l := mock_logger.NewMockLogger(mockCtrl)
	publisher := mock_publisher.NewMockPublisher(mockCtrl)

	return l, publisher
}

func TestNewTransactionPublisher(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg *Config
		log logger.Logger
		pub message.Publisher
	}

	log, pub := transactionPublisherHelper(t)

	testcases := []struct {
		name                         string
		args                         args
		expectedTransactionPublisher *transactionPublisher
		expectedErr                  error
	}{
		{
			name: "Successfully create new transaction publisher",
			args: args{
				cfg: &Config{},
				log: log,
				pub: pub,
			},
			expectedTransactionPublisher: &transactionPublisher{
				cfg: getDefaultConfig(),
				log: log,
				pub: pub,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualTransactionPublisher, err := NewTransactionPublisher(
				testcase.args.cfg,
				testcase.args.log,
				testcase.args.pub,
			)

			assert.Equal(t, testcase.expectedTransactionPublisher, actualTransactionPublisher)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestPublishSucceededTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		transaction *dto.ProcessedTransaction
	}

	processedTransaction := &dto.ProcessedTransaction{}

	someErr := NewPublishSucceededTransactionError("test err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_publisher.MockPublisher)
		expectedErr error
	}{
		{
			name: "Successfully publish succeeded transaction",
			args: args{
				transaction: processedTransaction,
			},
			mock: func(ml *mock_logger.MockLogger, mp *mock_publisher.MockPublisher) {
				ml.EXPECT().Debug("Start publish succeeded transaction", map[string]interface{}{
					"transaction": processedTransaction,
				})
				mp.EXPECT().Publish(testTransactionProcessedTopic, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to publish succeeded transaction",
			args: args{
				transaction: processedTransaction,
			},
			mock: func(ml *mock_logger.MockLogger, mp *mock_publisher.MockPublisher) {
				ml.EXPECT().Debug("Start publish succeeded transaction", map[string]interface{}{
					"transaction": processedTransaction,
				})
				mp.EXPECT().Publish(testTransactionProcessedTopic, gomock.Any()).Return(someErr).Times(1)
			},
			expectedErr: NewPublishSucceededTransactionError("failed to publish succeded transaction", someErr),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			log, pub := transactionPublisherHelper(t)

			testcase.mock(log, pub)

			transactionPublisher, err := NewTransactionPublisher(
				&Config{
					ProcessedTransactionTopic: testTransactionProcessedTopic,
				},
				log,
				pub,
			)
			assert.NoError(t, err)

			err = transactionPublisher.PublishSucceededTransaction(testcase.args.transaction)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}
