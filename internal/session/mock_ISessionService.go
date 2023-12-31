// Code generated by mockery v2.36.0. DO NOT EDIT.

package session

import mock "github.com/stretchr/testify/mock"

// MockISessionService is an autogenerated mock type for the ISessionService type
type MockISessionService struct {
	mock.Mock
}

// GetNodeProvisionTimeStamp provides a mock function with given fields: NodeName
func (_m *MockISessionService) GetNodeProvisionTimeStamp(NodeName string) (int64, error) {
	ret := _m.Called(NodeName)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(NodeName)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(NodeName)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(NodeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPodScheduleTimeStamp provides a mock function with given fields: sessionId
func (_m *MockISessionService) GetPodScheduleTimeStamp(sessionId string) (int64, error) {
	ret := _m.Called(sessionId)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(sessionId)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(sessionId)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sessionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetNodeProvisionTimeStamp provides a mock function with given fields: _a0
func (_m *MockISessionService) SetNodeProvisionTimeStamp(_a0 *SetNodeProvisionTimeStampActionPayload) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*SetNodeProvisionTimeStampActionPayload) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPodScheduleTimeStamp provides a mock function with given fields: _a0
func (_m *MockISessionService) SetPodScheduleTimeStamp(_a0 *UpdateSessionTimeStampLikeFieldActionPayload) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*UpdateSessionTimeStampLikeFieldActionPayload) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetSessionDeletable provides a mock function with given fields: _a0
func (_m *MockISessionService) SetSessionDeletable(_a0 *SetSessionDeletableActionPayload) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*SetSessionDeletableActionPayload) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetSessionReady provides a mock function with given fields: _a0
func (_m *MockISessionService) SetSessionReady(_a0 *SetSessionReadyActionPayload) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*SetSessionReadyActionPayload) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockISessionService creates a new instance of MockISessionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockISessionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockISessionService {
	mock := &MockISessionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
