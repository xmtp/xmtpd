// Code generated by mockery v2.52.2. DO NOT EDIT.

package blockchain

import (
	big "math/big"

	context "context"

	ecdsa "crypto/ecdsa"

	mock "github.com/stretchr/testify/mock"
)

// MockNodeRegistry is an autogenerated mock type for the NodeRegistry type
type MockNodeRegistry struct {
	mock.Mock
}

type MockNodeRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNodeRegistry) EXPECT() *MockNodeRegistry_Expecter {
	return &MockNodeRegistry_Expecter{mock: &_m.Mock}
}

// AddNode provides a mock function with given fields: ctx, owner, signingKeyPub, httpAddress, minMonthlyFee
func (_m *MockNodeRegistry) AddNode(ctx context.Context, owner string, signingKeyPub *ecdsa.PublicKey, httpAddress string, minMonthlyFee *big.Int) (uint32, error) {
	ret := _m.Called(ctx, owner, signingKeyPub, httpAddress, minMonthlyFee)

	if len(ret) == 0 {
		panic("no return value specified for AddNode")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *ecdsa.PublicKey, string, *big.Int) (uint32, error)); ok {
		return rf(ctx, owner, signingKeyPub, httpAddress, minMonthlyFee)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *ecdsa.PublicKey, string, *big.Int) uint32); ok {
		r0 = rf(ctx, owner, signingKeyPub, httpAddress, minMonthlyFee)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *ecdsa.PublicKey, string, *big.Int) error); ok {
		r1 = rf(ctx, owner, signingKeyPub, httpAddress, minMonthlyFee)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNodeRegistry_AddNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddNode'
type MockNodeRegistry_AddNode_Call struct {
	*mock.Call
}

// AddNode is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - signingKeyPub *ecdsa.PublicKey
//   - httpAddress string
//   - minMonthlyFee *big.Int
func (_e *MockNodeRegistry_Expecter) AddNode(ctx interface{}, owner interface{}, signingKeyPub interface{}, httpAddress interface{}, minMonthlyFee interface{}) *MockNodeRegistry_AddNode_Call {
	return &MockNodeRegistry_AddNode_Call{Call: _e.mock.On("AddNode", ctx, owner, signingKeyPub, httpAddress, minMonthlyFee)}
}

func (_c *MockNodeRegistry_AddNode_Call) Run(run func(ctx context.Context, owner string, signingKeyPub *ecdsa.PublicKey, httpAddress string, minMonthlyFee *big.Int)) *MockNodeRegistry_AddNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*ecdsa.PublicKey), args[3].(string), args[4].(*big.Int))
	})
	return _c
}

func (_c *MockNodeRegistry_AddNode_Call) Return(_a0 uint32, _a1 error) *MockNodeRegistry_AddNode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNodeRegistry_AddNode_Call) RunAndReturn(run func(context.Context, string, *ecdsa.PublicKey, string, *big.Int) (uint32, error)) *MockNodeRegistry_AddNode_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateActive provides a mock function with given fields: ctx, nodeId, isActive
func (_m *MockNodeRegistry) UpdateActive(ctx context.Context, nodeId uint32, isActive bool) error {
	ret := _m.Called(ctx, nodeId, isActive)

	if len(ret) == 0 {
		panic("no return value specified for UpdateActive")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, bool) error); ok {
		r0 = rf(ctx, nodeId, isActive)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockNodeRegistry_UpdateActive_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateActive'
type MockNodeRegistry_UpdateActive_Call struct {
	*mock.Call
}

// UpdateActive is a helper method to define mock.On call
//   - ctx context.Context
//   - nodeId uint32
//   - isActive bool
func (_e *MockNodeRegistry_Expecter) UpdateActive(ctx interface{}, nodeId interface{}, isActive interface{}) *MockNodeRegistry_UpdateActive_Call {
	return &MockNodeRegistry_UpdateActive_Call{Call: _e.mock.On("UpdateActive", ctx, nodeId, isActive)}
}

func (_c *MockNodeRegistry_UpdateActive_Call) Run(run func(ctx context.Context, nodeId uint32, isActive bool)) *MockNodeRegistry_UpdateActive_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32), args[2].(bool))
	})
	return _c
}

