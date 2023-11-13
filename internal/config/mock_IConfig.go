// Code generated by mockery v2.36.0. DO NOT EDIT.

package config

import mock "github.com/stretchr/testify/mock"

// MockIConfig is an autogenerated mock type for the IConfig type
type MockIConfig struct {
	mock.Mock
}

// Get provides a mock function with given fields: key
func (_m *MockIConfig) Get(key string) interface{} {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Set provides a mock function with given fields: key, value
func (_m *MockIConfig) Set(key string, value interface{}) {
	_m.Called(key, value)
}

// NewMockIConfig creates a new instance of MockIConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIConfig(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIConfig {
	mock := &MockIConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
