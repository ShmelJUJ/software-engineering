package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/ShmelJUJ/software-engineering/monitor/internal/generated/models"
	apiMonitor "github.com/ShmelJUJ/software-engineering/monitor/internal/generated/restapi/operations/monitor"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	mock_logger "github.com/ShmelJUJ/software-engineering/pkg/logger/mocks"
	gen "github.com/ShmelJUJ/software-engineering/user/gen"
	mock_user_client "github.com/ShmelJUJ/software-engineering/user/mocks"
	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func monitorHandlerHelper(t *testing.T) (*mock_logger.MockLogger, *mock_user_client.MockHandler) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	log := mock_logger.NewMockLogger(mockCtrl)
	userClient := mock_user_client.NewMockHandler(mockCtrl)

	return log, userClient
}

func TestNewMonitorHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		log        logger.Logger
		userClient gen.Handler
	}

	log, userClient := monitorHandlerHelper(t)

	testcases := []struct {
		name                   string
		args                   args
		expectedMonitorHandler *MonitorHandler
	}{
		{
			name: "Successfully create new monitor handler",
			args: args{
				log:        log,
				userClient: userClient,
			},
			expectedMonitorHandler: &MonitorHandler{
				log:        log,
				userClient: userClient,
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			actualMonitorHandler := NewMonitorHandler(
				testcase.args.log,
				testcase.args.userClient,
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
	testPayload := gen.GetClientByIdParams{}
	testUnknownService := "test-service"
	ctx := context.Background()
	res := &gen.User{}

	testcases := []struct {
		name             string
		args             args
		mock             func(*mock_logger.MockLogger, *mock_user_client.MockHandler)
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
			mock: func(ml *mock_logger.MockLogger, mh *mock_user_client.MockHandler) {
				ml.EXPECT().Debug("Process handler", map[string]interface{}{
					"from":    testPaymentGatewayService,
					"to":      testUserService,
					"method":  testGetUserMethod,
					"payload": testPayload,
				})
				mh.EXPECT().GetClientById(ctx, testPayload).Return(res, nil).Times(1)
			},
			expectedResponse: apiMonitor.NewProcessOK().
				WithPayload(res),
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
			mock: func(ml *mock_logger.MockLogger, _ *mock_user_client.MockHandler) {
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

			log, userClient := monitorHandlerHelper(t)

			testcase.mock(log, userClient)

			actualMonitorHandler := NewMonitorHandler(
				log,
				userClient,
			)

			actualResponse := actualMonitorHandler.ProcessHandler(testcase.args.params)

			if bad, ok := actualResponse.(*apiMonitor.ProcessBadRequest); ok {
				fmt.Println(bad.Payload)
			}

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
			name: "Successful verify from payment gateway to user service with getClient method",
			args: args{
				from:   paymentGatewayService,
				to:     userService,
				method: getUserMethod,
			},
			expectedVal: true,
		},
		{
			name: "Successful verify from payment gateway to user service with getWallet method",
			args: args{
				from:   paymentGatewayService,
				to:     userService,
				method: getWalletMethod,
			},
			expectedVal: true,
		},
		{
			name: "failed to verify from unknownService to user service with getClient method",
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

	var (
		ctx = context.Background()

		testPaymentGatewayService = paymentGatewayService
		testUserService           = userService

		testUnknownService = "test-service"
		testUnknownMethod  = "test-method"

		testGetUserMethod = getUserMethod
		testClientPayload = gen.GetClientByIdParams{}
		testClientRes     = &gen.User{}

		testGetWalletMethod = getWalletMethod
		testWalletPayload   = gen.GetWalletByIdParams{}
		testWalletRes       = &gen.Wallet{}
	)

	testcases := []struct {
		name             string
		args             args
		mock             func(*mock_user_client.MockHandler)
		expectedResponse middleware.Responder
	}{
		{
			name: "Successfully process request getClientByID",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testPaymentGatewayService,
						To:      &testUserService,
						Method:  &testGetUserMethod,
						Payload: testClientPayload,
					},
				},
			},
			mock: func(mh *mock_user_client.MockHandler) {
				mh.EXPECT().GetClientById(ctx, testClientPayload).Return(testClientRes, nil).Times(1)
			},
			expectedResponse: apiMonitor.NewProcessOK().
				WithPayload(testClientRes),
		},
		{
			name: "Successfully process request getWalletByID",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testPaymentGatewayService,
						To:      &testUserService,
						Method:  &testGetWalletMethod,
						Payload: testWalletPayload,
					},
				},
			},
			mock: func(mh *mock_user_client.MockHandler) {
				mh.EXPECT().GetWalletById(ctx, testWalletPayload).Return(testWalletRes, nil).Times(1)
			},
			expectedResponse: apiMonitor.NewProcessOK().
				WithPayload(testWalletRes),
		},
		{
			name: "Failed to process request with unknown destination service",
			args: args{
				params: apiMonitor.ProcessParams{
					Body: &models.ProcessRequest{
						From:    &testPaymentGatewayService,
						To:      &testUnknownService,
						Method:  &testGetUserMethod,
						Payload: testClientPayload,
					},
				},
			},
			mock: func(_ *mock_user_client.MockHandler) {},
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
						Payload: testClientPayload,
					},
				},
			},
			mock: func(_ *mock_user_client.MockHandler) {},
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

			log, userClient := monitorHandlerHelper(t)

			testcase.mock(userClient)

			monitorHandler := NewMonitorHandler(
				log,
				userClient,
			)

			actualResponse := monitorHandler.processRequest(testcase.args.params)

			assert.Equal(t, testcase.expectedResponse, actualResponse)
		})
	}
}
