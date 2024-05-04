package publisher

import "fmt"

// WorkerError represents an error encountered during creation payment worker.
type WorkerError struct {
	msg string
	err error
}

// NewWorkerError creates a new instance of WorkerError with the given message and error.
func NewWorkerError(msg string, err error) *WorkerError {
	return &WorkerError{
		msg: msg,
		err: err,
	}
}

func (e WorkerError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// StartError represents an error encountered during start payment worker.
type StartError struct {
	msg string
	err error
}

// NewStartError creates a new instance of StartError with the given message and error.
func NewStartError(msg string, err error) *StartError {
	return &StartError{
		msg: msg,
		err: err,
	}
}

func (e StartError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// ProcessPaymentError represents an error type associated with payment processing.
type ProccessPaymentError struct {
	msg string
	err error
}

// NewProcessPaymentError creates a new instance of ProcessPaymentError with the given message and error.
func NewProccessPaymentError(msg string, err error) *ProccessPaymentError {
	return &ProccessPaymentError{
		msg: msg,
		err: err,
	}
}

func (e ProccessPaymentError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}
