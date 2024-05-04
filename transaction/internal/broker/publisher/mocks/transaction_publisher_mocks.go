// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher (interfaces: TransactionPublisher)
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/transaction_publisher_mocks.go github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher TransactionPublisher
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	dto "github.com/ShmelJUJ/software-engineering/transaction/internal/broker/publisher/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockTransactionPublisher is a mock of TransactionPublisher interface.
type MockTransactionPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionPublisherMockRecorder
}

// MockTransactionPublisherMockRecorder is the mock recorder for MockTransactionPublisher.
type MockTransactionPublisherMockRecorder struct {
	mock *MockTransactionPublisher
}

// NewMockTransactionPublisher creates a new mock instance.
func NewMockTransactionPublisher(ctrl *gomock.Controller) *MockTransactionPublisher {
	mock := &MockTransactionPublisher{ctrl: ctrl}
	mock.recorder = &MockTransactionPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionPublisher) EXPECT() *MockTransactionPublisherMockRecorder {
	return m.recorder
}

// PublishProcessedTransaction mocks base method.
func (m *MockTransactionPublisher) PublishProcessedTransaction(arg0 *dto.ProcessedTransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishProcessedTransaction", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishProcessedTransaction indicates an expected call of PublishProcessedTransaction.
func (mr *MockTransactionPublisherMockRecorder) PublishProcessedTransaction(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishProcessedTransaction", reflect.TypeOf((*MockTransactionPublisher)(nil).PublishProcessedTransaction), arg0)
}
