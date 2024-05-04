package publisher

import (
	"testing"

	mock_publisher "github.com/ShmelJUJ/software-engineering/pkg/kafka/mocks"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	mock_logger "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func monitorPublisherHelper(t *testing.T) (*mock_logger.MockLogger, *mock_publisher.MockPublisher) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	l := mock_logger.NewMockLogger(mockCtrl)
	pub := mock_publisher.NewMockPublisher(mockCtrl)

	return l, pub
}

func TestNewMonitorPublisher(t *testing.T) {
	t.Parallel()

	type args struct {
		log logger.Logger
		pub message.Publisher
	}

	log, pub := monitorPublisherHelper(t)

	testcases := []struct {
		name                     string
		args                     args
		expectedMonitorPublisher *monitorPublisher
	}{
		{
			name: "Successfully create new monitor publisher",
			args: args{
				log: log,
				pub: pub,
			},
			expectedMonitorPublisher: &monitorPublisher{
				log: log,
				pub: pub,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualMonitorPublisher := NewMonitorPublisher(
				testcase.args.log,
				testcase.args.pub,
			)

			assert.Equal(t, testcase.expectedMonitorPublisher, actualMonitorPublisher)
		})
	}
}

func TestPublishProcess(t *testing.T) {
	t.Parallel()

	type args struct {
		toTopic string
		payload any
	}

	payload := "test"
	toTopic := "test-topic"

	someErr := NewPublishProcessError("test err", nil)

	testcases := []struct {
		name        string
		args        args
		mock        func(*mock_logger.MockLogger, *mock_publisher.MockPublisher)
		expectedErr error
	}{
		{
			name: "Successfully publish process",
			args: args{
				toTopic: toTopic,
				payload: payload,
			},
			mock: func(ml *mock_logger.MockLogger, mp *mock_publisher.MockPublisher) {
				ml.EXPECT().Debug("Publish process", map[string]interface{}{
					"to_topic": toTopic,
					"payload":  payload,
				})
				mp.EXPECT().Publish(toTopic, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Failed to publish process",
			args: args{
				toTopic: toTopic,
				payload: payload,
			},
			mock: func(ml *mock_logger.MockLogger, mp *mock_publisher.MockPublisher) {
				ml.EXPECT().Debug("Publish process", map[string]interface{}{
					"to_topic": toTopic,
					"payload":  payload,
				})
				mp.EXPECT().Publish(toTopic, gomock.Any()).Return(someErr).Times(1)
			},
			expectedErr: NewPublishProcessError("failed to publish process message", someErr),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			log, pub := monitorPublisherHelper(t)

			testcase.mock(log, pub)

			monitorPublisher := NewMonitorPublisher(
				log,
				pub,
			)

			err := monitorPublisher.PublishProcess(testcase.args.toTopic, testcase.args.payload)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}
