package subscriber

import "fmt"

type MonitorSubscriberError struct {
	msg string
	err error
}

func NewMonitorSubscriberError(msg string, err error) *MonitorSubscriberError {
	return &MonitorSubscriberError{
		msg: msg,
		err: err,
	}
}

func (e MonitorSubscriberError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

type HandleProcessError struct {
	msg string
	err error
}

func NewHandleProcessError(msg string, err error) *HandleProcessError {
	return &HandleProcessError{
		msg: msg,
		err: err,
	}
}

func (e HandleProcessError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}
