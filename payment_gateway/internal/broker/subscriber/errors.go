package subscriber

import "fmt"

// TransactionSubscriberError represents an error encountered during creation transaction subscriber.
type TransactionSubscriberError struct {
	msg string
	err error
}

// NewTransactionSubscriberError creates a new TransactionSubscriberError instance with the given message and error.
func NewTransactionSubscriberError(msg string, err error) *TransactionSubscriberError {
	return &TransactionSubscriberError{
		msg: msg,
		err: err,
	}
}

func (e TransactionSubscriberError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}
