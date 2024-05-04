package subscriber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFailedTransactionTopic    = "transaction.failed"
	testSucceededTransactionTopic = "transaction.succeeded"
)

func TestMergeWithDefault(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		cfg         *Config
		expectedCfg *Config
		expectedErr error
	}{
		{
			name: "With some config",
			cfg: &Config{
				FailedTransactionTopic: testFailedTransactionTopic,
			},
			expectedCfg: &Config{
				FailedTransactionTopic:    testFailedTransactionTopic,
				SucceededTransactionTopic: defaultSucceededTransactionTopic,
			},
		},
		{
			name: "With empty config",
			cfg:  &Config{},
			expectedCfg: &Config{
				FailedTransactionTopic:    defaultFailedTransactionTopic,
				SucceededTransactionTopic: defaultSucceededTransactionTopic,
			},
			expectedErr: nil,
		},
		{
			name:        "With nil config",
			cfg:         nil,
			expectedCfg: nil,
			expectedErr: ErrNilConfig,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualCfg, err := mergeWithDefault(testcase.cfg)

			assert.Equal(t, testcase.expectedCfg, actualCfg)
			assert.Equal(t, testcase.expectedErr, err)
		})
	}
}
