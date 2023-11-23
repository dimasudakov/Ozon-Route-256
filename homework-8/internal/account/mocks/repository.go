// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package mock_account is a generated GoMock package.
package mock_account

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	model "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateBankAccount mocks base method.
func (m *MockRepository) CreateBankAccount(ctx context.Context, account *model.BankAccount) (*model.BankAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBankAccount", ctx, account)
	ret0, _ := ret[0].(*model.BankAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBankAccount indicates an expected call of CreateBankAccount.
func (mr *MockRepositoryMockRecorder) CreateBankAccount(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBankAccount", reflect.TypeOf((*MockRepository)(nil).CreateBankAccount), ctx, account)
}

// DeleteBankAccount mocks base method.
func (m *MockRepository) DeleteBankAccount(ctx context.Context, id uuid.UUID) (*model.BankAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBankAccount", ctx, id)
	ret0, _ := ret[0].(*model.BankAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteBankAccount indicates an expected call of DeleteBankAccount.
func (mr *MockRepositoryMockRecorder) DeleteBankAccount(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBankAccount", reflect.TypeOf((*MockRepository)(nil).DeleteBankAccount), ctx, id)
}

// GetBankAccountByID mocks base method.
func (m *MockRepository) GetBankAccountByID(ctx context.Context, id uuid.UUID) (*model.BankAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBankAccountByID", ctx, id)
	ret0, _ := ret[0].(*model.BankAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBankAccountByID indicates an expected call of GetBankAccountByID.
func (mr *MockRepositoryMockRecorder) GetBankAccountByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBankAccountByID", reflect.TypeOf((*MockRepository)(nil).GetBankAccountByID), ctx, id)
}

// UpdateBankAccount mocks base method.
func (m *MockRepository) UpdateBankAccount(ctx context.Context, id uuid.UUID, account *model.BankAccount) (*model.BankAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBankAccount", ctx, id, account)
	ret0, _ := ret[0].(*model.BankAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBankAccount indicates an expected call of UpdateBankAccount.
func (mr *MockRepositoryMockRecorder) UpdateBankAccount(ctx, id, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBankAccount", reflect.TypeOf((*MockRepository)(nil).UpdateBankAccount), ctx, id, account)
}