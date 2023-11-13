// Code generated by mockery v2.36.0. DO NOT EDIT.

package ddd

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockIEventDispatcher is an autogenerated mock type for the IEventDispatcher type
type MockIEventDispatcher[T IEvent] struct {
	mock.Mock
}

// Publish provides a mock function with given fields: ctx, events
func (_m *MockIEventDispatcher[T]) Publish(ctx context.Context, events ...T) error {
	_va := make([]interface{}, len(events))
	for _i := range events {
		_va[_i] = events[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...T) error); ok {
		r0 = rf(ctx, events...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: handler, events
func (_m *MockIEventDispatcher[T]) Subscribe(handler IEventHandler[T], events ...string) {
	_va := make([]interface{}, len(events))
	for _i := range events {
		_va[_i] = events[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handler)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// NewMockIEventDispatcher creates a new instance of MockIEventDispatcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIEventDispatcher[T IEvent](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIEventDispatcher[T] {
	mock := &MockIEventDispatcher[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
