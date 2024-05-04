package subscriber

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				PoolCfg: &PoolConfig{
					MinWorkers:    10,
					TasksCapacity: 10000,
				},
				ProcessedTransactionTopic: "transaction.processed2",
			},
			expectedCfg: &Config{
				PoolCfg: &PoolConfig{
					IdleTimeout:   defaultIdleTimeout,
					MinWorkers:    10,
					NumWorkers:    defaultNumWorkers,
					TasksCapacity: 10000,
				},
				ProcessedTransactionTopic: "transaction.processed2",
				CancelledTransactionTopic: defaultCancelledTransactionTopic,
			},
		},
		{
			name: "With empty config",
			cfg:  &Config{},
			expectedCfg: &Config{
				PoolCfg: &PoolConfig{
					IdleTimeout:   defaultIdleTimeout,
					MinWorkers:    defaultMinWorkers,
					NumWorkers:    defaultNumWorkers,
					TasksCapacity: defaultTasksCapacity,
				},
				ProcessedTransactionTopic: defaultProcessedTransactionTopic,
				CancelledTransactionTopic: defaultCancelledTransactionTopic,
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
