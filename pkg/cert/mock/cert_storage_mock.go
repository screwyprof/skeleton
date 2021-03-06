// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/cert/usecase/issuecert/issue_certificate.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	issuecert "github.com/screwyprof/skeleton/pkg/cert/usecase/issuecert"
)

// MockCertStorage is a mock of CertStorage interface
type MockCertStorage struct {
	ctrl     *gomock.Controller
	recorder *MockCertStorageMockRecorder
}

// MockCertStorageMockRecorder is the mock recorder for MockCertStorage
type MockCertStorageMockRecorder struct {
	mock *MockCertStorage
}

// NewMockCertStorage creates a new mock instance
func NewMockCertStorage(ctrl *gomock.Controller) *MockCertStorage {
	mock := &MockCertStorage{ctrl: ctrl}
	mock.recorder = &MockCertStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCertStorage) EXPECT() *MockCertStorageMockRecorder {
	return m.recorder
}

// Store mocks base method
func (m *MockCertStorage) Store(arg0 context.Context, arg1 *issuecert.Certificate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockCertStorageMockRecorder) Store(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockCertStorage)(nil).Store), arg0, arg1)
}
