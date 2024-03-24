package kafka_test

import (
	"testing"

	"github.com/ShmelJUJ/software-engineering/pkg/kafka"
	"github.com/stretchr/testify/assert"
)

func TestNewBrokerRouter(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "successfully create new broker router",
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			_, err := kafka.NewBrokerRouter()

			if testcase.expectedErr != nil {
				assert.Errorf(t, err, testcase.expectedErr.Error())
			} else {
				assert.Equal(t, testcase.expectedErr, err)
			}
		})
	}
}
