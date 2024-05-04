package gateway

import "fmt"

// GatewayError represents a gateway error type with a specific message and underlying error.
type CreationGatewayError struct {
	msg string
	err error
}

// NewCreationGatewayError creates a new GatewayError instance with the given message and underlying error.
func NewCreationGatewayError(msg string, err error) *CreationGatewayError {
	return &CreationGatewayError{
		msg: msg,
		err: err,
	}
}

func (e *CreationGatewayError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// CreatePaymentError represents an error type specific to creating payments.
type CreatePaymentError struct {
	msg string
	err error
}

// NewCreatePaymentError creates a new CreatePaymentError instance with the given message and underlying error.
func NewCreatePaymentError(msg string, err error) *CreatePaymentError {
	return &CreatePaymentError{
		msg: msg,
		err: err,
	}
}

func (e *CreatePaymentError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// CheckStatusError represents an error type specific to checking status.
type CheckStatusError struct {
	msg string
	err error
}

// NewCheckStatusError creates a new CheckStatusError instance with the given message and underlying error.
func NewCheckStatusError(msg string, err error) *CheckStatusError {
	return &CheckStatusError{
		msg: msg,
		err: err,
	}
}

func (e *CheckStatusError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}
