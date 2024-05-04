package publisher

import "fmt"

type PublishProcessError struct {
	msg string
	err error
}

func NewPublishProcessError(msg string, err error) *PublishProcessError {
	return &PublishProcessError{
		msg: msg,
		err: err,
	}
}

func (e PublishProcessError) Error() string {
	var errMsg string
	if e.err != nil {
		errMsg = e.err.Error()
	}

	return fmt.Sprintf("%s: %s", e.msg, errMsg)
}
