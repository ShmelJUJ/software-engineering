package publisher

import "fmt"

// TransactionPublisherError represents an error that occurred during the creation transaction publishing.
type TransactionPublisherError struct {
	msg string
	err error
}

// NewTransactionPublisherError creates and returns a new instance of TransactionPublisherError.
func NewTransactionPublisherError(msg string, err error) *TransactionPublisherError {
	return &TransactionPublisherError{
		msg: msg,
		err: err,
	}
}

func (e TransactionPublisherError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

// PublishSucceededTransactionError represents an error when attempting to publish a successful transaction.
type PublishProcessedTransactionError struct {
	msg string
	err error
}

// NewPublishProcessedTransactionError creates and returns a new instance of PublishSucceededTransactionError.
func NewPublishProcessedTransactionError(msg string, err error) *PublishProcessedTransactionError {
	return &PublishProcessedTransactionError{
		msg: msg,
		err: err,
	}
}

func (e PublishProcessedTransactionError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}
