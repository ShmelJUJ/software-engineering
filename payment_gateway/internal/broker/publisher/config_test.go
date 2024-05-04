package publisher

import (
	"testing"
	"time"

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
				PaymentProccessingTime: time.Minute,
				FailedTransactionTopic: "transaction.failed2",
			},
			expectedCfg: &Config{
				MonitorProcessTopic:       defaultMonitorProcessTopic,
				PaymentProccessingTime:    time.Minute,
				SucceededTransactionTopic: defaultSucceededTransactionTopic,
				FailedTransactionTopic:    "transaction.failed2",
			},
		},
		{
			name: "With empty config",
			cfg:  &Config{},
			expectedCfg: &Config{
				MonitorProcessTopic:       defaultMonitorProcessTopic,
				PaymentProccessingTime:    defaultPaymentProccessingTime,
				SucceededTransactionTopic: defaultSucceededTransactionTopic,
				FailedTransactionTopic:    defaultFailedTransactionTopic,
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
