package repository

import "fmt"

// GetTransactionError represents an error encountered while getting a transaction.
type GetTransactionError struct {
	msg string
	err error
}

// NewGetTransactionError creates a new GetTransactionError instance with the provided message and error.
func NewGetTransactionError(msg string, err error) *GetTransactionError {
	return &GetTransactionError{
		msg: msg,
		err: err,
	}
}

func (e GetTransactionError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// GetTransactionUserError represents a user-related error encountered while getting a transaction.
type GetTransactionUserError struct {
	msg string
	err error
}

// NewGetTransactionUserError creates a new GetTransactionUserError instance with the provided message and error.
func NewGetTransactionUserError(msg string, err error) *GetTransactionUserError {
	return &GetTransactionUserError{
		msg: msg,
		err: err,
	}
}

func (e GetTransactionUserError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// CreateTransactionError represents an error encountered while creating a transaction.
type CreateTransactionError struct {
	msg string
	err error
}

// NewCreateTransactionError creates a new CreateTransactionError instance with the provided message and error.
func NewCreateTransactionError(msg string, err error) *CreateTransactionError {
	return &CreateTransactionError{
		msg: msg,
		err: err,
	}
}

func (e CreateTransactionError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// GetTransactionStatusError represents an error encountered while getting the status of a transaction.
type GetTransactionStatusError struct {
	msg string
	err error
}

// NewGetTransactionStatusError creates a new GetTransactionStatusError instance with the provided message and error.
func NewGetTransactionStatusError(msg string, err error) *GetTransactionStatusError {
	return &GetTransactionStatusError{
		msg: msg,
		err: err,
	}
}

func (e GetTransactionStatusError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// CancelTransactionError represents an error encountered while canceling a transaction.
type CancelTransactionError struct {
	msg string
	err error
}

// NewCancelTransactionError creates a new CancelTransactionError instance with the provided message and error.
func NewCancelTransactionError(msg string, err error) *CancelTransactionError {
	return &CancelTransactionError{
		msg: msg,
		err: err,
	}
}

func (e CancelTransactionError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// AcceptTransactionError represents an error encountered while accepting a transaction.
type AcceptTransactionError struct {
	msg string
	err error
}

// NewAcceptTransactionError creates a new AcceptTransactionError instance with the provided message and error.
func NewAcceptTransactionError(msg string, err error) *AcceptTransactionError {
	return &AcceptTransactionError{
		msg: msg,
		err: err,
	}
}

func (e AcceptTransactionError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// UpdateTransactionError represents an error encountered while updating a transaction.
type UpdateTransactionError struct {
	msg string
	err error
}

// NewUpdateTransactionError creates a new UpdateTransactionError instance with the provided message and error.
func NewUpdateTransactionError(msg string, err error) *UpdateTransactionError {
	return &UpdateTransactionError{
		msg: msg,
		err: err,
	}
}

func (e UpdateTransactionError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// CreateTransactionUserError represents a user-related error encountered while creating a transaction.
type CreateTransactionUserError struct {
	msg string
	err error
}

// NewCreateTransactionUserError creates a new CreateTransactionUserError instance with the provided message and error.
func NewCreateTransactionUserError(msg string, err error) *CreateTransactionUserError {
	return &CreateTransactionUserError{
		msg: msg,
		err: err,
	}
}

func (e CreateTransactionUserError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// ChangeTransactionStatusError represents an error encountered while changing the status of a transaction.
type ChangeTransactionStatusError struct {
	msg string
	err error
}

// NewChangeTransactionStatusError creates a new ChangeTransactionStatusError instance with the provided message and error.
func NewChangeTransactionStatusError(msg string, err error) *ChangeTransactionStatusError {
	return &ChangeTransactionStatusError{
		msg: msg,
		err: err,
	}
}

func (e ChangeTransactionStatusError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}
