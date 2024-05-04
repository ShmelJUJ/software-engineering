// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ShmelJUJ/software-engineering/monitor/internal/broker/publisher (interfaces: MonitorPublisher)
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/monitor_publisher_mocks.go github.com/ShmelJUJ/software-engineering/monitor/internal/broker/publisher MonitorPublisher
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMonitorPublisher is a mock of MonitorPublisher interface.
type MockMonitorPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockMonitorPublisherMockRecorder
}

// MockMonitorPublisherMockRecorder is the mock recorder for MockMonitorPublisher.
type MockMonitorPublisherMockRecorder struct {
	mock *MockMonitorPublisher
}

// NewMockMonitorPublisher creates a new mock instance.
func NewMockMonitorPublisher(ctrl *gomock.Controller) *MockMonitorPublisher {
	mock := &MockMonitorPublisher{ctrl: ctrl}
	mock.recorder = &MockMonitorPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMonitorPublisher) EXPECT() *MockMonitorPublisherMockRecorder {
	return m.recorder
}

// PublishProcess mocks base method.
func (m *MockMonitorPublisher) PublishProcess(arg0 string, arg1 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishProcess", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishProcess indicates an expected call of PublishProcess.
func (mr *MockMonitorPublisherMockRecorder) PublishProcess(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishProcess", reflect.TypeOf((*MockMonitorPublisher)(nil).PublishProcess), arg0, arg1)
}