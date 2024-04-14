package algorand

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
				AlgodAddress: "test-algod-address",
				Retries:      81,
			},
			expectedCfg: &Config{
				AlgodAddress:           "test-algod-address",
				AlgodToken:             defaultAlgodToken,
				ConfirmationWaitRounds: defaultConfirmationWaitRounds,
				Timeout:                defaultTimeout,
				Retries:                81,
				IsTest:                 defaultIsTest,
			},
		},
		{
			name: "With empty config",
			cfg:  &Config{},
			expectedCfg: &Config{
				AlgodAddress:           defaultAlgodAddress,
				AlgodToken:             defaultAlgodToken,
				ConfirmationWaitRounds: defaultConfirmationWaitRounds,
				Timeout:                defaultTimeout,
				Retries:                defaultRetries,
				IsTest:                 defaultIsTest,
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
