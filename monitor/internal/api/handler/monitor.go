package handler

import (
	"fmt"

	"github.com/ShmelJUJ/software-engineering/monitor/internal/generated/models"
	apiMonitor "github.com/ShmelJUJ/software-engineering/monitor/internal/generated/restapi/operations/monitor"
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/go-openapi/runtime/middleware"
)

const (
	paymentGatewayService = "payment_gateway"
	userService           = "user"

	getUserMethod = "getUserData"
)

// MonitorHandler handles incoming requests for monitoring.
type MonitorHandler struct {
	log logger.Logger
}

// NewMonitorHandler creates a new instance of MonitorHandler.
func NewMonitorHandler(log logger.Logger) *MonitorHandler {
	return &MonitorHandler{
		log: log,
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

	return processRequest(params)
}

func verify(from, to, method string) bool {
	if from == paymentGatewayService && to == userService && method == getUserMethod {
		return true
	}

	return false
}

func processRequest(params apiMonitor.ProcessParams) middleware.Responder {
	to := *params.Body.To
	method := *params.Body.Method

	switch to {
	case userService:
		switch method {
		case getUserMethod:
			// TODO: add user client.
			return apiMonitor.NewProcessOK()
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
