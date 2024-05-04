// Code generated by go-swagger; DO NOT EDIT.

package transaction

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ShmelJUJ/software-engineering/transaction/internal/generated/models"
)

// AcceptTransactionOKCode is the HTTP code returned for type AcceptTransactionOK
const AcceptTransactionOKCode int = 200

/*
AcceptTransactionOK Transaction successfully accepted.

swagger:response acceptTransactionOK
*/
type AcceptTransactionOK struct {
}

// NewAcceptTransactionOK creates AcceptTransactionOK with default headers values
func NewAcceptTransactionOK() *AcceptTransactionOK {

	return &AcceptTransactionOK{}
}

// WriteResponse to the client
func (o *AcceptTransactionOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// AcceptTransactionForbiddenCode is the HTTP code returned for type AcceptTransactionForbidden
const AcceptTransactionForbiddenCode int = 403

/*
AcceptTransactionForbidden Forbidden error.

swagger:response acceptTransactionForbidden
*/
type AcceptTransactionForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewAcceptTransactionForbidden creates AcceptTransactionForbidden with default headers values
func NewAcceptTransactionForbidden() *AcceptTransactionForbidden {

	return &AcceptTransactionForbidden{}
}

// WithPayload adds the payload to the accept transaction forbidden response
func (o *AcceptTransactionForbidden) WithPayload(payload *models.ErrorResponse) *AcceptTransactionForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the accept transaction forbidden response
func (o *AcceptTransactionForbidden) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcceptTransactionForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AcceptTransactionNotFoundCode is the HTTP code returned for type AcceptTransactionNotFound
const AcceptTransactionNotFoundCode int = 404

/*
AcceptTransactionNotFound Not found error.

swagger:response acceptTransactionNotFound
*/
type AcceptTransactionNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewAcceptTransactionNotFound creates AcceptTransactionNotFound with default headers values
func NewAcceptTransactionNotFound() *AcceptTransactionNotFound {

	return &AcceptTransactionNotFound{}
}

// WithPayload adds the payload to the accept transaction not found response
func (o *AcceptTransactionNotFound) WithPayload(payload *models.ErrorResponse) *AcceptTransactionNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the accept transaction not found response
func (o *AcceptTransactionNotFound) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcceptTransactionNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AcceptTransactionInternalServerErrorCode is the HTTP code returned for type AcceptTransactionInternalServerError
const AcceptTransactionInternalServerErrorCode int = 500

/*
AcceptTransactionInternalServerError Internal server error.

swagger:response acceptTransactionInternalServerError
*/
type AcceptTransactionInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewAcceptTransactionInternalServerError creates AcceptTransactionInternalServerError with default headers values
func NewAcceptTransactionInternalServerError() *AcceptTransactionInternalServerError {

	return &AcceptTransactionInternalServerError{}
}

// WithPayload adds the payload to the accept transaction internal server error response
func (o *AcceptTransactionInternalServerError) WithPayload(payload *models.ErrorResponse) *AcceptTransactionInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the accept transaction internal server error response
func (o *AcceptTransactionInternalServerError) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AcceptTransactionInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