func (_c *MockNodeRegistry_UpdateActive_Call) Return(_a0 error) *MockNodeRegistry_UpdateActive_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNodeRegistry_UpdateActive_Call) RunAndReturn(run func(context.Context, uint32, bool) error) *MockNodeRegistry_UpdateActive_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateIsApiEnabled provides a mock function with given fields: ctx, nodeId
func (_m *MockNodeRegistry) UpdateIsApiEnabled(ctx context.Context, nodeId uint32) error {
	ret := _m.Called(ctx, nodeId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateIsApiEnabled")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) error); ok {
		r0 = rf(ctx, nodeId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockNodeRegistry_UpdateIsApiEnabled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateIsApiEnabled'
type MockNodeRegistry_UpdateIsApiEnabled_Call struct {
	*mock.Call
}

// UpdateIsApiEnabled is a helper method to define mock.On call
//   - ctx context.Context
//   - nodeId uint32
func (_e *MockNodeRegistry_Expecter) UpdateIsApiEnabled(ctx interface{}, nodeId interface{}) *MockNodeRegistry_UpdateIsApiEnabled_Call {
	return &MockNodeRegistry_UpdateIsApiEnabled_Call{Call: _e.mock.On("UpdateIsApiEnabled", ctx, nodeId)}
}

func (_c *MockNodeRegistry_UpdateIsApiEnabled_Call) Run(run func(ctx context.Context, nodeId uint32)) *MockNodeRegistry_UpdateIsApiEnabled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32))
	})
	return _c
}

func (_c *MockNodeRegistry_UpdateIsApiEnabled_Call) Return(_a0 error) *MockNodeRegistry_UpdateIsApiEnabled_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNodeRegistry_UpdateIsApiEnabled_Call) RunAndReturn(run func(context.Context, uint32) error) *MockNodeRegistry_UpdateIsApiEnabled_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateIsReplicationEnabled provides a mock function with given fields: ctx, nodeId, isReplicationEnabled
func (_m *MockNodeRegistry) UpdateIsReplicationEnabled(ctx context.Context, nodeId uint32, isReplicationEnabled bool) error {
	ret := _m.Called(ctx, nodeId, isReplicationEnabled)

	if len(ret) == 0 {
		panic("no return value specified for UpdateIsReplicationEnabled")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, bool) error); ok {
		r0 = rf(ctx, nodeId, isReplicationEnabled)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockNodeRegistry_UpdateIsReplicationEnabled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateIsReplicationEnabled'
type MockNodeRegistry_UpdateIsReplicationEnabled_Call struct {
	*mock.Call
}

// UpdateIsReplicationEnabled is a helper method to define mock.On call
//   - ctx context.Context
//   - nodeId uint32
//   - isReplicationEnabled bool
func (_e *MockNodeRegistry_Expecter) UpdateIsReplicationEnabled(ctx interface{}, nodeId interface{}, isReplicationEnabled interface{}) *MockNodeRegistry_UpdateIsReplicationEnabled_Call {
	return &MockNodeRegistry_UpdateIsReplicationEnabled_Call{Call: _e.mock.On("UpdateIsReplicationEnabled", ctx, nodeId, isReplicationEnabled)}
}

func (_c *MockNodeRegistry_UpdateIsReplicationEnabled_Call) Run(run func(ctx context.Context, nodeId uint32, isReplicationEnabled bool)) *MockNodeRegistry_UpdateIsReplicationEnabled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32), args[2].(bool))
	})
	return _c
}

func (_c *MockNodeRegistry_UpdateIsReplicationEnabled_Call) Return(_a0 error) *MockNodeRegistry_UpdateIsReplicationEnabled_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNodeRegistry_UpdateIsReplicationEnabled_Call) RunAndReturn(run func(context.Context, uint32, bool) error) *MockNodeRegistry_UpdateIsReplicationEnabled_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNodeRegistry creates a new instance of MockNodeRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNodeRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNodeRegistry {
	mock := &MockNodeRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
