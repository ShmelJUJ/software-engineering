package handler

import (
	"github.com/ShmelJUJ/software-engineering/pkg/logger"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/generated/models"
	apiTransaction "github.com/ShmelJUJ/software-engineering/transaction/internal/generated/restapi/operations/transaction"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	"github.com/ShmelJUJ/software-engineering/transaction/internal/usecase"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
	log                logger.Logger
}

// NewTransactionHandler creates a new instance of transactionHandler.
func NewTransactionHandler(
	transactionUsecase usecase.TransactionUsecase,
	log logger.Logger,
) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
		log:                log,
	}
}

// AcceptTransactionHandler handles the request to accept a transaction.
func (th *TransactionHandler) AcceptTransactionHandler(params apiTransaction.AcceptTransactionParams) middleware.Responder {
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
func (th *TransactionHandler) CancelTransactionHandler(params apiTransaction.CancelTransactionParams) middleware.Responder {
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
func (th *TransactionHandler) CreateTransactionHandler(params apiTransaction.CreateTransactionParams) middleware.Responder {
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
func (th *TransactionHandler) EditTransactionHandler(params apiTransaction.EditTransactionParams) middleware.Responder {
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
func (th *TransactionHandler) RetrieveTransactionHandler(params apiTransaction.RetrieveTransactionParams) middleware.Responder {
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
func (th *TransactionHandler) RetrieveTransactionStatusHandler(params apiTransaction.RetrieveTransactionStatusParams) middleware.Responder {
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
