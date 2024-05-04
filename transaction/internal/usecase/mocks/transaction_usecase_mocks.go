// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ShmelJUJ/software-engineering/transaction/internal/usecase (interfaces: TransactionUsecase)
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/transaction_usecase_mocks.go github.com/ShmelJUJ/software-engineering/transaction/internal/usecase TransactionUsecase
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/ShmelJUJ/software-engineering/transaction/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockTransactionUsecase is a mock of TransactionUsecase interface.
type MockTransactionUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionUsecaseMockRecorder
}

// MockTransactionUsecaseMockRecorder is the mock recorder for MockTransactionUsecase.
type MockTransactionUsecaseMockRecorder struct {
	mock *MockTransactionUsecase
}

// NewMockTransactionUsecase creates a new mock instance.
func NewMockTransactionUsecase(ctrl *gomock.Controller) *MockTransactionUsecase {
	mock := &MockTransactionUsecase{ctrl: ctrl}
	mock.recorder = &MockTransactionUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionUsecase) EXPECT() *MockTransactionUsecaseMockRecorder {
	return m.recorder
}

// AcceptTransaction mocks base method.
func (m *MockTransactionUsecase) AcceptTransaction(arg0 context.Context, arg1 string, arg2 *model.TransactionUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptTransaction", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcceptTransaction indicates an expected call of AcceptTransaction.
func (mr *MockTransactionUsecaseMockRecorder) AcceptTransaction(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptTransaction", reflect.TypeOf((*MockTransactionUsecase)(nil).AcceptTransaction), arg0, arg1, arg2)
}

// CancelTransaction mocks base method.
func (m *MockTransactionUsecase) CancelTransaction(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelTransaction", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelTransaction indicates an expected call of CancelTransaction.
func (mr *MockTransactionUsecaseMockRecorder) CancelTransaction(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelTransaction", reflect.TypeOf((*MockTransactionUsecase)(nil).CancelTransaction), arg0, arg1, arg2)
}

// ChangeTransactionStatus mocks base method.
func (m *MockTransactionUsecase) ChangeTransactionStatus(arg0 context.Context, arg1 string, arg2 model.TransactionStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeTransactionStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeTransactionStatus indicates an expected call of ChangeTransactionStatus.
func (mr *MockTransactionUsecaseMockRecorder) ChangeTransactionStatus(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeTransactionStatus", reflect.TypeOf((*MockTransactionUsecase)(nil).ChangeTransactionStatus), arg0, arg1, arg2)
}

// CreateTransaction mocks base method.
func (m *MockTransactionUsecase) CreateTransaction(arg0 context.Context, arg1 *model.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockTransactionUsecaseMockRecorder) CreateTransaction(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionUsecase)(nil).CreateTransaction), arg0, arg1)
}

// GetTransaction mocks base method.
func (m *MockTransactionUsecase) GetTransaction(arg0 context.Context, arg1 string) (*model.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", arg0, arg1)
	ret0, _ := ret[0].(*model.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockTransactionUsecaseMockRecorder) GetTransaction(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockTransactionUsecase)(nil).GetTransaction), arg0, arg1)
}

// GetTransactionStatus mocks base method.
func (m *MockTransactionUsecase) GetTransactionStatus(arg0 context.Context, arg1 string) (model.TransactionStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionStatus", arg0, arg1)
	ret0, _ := ret[0].(model.TransactionStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionStatus indicates an expected call of GetTransactionStatus.
func (mr *MockTransactionUsecaseMockRecorder) GetTransactionStatus(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionStatus", reflect.TypeOf((*MockTransactionUsecase)(nil).GetTransactionStatus), arg0, arg1)
}

// UpdateTransaction mocks base method.
func (m *MockTransactionUsecase) UpdateTransaction(arg0 context.Context, arg1 *model.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockTransactionUsecaseMockRecorder) UpdateTransaction(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockTransactionUsecase)(nil).UpdateTransaction), arg0, arg1)
}