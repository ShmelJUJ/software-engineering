package kafka_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	"github.com/stretchr/testify/assert"
)

var errMissingUnmarhaler = errors.New("missing unmarshaler")

func TestNewSubscriber(t *testing.T) {
	t.Parallel()

	type args struct {
		brokers []string
		opts    []kafka.SubscriberOption
	}

	testcases := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "failed with empty brokers",
			args: args{
				brokers: make([]string, 0),
			},
			expectedErr: fmt.Errorf("failed to create a new subscriber: %w", errMissingBrokers),
		},
		{
			name: "failed with nil unmarhaler",
			args: args{
				brokers: validBrokers,
				opts:    []kafka.SubscriberOption{kafka.WithUnmarshaler(nil)},
			},
			expectedErr: fmt.Errorf("failed to create a new subscriber: %w", errMissingUnmarhaler),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			_, err := kafka.NewSubscriber(testcase.args.brokers, testcase.args.opts...)

			if testcase.expectedErr != nil {
				assert.Errorf(t, err, testcase.expectedErr.Error())
			} else {
				assert.Equal(t, testcase.expectedErr, err)
			}
		})
	}
}
