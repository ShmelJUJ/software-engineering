package handler

import (
	"fmt"
	"testing"

	"github.com/ShmelJUJ/software-engineering/monitor/internal/generated/models"
	apiMonitor "github.com/ShmelJUJ/software-engineering/monitor/internal/generated/restapi/operations/monitor"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	mock_logger "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func monitorHandlerHelper(t *testing.T) *mock_logger.MockLogger {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	log := mock_logger.NewMockLogger(mockCtrl)

	return log
}

func TestNewMonitorHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		log logger.Logger
	}

	log := monitorHandlerHelper(t)

	testcases := []struct {
		name                   string
		args                   args
		expectedMonitorHandler *MonitorHandler
	}{
		{
			name: "Successfully create new monitor handler",
			args: args{
				log: log,
			},
			expectedMonitorHandler: &MonitorHandler{
				log: log,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualMonitorHandler := NewMonitorHandler(
				testcase.args.log,
			)

			assert.Equal(t, testcase.expectedMonitorHandler, actualMonitorHandler)
		})
	}
}

func TestProcessHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		params apiMonitor.ProcessParams
	}

	testPaymentGatewayService := paymentGatewayService
	testUserService := userService
	testGetUserMethod := getUserMethod
	testPayload := "test-payload"
	testUnknownService := "test-service"

	testcases := []struct {
		name             string
		args             args
		mock             func(*mock_logger.MockLogger)
		expectedResponse middleware.Responder
	}{
		{
			name: "Successfully process request",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testPaymentGatewayService,
						To:      &testUserService,
						Method:  &testGetUserMethod,
						Payload: testPayload,
					},
				},
			},
			mock: func(ml *mock_logger.MockLogger) {
				ml.EXPECT().Debug("Process handler", map[string]interface{}{
					"from":    testPaymentGatewayService,
					"to":      testUserService,
					"method":  testGetUserMethod,
					"payload": testPayload,
				})
			},
			expectedResponse: apiMonitor.NewProcessOK(),
		},
		{
			name: "Failed to verify process request",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testUnknownService,
						To:      &testUserService,
						Method:  &testGetUserMethod,
						Payload: testPayload,
					},
				},
			},
			mock: func(ml *mock_logger.MockLogger) {
				ml.EXPECT().Debug("Process handler", map[string]interface{}{
					"from":    testUnknownService,
					"to":      testUserService,
					"method":  testGetUserMethod,
					"payload": testPayload,
				})
			},
			expectedResponse: apiMonitor.NewProcessForbidden().
				WithPayload(&models.ErrorResponse{
					Code:    int32(apiMonitor.ProcessForbiddenCode),
					Message: fmt.Sprintf("cannot process request from %s to %s with method %s", testUnknownService, testUserService, testGetUserMethod),
				}),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			log := monitorHandlerHelper(t)

			testcase.mock(log)

			actualMonitorHandler := NewMonitorHandler(
				log,
			)

			actualResponse := actualMonitorHandler.ProcessHandler(testcase.args.params)

			assert.Equal(t, testcase.expectedResponse, actualResponse)
		})
	}
}

func TestVerify(t *testing.T) {
	t.Parallel()

	type args struct {
		from, to, method string
	}

	unknownService := "test-service"
	unknownMethod := "test-method"

	testcases := []struct {
		name        string
		args        args
		expectedVal bool
	}{
		{
			name: "Successful verify from payment gateway to user service with getUser method",
			args: args{
				from:   paymentGatewayService,
				to:     userService,
				method: getUserMethod,
			},
			expectedVal: true,
		},
		{
			name: "failed to verify from unknownService to user service with getUser method",
			args: args{
				from:   unknownService,
				to:     userService,
				method: getUserMethod,
			},
			expectedVal: false,
		},
		{
			name: "failed to verify from payment_gateway to user service with unknown method",
			args: args{
				from:   paymentGatewayService,
				to:     userService,
				method: unknownMethod,
			},
			expectedVal: false,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			val := verify(
				testcase.args.from,
				testcase.args.to,
				testcase.args.method,
			)
			assert.Equal(t, testcase.expectedVal, val)
		})
	}
}

func TestProcessRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		params apiMonitor.ProcessParams
	}

	testPaymentGatewayService := paymentGatewayService
	testUserService := userService
	testGetUserMethod := getUserMethod
	testPayload := "test-payload"
	testUnknownService := "test-service"
	testUnknownMethod := "test-method"

	testcases := []struct {
		name             string
		args             args
		expectedResponse middleware.Responder
	}{
		{
			name: "Successfully process request",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testPaymentGatewayService,
						To:      &testUserService,
						Method:  &testGetUserMethod,
						Payload: testPayload,
					},
				},
			},
			expectedResponse: apiMonitor.NewProcessOK(),
		},
		{
			name: "Failed to process request with unknown destination service",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testPaymentGatewayService,
						To:      &testUnknownService,
						Method:  &testGetUserMethod,
						Payload: testPayload,
					},
				},
			},
			expectedResponse: apiMonitor.NewProcessBadRequest().WithPayload(
				&models.ErrorResponse{
					Code:    int32(apiMonitor.ProcessBadRequestCode),
					Message: fmt.Sprintf("unknown destination: %s", testUnknownService),
				},
			),
		},
		{
			name: "Failed to process request with unknown method",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testPaymentGatewayService,
						To:      &testUserService,
						Method:  &testUnknownMethod,
						Payload: testPayload,
					},
				},
			},
			expectedResponse: apiMonitor.NewProcessBadRequest().WithPayload(
				&models.ErrorResponse{
					Code:    int32(apiMonitor.ProcessBadRequestCode),
					Message: fmt.Sprintf("unknown method: %s", testUnknownMethod),
				},
			),
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualResponse := processRequest(testcase.args.params)

			assert.Equal(t, testcase.expectedResponse, actualResponse)
		})
	}
}
