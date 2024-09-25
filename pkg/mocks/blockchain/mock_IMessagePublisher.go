// Code generated by mockery v2.44.1. DO NOT EDIT.

package blockchain

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockIMessagePublisher is an autogenerated mock type for the IMessagePublisher type
type MockIMessagePublisher struct {
	mock.Mock
}

type MockIMessagePublisher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIMessagePublisher) EXPECT() *MockIMessagePublisher_Expecter {
	return &MockIMessagePublisher_Expecter{mock: &_m.Mock}
}

// PublishGroupMessage provides a mock function with given fields: ctx, groupdId, message
func (_m *MockIMessagePublisher) PublishGroupMessage(ctx context.Context, groupdId [32]byte, message []byte) error {
	ret := _m.Called(ctx, groupdId, message)

	if len(ret) == 0 {
		panic("no return value specified for PublishGroupMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, [32]byte, []byte) error); ok {
		r0 = rf(ctx, groupdId, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIMessagePublisher_PublishGroupMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishGroupMessage'
type MockIMessagePublisher_PublishGroupMessage_Call struct {
	*mock.Call
}

// PublishGroupMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - groupdId [32]byte
//   - message []byte
func (_e *MockIMessagePublisher_Expecter) PublishGroupMessage(ctx interface{}, groupdId interface{}, message interface{}) *MockIMessagePublisher_PublishGroupMessage_Call {
	return &MockIMessagePublisher_PublishGroupMessage_Call{Call: _e.mock.On("PublishGroupMessage", ctx, groupdId, message)}
}

func (_c *MockIMessagePublisher_PublishGroupMessage_Call) Run(run func(ctx context.Context, groupdId [32]byte, message []byte)) *MockIMessagePublisher_PublishGroupMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([32]byte), args[2].([]byte))
	})
	return _c
}

func (_c *MockIMessagePublisher_PublishGroupMessage_Call) Return(_a0 error) *MockIMessagePublisher_PublishGroupMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIMessagePublisher_PublishGroupMessage_Call) RunAndReturn(run func(context.Context, [32]byte, []byte) error) *MockIMessagePublisher_PublishGroupMessage_Call {
	_c.Call.Return(run)
	return _c
}

// PublishIdentityUpdate provides a mock function with given fields: ctx, inboxId, identityUpdate
func (_m *MockIMessagePublisher) PublishIdentityUpdate(ctx context.Context, inboxId [32]byte, identityUpdate []byte) error {
	ret := _m.Called(ctx, inboxId, identityUpdate)

	if len(ret) == 0 {
		panic("no return value specified for PublishIdentityUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, [32]byte, []byte) error); ok {
		r0 = rf(ctx, inboxId, identityUpdate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIMessagePublisher_PublishIdentityUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishIdentityUpdate'
type MockIMessagePublisher_PublishIdentityUpdate_Call struct {
	*mock.Call
}

// PublishIdentityUpdate is a helper method to define mock.On call
//   - ctx context.Context
//   - inboxId [32]byte
//   - identityUpdate []byte
func (_e *MockIMessagePublisher_Expecter) PublishIdentityUpdate(ctx interface{}, inboxId interface{}, identityUpdate interface{}) *MockIMessagePublisher_PublishIdentityUpdate_Call {
	return &MockIMessagePublisher_PublishIdentityUpdate_Call{Call: _e.mock.On("PublishIdentityUpdate", ctx, inboxId, identityUpdate)}
}

func (_c *MockIMessagePublisher_PublishIdentityUpdate_Call) Run(run func(ctx context.Context, inboxId [32]byte, identityUpdate []byte)) *MockIMessagePublisher_PublishIdentityUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([32]byte), args[2].([]byte))
	})
	return _c
}

func (_c *MockIMessagePublisher_PublishIdentityUpdate_Call) Return(_a0 error) *MockIMessagePublisher_PublishIdentityUpdate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIMessagePublisher_PublishIdentityUpdate_Call) RunAndReturn(run func(context.Context, [32]byte, []byte) error) *MockIMessagePublisher_PublishIdentityUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIMessagePublisher creates a new instance of MockIMessagePublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIMessagePublisher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIMessagePublisher {
	mock := &MockIMessagePublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
