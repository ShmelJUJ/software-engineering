package subscriber

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ShmelJUJ/software-engineering/monitor/internal/broker/publisher"
	mock_monitor_publisher "github.com/ShmelJUJ/software-engineering/monitor/internal/broker/publisher/mocks"
	"github.com/ShmelJUJ/software-engineering/monitor/internal/broker/subscriber/dto"
	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	mock_subscriber "github.com/ShmelJUJ/software-engineering/pkg/kafka/mocks"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	mock_logger "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func monitorSubscriberHelper(t *testing.T) (
	*mock_logger.MockLogger,
	*mock_subscriber.MockSubscriber,
	*mock_monitor_publisher.MockMonitorPublisher,
) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	l := mock_logger.NewMockLogger(mockCtrl)
	sub := mock_subscriber.NewMockSubscriber(mockCtrl)
	monitorPub := mock_monitor_publisher.NewMockMonitorPublisher(mockCtrl)

	return l, sub, monitorPub
}

func TestNewMonitorSubscriber(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg        *Config
		log        logger.Logger
		sub        message.Subscriber
		router     *message.Router
		monitorPub publisher.MonitorPublisher
	}

	log, sub, monitorPub := monitorSubscriberHelper(t)

	router, err := kafka.NewBrokerRouter()
	assert.NoError(t, err)

	testcases := []struct {
		name                      string
		args                      args
		expectedMonitorSubscriber *MonitorSubscriber
		expectedErr               error
	}{
		{
			name: "Successfully create new monitor subscriber",
			args: args{
				cfg:        &Config{},
				log:        log,
				sub:        sub,
				router:     router,
				monitorPub: monitorPub,
			},
			expectedMonitorSubscriber: &MonitorSubscriber{
				cfg:        getDefaultConfig(),
				log:        log,
				sub:        sub,
				router:     router,
				monitorPub: monitorPub,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualMonitorSubscriber, err := NewMonitorSubscriber(
				testcase.args.cfg,
				testcase.args.log,
				testcase.args.sub,
				testcase.args.router,
				testcase.args.monitorPub,
			)

			assert.Equal(t, testcase.expectedMonitorSubscriber, actualMonitorSubscriber)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestHandleProcess(t *testing.T) {
	t.Parallel()

	type args struct {
		msg *message.Message
	}

	router, err := kafka.NewBrokerRouter()
	assert.NoError(t, err)

	processDTO := &dto.Process{
		From:    fromTransaction,
		ToTopic: toProcessedTransactionTopic,
		Payload: "test-payload",
	}

	processDTOData, err := json.Marshal(processDTO)
	assert.NoError(t, err)

	someErr := errors.New("test-err")

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_monitor_publisher.MockMonitorPublisher)
		expectedErr error
	}{
		{
			name: "Successfully handle process message",
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), processDTOData),
			},
			mock: func(ml *mock_logger.MockLogger, mmp *mock_monitor_publisher.MockMonitorPublisher) {
				ml.EXPECT().Debug("Start handle process", map[string]interface{}{
					"from":     processDTO.From,
					"to_topic": processDTO.ToTopic,
					"payload":  processDTO.Payload,
				})
				mmp.EXPECT().PublishProcess(processDTO.ToTopic, processDTO.Payload).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to handle process message",
			args: args{
				msg: message.NewMessage(watermill.NewUUID(), processDTOData),
			},
			mock: func(ml *mock_logger.MockLogger, mmp *mock_monitor_publisher.MockMonitorPublisher) {
				ml.EXPECT().Debug("Start handle process", map[string]interface{}{
					"from":     processDTO.From,
					"to_topic": processDTO.ToTopic,
					"payload":  processDTO.Payload,
				})
				mmp.EXPECT().PublishProcess(processDTO.ToTopic, processDTO.Payload).Return(someErr).Times(1)
				ml.EXPECT().Error("failed to publish process", map[string]interface{}{
					"error":    someErr,
					"from":     processDTO.From,
					"to_topic": processDTO.ToTopic,
					"payload":  processDTO.Payload,
				})
			},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			log, sub, monitorPub := monitorSubscriberHelper(t)

			testcase.mock(log, monitorPub)

			monitorSubscriber, err := NewMonitorSubscriber(
				&Config{},
				log,
				sub,
				router,
				monitorPub,
			)
			assert.NoError(t, err)

			err = monitorSubscriber.handleProcess(testcase.args.msg)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}

func TestVerify(t *testing.T) {
	t.Parallel()

	type args struct {
		from, toTopic string
	}

	testcases := []struct {
		name        string
		args        args
		expectedVal bool
	}{
		{
			name: "Successful verify from transaction to transaction.processed topic",
			args: args{
				from:    fromTransaction,
				toTopic: toProcessedTransactionTopic,
			},
			expectedVal: true,
		},
		{
			name: "Successful verify from payment_gateway to transaction.failed topic",
			args: args{
				from:    fromPaymentGateway,
				toTopic: toFailedTransactionTopic,
			},
			expectedVal: true,
		},
		{
			name: "Successful verify from payment_gateway to transaction.succeeded topic",
			args: args{
				from:    fromPaymentGateway,
				toTopic: toSucceededTransactionTopic,
			},
			expectedVal: true,
		},
		{
			name: "failed to verify from someService to transaction.failed topic",
			args: args{
				from:    "someService",
				toTopic: toFailedTransactionTopic,
			},
			expectedVal: false,
		},
		{
			name: "failed to verify from payment_gateway to transaction.someTopic topic",
			args: args{
				from:    fromPaymentGateway,
				toTopic: "transaction.someTopic",
			},
			expectedVal: false,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			val := verify(testcase.args.from, testcase.args.toTopic)
			assert.Equal(t, testcase.expectedVal, val)
		})
	}
}
