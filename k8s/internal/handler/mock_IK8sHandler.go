// Code generated by mockery v2.36.0. DO NOT EDIT.

package handler

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockIK8sHandler is an autogenerated mock type for the IK8sHandler type
type MockIK8sHandler struct {
	mock.Mock
}

// GetLivenessProbe provides a mock function with given fields: w, r
func (_m *MockIK8sHandler) GetLivenessProbe(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetReadinessProbe provides a mock function with given fields: w, r
func (_m *MockIK8sHandler) GetReadinessProbe(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// NewMockIK8sHandler creates a new instance of MockIK8sHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIK8sHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIK8sHandler {
	mock := &MockIK8sHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
