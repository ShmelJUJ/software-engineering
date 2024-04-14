package subscriber

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/publisher"
	worker_mocks "github.com/ShmelJUJ/software-engineering/payment_gateway/internal/broker/publisher/mocks"
	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway/algorand"
	kafka_mocks "github.com/ShmelJUJ/software-engineering/pkg/kafka/mocks"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	logger_mocks "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alitto/pond"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	transactionID  = "test-transaction-id"
	processedTopic = "processed"
	cancelledTopic = "cancelled"
)

func subscriberHelper(t *testing.T) (
	*logger_mocks.MockLogger,
	*kafka_mocks.MockSubscriber,
	*kafka_mocks.MockPublisher,
	*worker_mocks.MockPaymentWorker,
) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	log := logger_mocks.NewMockLogger(mockCtrl)
	s := kafka_mocks.NewMockSubscriber(mockCtrl)
	p := kafka_mocks.NewMockPublisher(mockCtrl)
	w := worker_mocks.NewMockPaymentWorker(mockCtrl)

	return log, s, p, w
}

func TestNewTransactionSubscriber(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	log := logger_mocks.NewMockLogger(mockCtrl)
	p := kafka_mocks.NewMockPublisher(mockCtrl)
	s := kafka_mocks.NewMockSubscriber(mockCtrl)

	type args struct {
		cfg          *Config
		logger       logger.Logger
		router       *message.Router
		subscriber   message.Subscriber
		publisher    message.Publisher
		publisherCfg *publisher.Config
		algorandCfg  *algorand.Config
	}

	testcases := []struct {
		name                          string
		args                          args
		expectedTransactionSubscriber *TransactionSubscriber
		expectedErr                   error
	}{
		{
			name: "With default config",
			args: args{
				cfg:          &Config{},
				logger:       log,
				router:       &message.Router{},
				subscriber:   s,
				publisher:    p,
				publisherCfg: &publisher.Config{},
				algorandCfg:  &algorand.Config{},
			},
			expectedTransactionSubscriber: &TransactionSubscriber{
				paymentWorkers: sync.Map{},
				router:         &message.Router{},
				sub:            s,
				pub:            p,
				publisherCfg:   &publisher.Config{},
				pool: pond.New(
					defaultNumWorkers,
					defaultTasksCapacity,
					pond.IdleTimeout(defaultIdleTimeout),
					pond.MinWorkers(defaultMinWorkers),
				),
				cfg:         getDefaultConfig(),
				algorandCfg: &algorand.Config{},
				log:         log,
			},
			expectedErr: nil,
		},
		{
			name: "With some config",
			args: args{
				cfg: &Config{
					PoolCfg: &PoolConfig{
						NumWorkers: 10000,
					},
					ProcessedTransactionTopic: processedTopic,
				},
				logger:       log,
				router:       &message.Router{},
				subscriber:   s,
				publisher:    p,
				publisherCfg: &publisher.Config{},
				algorandCfg:  &algorand.Config{},
			},
			expectedTransactionSubscriber: &TransactionSubscriber{
				paymentWorkers: sync.Map{},
				router:         &message.Router{},
				sub:            s,
				pub:            p,
				publisherCfg:   &publisher.Config{},
				pool: pond.New(
					10000,
					defaultTasksCapacity,
					pond.IdleTimeout(defaultIdleTimeout),
					pond.MinWorkers(defaultMinWorkers),
				),
				cfg: &Config{
					PoolCfg: &PoolConfig{
						IdleTimeout:   defaultIdleTimeout,
						MinWorkers:    defaultMinWorkers,
						NumWorkers:    10000,
						TasksCapacity: defaultTasksCapacity,
					},
					ProcessedTransactionTopic: processedTopic,
					CancelledTransactionTopic: defaultCancelledTransactionTopic,
				},
				algorandCfg: &algorand.Config{},
				log:         log,
			},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualTransactionSubscriber, err := NewTransactionSubscriber(
				testcase.args.cfg,
				testcase.args.logger,
				testcase.args.router,
				testcase.args.subscriber,
				testcase.args.publisher,
				testcase.args.publisherCfg,
				testcase.args.algorandCfg,
			)

			assert.Equal(t, testcase.expectedTransactionSubscriber.cfg, actualTransactionSubscriber.cfg)
			assert.Equal(t, testcase.expectedTransactionSubscriber.log, actualTransactionSubscriber.log)
			assert.Equal(t, testcase.expectedTransactionSubscriber.pub, actualTransactionSubscriber.pub)
			assert.Equal(t, testcase.expectedTransactionSubscriber.publisherCfg, actualTransactionSubscriber.publisherCfg)
			assert.Equal(t, testcase.expectedTransactionSubscriber.router, actualTransactionSubscriber.router)
			assert.Equal(t, testcase.expectedTransactionSubscriber.sub, actualTransactionSubscriber.sub)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestHandleProcessedTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		payload message.Payload
	}

	testcases := []struct {
		name string
		args args
		mock func(*logger_mocks.MockLogger)
	}{
		{
			name: "Get payment gateway error",
			args: args{
				payload: []byte(`{"transaction":{}}`),
			},
			mock: func(ml *logger_mocks.MockLogger) {
				ml.EXPECT().Debug("Start handle processed transaction", map[string]interface{}{
					"transaction_id": "",
				}).Times(1)
				ml.EXPECT().Error("failed to get payment gateway", map[string]interface{}{
					"transaction_id": "",
					"error":          fmt.Errorf("cannot handle %s payment gateway", ""),
				})
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			mockLog, mockSubscriber, mockPublisher, _ := subscriberHelper(t)
			testcase.mock(mockLog)

			transactionSubscriber, err := NewTransactionSubscriber(
				&Config{},
				mockLog,
				&message.Router{},
				mockSubscriber,
				mockPublisher,
				&publisher.Config{},
				&algorand.Config{},
			)
			assert.NoError(t, err)

			err = transactionSubscriber.handleProcessedTransaction(
				message.NewMessage(watermill.NewUUID(), testcase.args.payload),
			)
			assert.NoError(t, err)
		})
	}
}

func lenSyncMap(m *sync.Map) int {
	var i int

	m.Range(func(_, _ interface{}) bool {
		i++
		return true
	})

	return i
}

func TestHandleCancelledTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		payload message.Payload
	}

	testcases := []struct {
		name               string
		args               args
		mock               func(*logger_mocks.MockLogger, *worker_mocks.MockPaymentWorker)
		preAction          func(*TransactionSubscriber, string, publisher.PaymentWorker)
		expectedNumWorkers int
	}{
		{
			name: "Successfully cancelled transaction",
			args: args{
				payload: []byte(`{"transaction_id":"test-transaction-id"}`),
			},
			preAction: func(transactionSubscriber *TransactionSubscriber, transactionID string, worker publisher.PaymentWorker) {
				transactionSubscriber.paymentWorkers.Store(transactionID, worker)
			},
			mock: func(ml *logger_mocks.MockLogger, mw *worker_mocks.MockPaymentWorker) {
				ml.EXPECT().Debug("Start handle cancelled transaction", map[string]interface{}{
					"transaction_id": transactionID,
				}).Times(1)
				mw.EXPECT().Stop(publisher.CancelledTransaction).Times(1)
			},
			expectedNumWorkers: 0,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			mockLog, mockSubscriber, mockPublisher, mockWorker := subscriberHelper(t)
			testcase.mock(mockLog, mockWorker)

			transactionSubscriber, err := NewTransactionSubscriber(
				&Config{},
				mockLog,
				&message.Router{},
				mockSubscriber,
				mockPublisher,
				&publisher.Config{},
				&algorand.Config{},
			)
			assert.NoError(t, err)

			testcase.preAction(transactionSubscriber, transactionID, mockWorker)

			err = transactionSubscriber.handleCancelledTransaction(
				message.NewMessage(watermill.NewUUID(), testcase.args.payload),
			)

			assert.NoError(t, err)
			assert.Equal(t, testcase.expectedNumWorkers, lenSyncMap(&transactionSubscriber.paymentWorkers))
		})
	}
}
