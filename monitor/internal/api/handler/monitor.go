package handler

import (
	"context"
	"fmt"

	"github.com/ShmelJUJ/software-engineering/monitor/internal/generated/models"
	apiMonitor "github.com/ShmelJUJ/software-engineering/monitor/internal/generated/restapi/operations/monitor"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	gen "github.com/ShmelJUJ/software-engineering/user/gen"
	"github.com/go-openapi/runtime/middleware"
	"github.com/mitchellh/mapstructure"
)

const (
	paymentGatewayService = "payment_gateway"
	userService           = "user"
	transactionService    = "transaction"

	getUserMethod   = "getClientByID"
	getWalletMethod = "getWalletByID"
	loginMethod     = "login"
)

// MonitorHandler handles incoming requests for monitoring.
type MonitorHandler struct {
	log        logger.Logger
	userClient gen.Handler
}

// NewMonitorHandler creates a new instance of MonitorHandler.
func NewMonitorHandler(log logger.Logger, userClient gen.Handler) *MonitorHandler {
	return &MonitorHandler{
		log:        log,
		userClient: userClient,
	}
}

// ProcessHandler processes incoming requests and returns a middleware.Responder.
func (mh *MonitorHandler) ProcessHandler(params apiMonitor.ProcessParams) middleware.Responder {
	from := *params.Body.From
	to := *params.Body.To
	method := *params.Body.Method

	mh.log.Debug("Process handler", map[string]interface{}{
		"from":    from,
		"to":      to,
		"method":  method,
		"payload": params.Body.Payload,
	})

	if !verify(from, to, method) {
		return apiMonitor.NewProcessForbidden().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiMonitor.ProcessForbiddenCode),
				Message: fmt.Sprintf("cannot process request from %s to %s with method %s", from, to, method),
			})
	}

	return mh.processRequest(params)
}

func verify(from, to, method string) bool {
	if from == paymentGatewayService && to == userService && (method == getUserMethod || method == getWalletMethod) {
		return true
	}

	if from == transactionService && to == userService && method == loginMethod {
		return true
	}

	return false
}

func (mh *MonitorHandler) processRequest(params apiMonitor.ProcessParams) middleware.Responder {
	ctx := context.Background()
	to := *params.Body.To
	method := *params.Body.Method

	switch to {
	case userService:
		switch method {
		case getUserMethod:
			dto := &gen.GetClientByIdParams{}

			if err := mapstructure.Decode(params.Body.Payload, dto); err != nil {
				return apiMonitor.NewProcessBadRequest().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessBadRequestCode),
						Message: fmt.Sprintf("failed to decode payload to gen.GetClientByIdParams: %s", err.Error()),
					})
			}

			client, err := mh.userClient.GetClientById(ctx, *dto)
			if err != nil {
				return apiMonitor.NewProcessInternalServerError().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessInternalServerErrorCode),
						Message: fmt.Sprintf("failed to get client by id: %s", err.Error()),
					})
			}

			switch c := client.(type) {
			case *gen.User:
				return apiMonitor.NewProcessOK().
					WithPayload(c)
			default:
				return apiMonitor.NewProcessInternalServerError().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessInternalServerErrorCode),
						Message: "failed to cast method info type to *get.User",
					})
			}

		case getWalletMethod:
			dto := &gen.GetWalletByIdParams{}

			if err := mapstructure.Decode(params.Body.Payload, dto); err != nil {
				return apiMonitor.NewProcessBadRequest().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessBadRequestCode),
						Message: fmt.Sprintf("failed to decode payload to gen.GetWalletByIdParams: %s", err.Error()),
					})
			}

			wallet, err := mh.userClient.GetWalletById(ctx, *dto)
			if err != nil {
				return apiMonitor.NewProcessInternalServerError().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessInternalServerErrorCode),
						Message: fmt.Sprintf("failed to get wallet by id: %s", err.Error()),
					})
			}

			switch w := wallet.(type) {
			case *gen.Wallet:
				return apiMonitor.NewProcessOK().
					WithPayload(w)
			default:
				return apiMonitor.NewProcessInternalServerError().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessInternalServerErrorCode),
						Message: "failed to cast method info type to *get.Wallet",
					})
			}

		case loginMethod:
			dto := &gen.AuthRequest{}

			if err := mapstructure.Decode(params.Body.Payload, dto); err != nil {
				return apiMonitor.NewProcessBadRequest().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessBadRequestCode),
						Message: fmt.Sprintf("failed to decode payload to gen.AuthRequest: %s", err.Error()),
					})
			}

			getTokenRes, err := mh.userClient.GetAuthToken(ctx, gen.OptAuthRequest{
				Value: *dto,
				Set:   true,
			})
			if err != nil {
				return apiMonitor.NewProcessInternalServerError().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessInternalServerErrorCode),
						Message: fmt.Sprintf("failed to get auth token: %s", err.Error()),
					})
			}

			switch t := getTokenRes.(type) {
			case *gen.AuthResponse:
				return apiMonitor.NewProcessOK().
					WithPayload(t)
			case *gen.GetAuthTokenBadRequest, *gen.GetAuthTokenForbidden, *gen.GetAuthTokenInternalServerError:
				return apiMonitor.NewProcessInternalServerError().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessInternalServerErrorCode),
						Message: fmt.Sprintf("failed to get auth token: %s", t),
					})
			default:
				return apiMonitor.NewProcessInternalServerError().
					WithPayload(&models.ErrorResponse{
						Code:    int32(apiMonitor.ProcessInternalServerErrorCode),
						Message: "failed to cast method info type to *get.AuthResponse",
					})
			}

		default:
			return apiMonitor.NewProcessBadRequest().WithPayload(
				&models.ErrorResponse{
					Code:    int32(apiMonitor.ProcessBadRequestCode),
					Message: fmt.Sprintf("unknown method: %s", method),
				},
			)
		}
	default:
		return apiMonitor.NewProcessBadRequest().WithPayload(
			&models.ErrorResponse{
				Code:    int32(apiMonitor.ProcessBadRequestCode),
				Message: fmt.Sprintf("unknown destination: %s", to),
			},
		)
	}
}
