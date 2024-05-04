package publisher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTransactionProcessedTopic = "transaction.processed"
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
				ProcessedTransactionTopic: testTransactionProcessedTopic,
			},
			expectedCfg: &Config{
				ProcessedTransactionTopic: testTransactionProcessedTopic,
				ProcessMonitorTopic:       defaultProcessMonitorTopic,
			},
		},
		{
			name: "With empty config",
			cfg:  &Config{},
			expectedCfg: &Config{
				ProcessedTransactionTopic: defaultProcessedTransactionTopic,
				ProcessMonitorTopic:       defaultProcessMonitorTopic,
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
