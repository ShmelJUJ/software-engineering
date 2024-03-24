package logger_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	debugLevel   = "debug"
	traceLevel   = "trace"
	panicLevel   = "panic"
	fatalLevel   = "fatal"
	errorLevel   = "error"
	warnLevel    = "warn"
	warningLevel = "warning"
	infoLevel    = "info"
	invalidLevel = "ttt"
)

func TestNewLogrusLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		level string
	}

	unknownLevelErr := fmt.Errorf("failed to parse logrus level: %w", errors.New("not a valid logrus Level: \"ttt\""))

	testcases := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "success with debug level",
			args: args{
				level: debugLevel,
			},
			expectedErr: nil,
		},
		{
			name: "success with info level",
			args: args{
				level: infoLevel,
			},
			expectedErr: nil,
		},
		{
			name: "success with warn level",
			args: args{
				level: warnLevel,
			},
			expectedErr: nil,
		},
		{
			name: "success with warning level",
			args: args{
				level: warningLevel,
			},
			expectedErr: nil,
		},
		{
			name: "success with error level",
			args: args{
				level: errorLevel,
			},
			expectedErr: nil,
		},
		{
			name: "success with fatal level",
			args: args{
				level: fatalLevel,
			},
			expectedErr: nil,
		},
		{
			name: "success with panic level",
			args: args{
				level: panicLevel,
			},
			expectedErr: nil,
		},
		{
			name: "success with trace level",
			args: args{
				level: traceLevel,
			},
			expectedErr: nil,
		},
		{
			name: "failed with unknown level",
			args: args{
				level: invalidLevel,
			},
			expectedErr: unknownLevelErr,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			_, err := logger.NewLogrusLogger(testcase.args.level)

			if testcase.expectedErr != nil {
				assert.Errorf(t, err, testcase.expectedErr.Error())
			} else {
				assert.Equal(t, testcase.expectedErr, err)
			}
		})
	}
}
