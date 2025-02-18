// Code generated by mockery v2.52.2. DO NOT EDIT.

package storer

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	storer "github.com/xmtp/xmtpd/pkg/indexer/storer"

	types "github.com/ethereum/go-ethereum/core/types"
)

// MockLogStorer is an autogenerated mock type for the LogStorer type
type MockLogStorer struct {
	mock.Mock
}

type MockLogStorer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLogStorer) EXPECT() *MockLogStorer_Expecter {
	return &MockLogStorer_Expecter{mock: &_m.Mock}
}

// StoreLog provides a mock function with given fields: ctx, event
func (_m *MockLogStorer) StoreLog(ctx context.Context, event types.Log) storer.LogStorageError {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for StoreLog")
	}

	var r0 storer.LogStorageError
	if rf, ok := ret.Get(0).(func(context.Context, types.Log) storer.LogStorageError); ok {
		r0 = rf(ctx, event)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(storer.LogStorageError)
		}
	}

	return r0
}

// MockLogStorer_StoreLog_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreLog'
type MockLogStorer_StoreLog_Call struct {
	*mock.Call
}

// StoreLog is a helper method to define mock.On call
//   - ctx context.Context
//   - event types.Log
func (_e *MockLogStorer_Expecter) StoreLog(ctx interface{}, event interface{}) *MockLogStorer_StoreLog_Call {
	return &MockLogStorer_StoreLog_Call{Call: _e.mock.On("StoreLog", ctx, event)}
}

func (_c *MockLogStorer_StoreLog_Call) Run(run func(ctx context.Context, event types.Log)) *MockLogStorer_StoreLog_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Log))
	})
	return _c
}

func (_c *MockLogStorer_StoreLog_Call) Return(_a0 storer.LogStorageError) *MockLogStorer_StoreLog_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogStorer_StoreLog_Call) RunAndReturn(run func(context.Context, types.Log) storer.LogStorageError) *MockLogStorer_StoreLog_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLogStorer creates a new instance of MockLogStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLogStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLogStorer {
	mock := &MockLogStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
