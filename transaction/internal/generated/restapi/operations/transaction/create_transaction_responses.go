// Code generated by go-swagger; DO NOT EDIT.

package transaction

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ShmelJUJ/software-engineering/transaction/internal/generated/models"
)

// CreateTransactionOKCode is the HTTP code returned for type CreateTransactionOK
const CreateTransactionOKCode int = 200

/*
CreateTransactionOK Transaction successfully created.

swagger:response createTransactionOK
*/
type CreateTransactionOK struct {

	/*
	  In: Body
	*/
	Payload *models.CreateTransactionResponse `json:"body,omitempty"`
}

// NewCreateTransactionOK creates CreateTransactionOK with default headers values
func NewCreateTransactionOK() *CreateTransactionOK {

	return &CreateTransactionOK{}
}

// WithPayload adds the payload to the create transaction o k response
func (o *CreateTransactionOK) WithPayload(payload *models.CreateTransactionResponse) *CreateTransactionOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create transaction o k response
func (o *CreateTransactionOK) SetPayload(payload *models.CreateTransactionResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTransactionOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTransactionBadRequestCode is the HTTP code returned for type CreateTransactionBadRequest
const CreateTransactionBadRequestCode int = 400

/*
CreateTransactionBadRequest Validation error.

swagger:response createTransactionBadRequest
*/
type CreateTransactionBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewCreateTransactionBadRequest creates CreateTransactionBadRequest with default headers values
func NewCreateTransactionBadRequest() *CreateTransactionBadRequest {

	return &CreateTransactionBadRequest{}
}

// WithPayload adds the payload to the create transaction bad request response
func (o *CreateTransactionBadRequest) WithPayload(payload *models.ErrorResponse) *CreateTransactionBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create transaction bad request response
func (o *CreateTransactionBadRequest) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTransactionBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTransactionForbiddenCode is the HTTP code returned for type CreateTransactionForbidden
const CreateTransactionForbiddenCode int = 403

/*
CreateTransactionForbidden Forbidden error.

swagger:response createTransactionForbidden
*/
type CreateTransactionForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewCreateTransactionForbidden creates CreateTransactionForbidden with default headers values
func NewCreateTransactionForbidden() *CreateTransactionForbidden {

	return &CreateTransactionForbidden{}
}

// WithPayload adds the payload to the create transaction forbidden response
func (o *CreateTransactionForbidden) WithPayload(payload *models.ErrorResponse) *CreateTransactionForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create transaction forbidden response
func (o *CreateTransactionForbidden) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTransactionForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTransactionInternalServerErrorCode is the HTTP code returned for type CreateTransactionInternalServerError
const CreateTransactionInternalServerErrorCode int = 500

/*
CreateTransactionInternalServerError Internal server error.

swagger:response createTransactionInternalServerError
*/
type CreateTransactionInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewCreateTransactionInternalServerError creates CreateTransactionInternalServerError with default headers values
func NewCreateTransactionInternalServerError() *CreateTransactionInternalServerError {

	return &CreateTransactionInternalServerError{}
}

// WithPayload adds the payload to the create transaction internal server error response
func (o *CreateTransactionInternalServerError) WithPayload(payload *models.ErrorResponse) *CreateTransactionInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create transaction internal server error response
func (o *CreateTransactionInternalServerError) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTransactionInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
