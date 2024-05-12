package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	monitor_client "github.com/ShmelJUJ/software-engineering/pkg/monitor_client/client/monitor"
	monitor_models "github.com/ShmelJUJ/software-engineering/pkg/monitor_client/models"
	"github.com/ShmelJUJ/software-engineering/pkg/redis"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/generated/models"
	apiTransaction "github.com/ShmelJUJ/software-engineering/transaction/internal/generated/restapi/operations/transaction"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/usecase"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

const (
	transactionService = "transaction"
	userService        = "user"

	loginMethod = "login"

	transactionAuthKey    = "transaction-auth"
	defaultAuthExpiration = time.Hour
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
	monitorClient      monitor_client.ClientService
	r                  *redis.Redis
	log                logger.Logger
}

// NewTransactionHandler creates a new instance of transactionHandler.
func NewTransactionHandler(
	transactionUsecase usecase.TransactionUsecase,
	log logger.Logger,
	monitorClient monitor_client.ClientService,
	r *redis.Redis,
) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
		log:                log,
		monitorClient:      monitorClient,
		r:                  r,
	}
}

// AcceptTransactionHandler handles the request to accept a transaction.
func (th *TransactionHandler) AcceptTransactionHandler(params apiTransaction.AcceptTransactionParams, _ interface{}) middleware.Responder {
	th.log.Debug("Accept transaction handler", map[string]interface{}{
		"transaction_id": params.ID.String(),
		"body":           params.Body,
	})

	sender := model.FromAcceptTransactionUserDTO(params.Body.Sender)

	if err := th.transactionUsecase.AcceptTransaction(
		params.HTTPRequest.Context(),
		params.ID.String(),
		sender,
	); err != nil {
		return apiTransaction.NewAcceptTransactionInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.AcceptTransactionInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	return apiTransaction.NewAcceptTransactionOK()
}

// CancelTransactionHandler handles the request to cancel a transaction.
func (th *TransactionHandler) CancelTransactionHandler(params apiTransaction.CancelTransactionParams, _ interface{}) middleware.Responder {
	th.log.Debug("Cancel transaction handler", map[string]interface{}{
		"transaction_id": params.ID.String(),
		"body":           params.Body,
	})

	if err := th.transactionUsecase.CancelTransaction(
		params.HTTPRequest.Context(),
		params.ID.String(),
		params.Body.Reason,
	); err != nil {
		return apiTransaction.NewCancelTransactionInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.CancelTransactionInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	return apiTransaction.NewCancelTransactionOK()
}

// CreateTransactionHandler handles the request to create a transaction.
func (th *TransactionHandler) CreateTransactionHandler(params apiTransaction.CreateTransactionParams, _ interface{}) middleware.Responder {
	th.log.Debug("Create transaction handler", map[string]interface{}{
		"body": params.Body,
	})

	transaction := model.FromCreateTransactionDTO(params.Body)

	if err := th.transactionUsecase.CreateTransaction(
		params.HTTPRequest.Context(),
		transaction,
	); err != nil {
		return apiTransaction.NewCreateTransactionInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.CreateTransactionInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	return apiTransaction.NewCreateTransactionOK().
		WithPayload(&models.CreateTransactionResponse{
			TransactionID: (*strfmt.UUID)(&transaction.ID),
		})
}

// EditTransactionHandler handles the request to edit a transaction.
func (th *TransactionHandler) EditTransactionHandler(params apiTransaction.EditTransactionParams, _ interface{}) middleware.Responder {
	th.log.Debug("Edit transaction handler", map[string]interface{}{
		"transaction_id": params.ID.String(),
		"body":           params.Body,
	})

	transaction := model.FromEditTransactionDTO(params.ID.String(), params.Body)

	if err := th.transactionUsecase.UpdateTransaction(
		params.HTTPRequest.Context(),
		transaction,
	); err != nil {
		return apiTransaction.NewEditTransactionInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.EditTransactionInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	return apiTransaction.NewEditTransactionOK()
}

// RetrieveTransactionHandler handles the request to retrieve a transaction.
func (th *TransactionHandler) RetrieveTransactionHandler(params apiTransaction.RetrieveTransactionParams, _ interface{}) middleware.Responder {
	th.log.Debug("Retrieve transaction handler", map[string]interface{}{
		"transaction_id": params.ID.String(),
	})

	transaction, err := th.transactionUsecase.GetTransaction(
		params.HTTPRequest.Context(),
		params.ID.String(),
	)
	if err != nil {
		return apiTransaction.NewRetrieveTransactionInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.RetrieveTransactionInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	return apiTransaction.NewRetrieveTransactionOK().
		WithPayload(transaction.ToGetTransactionDTO())
}

// RetrieveTransactionStatusHandler handles the request to retrieve the status of a transaction.
func (th *TransactionHandler) RetrieveTransactionStatusHandler(params apiTransaction.RetrieveTransactionStatusParams, _ interface{}) middleware.Responder {
	th.log.Debug("Retrieve transaction status handler", map[string]interface{}{
		"transaction_id": params.ID.String(),
	})

	transactionStatus, err := th.transactionUsecase.GetTransactionStatus(
		params.HTTPRequest.Context(),
		params.ID.String(),
	)
	if err != nil {
		return apiTransaction.NewRetrieveTransactionInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.RetrieveTransactionInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	transactionStatusStr := transactionStatus.String()

	return apiTransaction.NewRetrieveTransactionStatusOK().
		WithPayload(&models.GetTransactionStatusResponse{
			TransactionStatus: &transactionStatusStr,
		})
}

func (th *TransactionHandler) LoginHandler(params apiTransaction.LoginParams) middleware.Responder {
	from := transactionService
	to := userService
	method := loginMethod

	th.log.Debug("Login handler", map[string]interface{}{
		"from":   from,
		"to":     to,
		"method": method,
		"body":   params.Body,
	})

	resp, err := th.monitorClient.Process(&monitor_client.ProcessParams{
		Body: &monitor_models.ProcessRequest{
			From:    &from,
			To:      &to,
			Method:  &method,
			Payload: params.Body,
		},
	})
	if err != nil {
		return apiTransaction.NewLoginInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.LoginInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	payload := resp.Payload.(map[string]interface{}) //nolint:errcheck // blya budu tut chto nado

	token := payload["auth_token"].(string) //nolint:errcheck // blya budu tut chto nado

	if err := th.r.Client.SAdd(params.HTTPRequest.Context(), transactionAuthKey, token).Err(); err != nil {
		return apiTransaction.NewLoginInternalServerError().
			WithPayload(&models.ErrorResponse{
				Code:    int32(apiTransaction.LoginInternalServerErrorCode),
				Message: err.Error(),
			})
	}

	return apiTransaction.NewLoginOK().
		WithPayload(&models.LoginResponse{
			AuthToken: &token,
		})
}

func (th *TransactionHandler) VerifyAuthToken(token string) (interface{}, error) {
	ctx := context.Background()

	th.log.Debug("Verify auth token", map[string]interface{}{
		"token": token,
	})

	res, err := th.r.Client.SIsMember(ctx, transactionAuthKey, token).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check that token contains in transactionAuthKey: %w", err)
	}

	if !res {
		return nil, nil //nolint:nilnil // to get 401 error
	}

	return true, nil
}
