package kafka_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	"github.com/stretchr/testify/assert"
)

var (
	errMissingBrokers   = errors.New("missing brokers")
	errMissingMarshaler = errors.New("missing marshaler")

	validBrokers = []string{"kafka:9092"}
)

func TestNewPublisher(t *testing.T) {
	t.Parallel()

	type args struct {
		brokers []string
		opts    []kafka.PublisherOption
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
			expectedErr: fmt.Errorf("failed to create a new publisher: %w", errMissingBrokers),
		},
		{
			name: "failed with nil marhaler",
			args: args{
				brokers: validBrokers,
				opts:    []kafka.PublisherOption{kafka.WithMarshaler(nil)},
			},
			expectedErr: fmt.Errorf("failed to create a new publisher: %w", errMissingMarshaler),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			_, err := kafka.NewPublisher(testcase.args.brokers, testcase.args.opts...)

			if testcase.expectedErr != nil {
				assert.Errorf(t, err, testcase.expectedErr.Error())
			} else {
				assert.Equal(t, testcase.expectedErr, err)
			}
		})
	}
}
