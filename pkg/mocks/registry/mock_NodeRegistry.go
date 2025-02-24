// Code generated by mockery v2.52.2. DO NOT EDIT.

package registry

import (
	mock "github.com/stretchr/testify/mock"
	registry "github.com/xmtp/xmtpd/pkg/registry"
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

// GetNode provides a mock function with given fields: _a0
func (_m *MockNodeRegistry) GetNode(_a0 uint32) (*registry.Node, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetNode")
	}

	var r0 *registry.Node
	var r1 error
	if rf, ok := ret.Get(0).(func(uint32) (*registry.Node, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(uint32) *registry.Node); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*registry.Node)
		}
	}

	if rf, ok := ret.Get(1).(func(uint32) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNodeRegistry_GetNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNode'
type MockNodeRegistry_GetNode_Call struct {
	*mock.Call
}

// GetNode is a helper method to define mock.On call
//   - _a0 uint32
func (_e *MockNodeRegistry_Expecter) GetNode(_a0 interface{}) *MockNodeRegistry_GetNode_Call {
	return &MockNodeRegistry_GetNode_Call{Call: _e.mock.On("GetNode", _a0)}
}

func (_c *MockNodeRegistry_GetNode_Call) Run(run func(_a0 uint32)) *MockNodeRegistry_GetNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint32))
	})
	return _c
}

func (_c *MockNodeRegistry_GetNode_Call) Return(_a0 *registry.Node, _a1 error) *MockNodeRegistry_GetNode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNodeRegistry_GetNode_Call) RunAndReturn(run func(uint32) (*registry.Node, error)) *MockNodeRegistry_GetNode_Call {
	_c.Call.Return(run)
	return _c
}

// GetNodes provides a mock function with no fields
func (_m *MockNodeRegistry) GetNodes() ([]registry.Node, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetNodes")
	}

	var r0 []registry.Node
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]registry.Node, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []registry.Node); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]registry.Node)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNodeRegistry_GetNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNodes'
type MockNodeRegistry_GetNodes_Call struct {
	*mock.Call
}

// GetNodes is a helper method to define mock.On call
func (_e *MockNodeRegistry_Expecter) GetNodes() *MockNodeRegistry_GetNodes_Call {
	return &MockNodeRegistry_GetNodes_Call{Call: _e.mock.On("GetNodes")}
}

func (_c *MockNodeRegistry_GetNodes_Call) Run(run func()) *MockNodeRegistry_GetNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNodeRegistry_GetNodes_Call) Return(_a0 []registry.Node, _a1 error) *MockNodeRegistry_GetNodes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNodeRegistry_GetNodes_Call) RunAndReturn(run func() ([]registry.Node, error)) *MockNodeRegistry_GetNodes_Call {
	_c.Call.Return(run)
	return _c
}

// OnChangedNode provides a mock function with given fields: _a0
func (_m *MockNodeRegistry) OnChangedNode(_a0 uint32) (<-chan registry.Node, registry.CancelSubscription) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for OnChangedNode")
	}

	var r0 <-chan registry.Node
	var r1 registry.CancelSubscription
	if rf, ok := ret.Get(0).(func(uint32) (<-chan registry.Node, registry.CancelSubscription)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(uint32) <-chan registry.Node); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan registry.Node)
		}
	}

	if rf, ok := ret.Get(1).(func(uint32) registry.CancelSubscription); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(registry.CancelSubscription)
		}
	}

	return r0, r1
}

// MockNodeRegistry_OnChangedNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnChangedNode'
type MockNodeRegistry_OnChangedNode_Call struct {
	*mock.Call
}

// OnChangedNode is a helper method to define mock.On call
//   - _a0 uint32
func (_e *MockNodeRegistry_Expecter) OnChangedNode(_a0 interface{}) *MockNodeRegistry_OnChangedNode_Call {
	return &MockNodeRegistry_OnChangedNode_Call{Call: _e.mock.On("OnChangedNode", _a0)}
}

func (_c *MockNodeRegistry_OnChangedNode_Call) Run(run func(_a0 uint32)) *MockNodeRegistry_OnChangedNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint32))
	})
	return _c
}

func (_c *MockNodeRegistry_OnChangedNode_Call) Return(_a0 <-chan registry.Node, _a1 registry.CancelSubscription) *MockNodeRegistry_OnChangedNode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNodeRegistry_OnChangedNode_Call) RunAndReturn(run func(uint32) (<-chan registry.Node, registry.CancelSubscription)) *MockNodeRegistry_OnChangedNode_Call {
	_c.Call.Return(run)
	return _c
}

// OnNewNodes provides a mock function with no fields
func (_m *MockNodeRegistry) OnNewNodes() (<-chan []registry.Node, registry.CancelSubscription) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for OnNewNodes")
	}

	var r0 <-chan []registry.Node
	var r1 registry.CancelSubscription
	if rf, ok := ret.Get(0).(func() (<-chan []registry.Node, registry.CancelSubscription)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() <-chan []registry.Node); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan []registry.Node)
		}
	}

	if rf, ok := ret.Get(1).(func() registry.CancelSubscription); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(registry.CancelSubscription)
		}
	}

	return r0, r1
}

// MockNodeRegistry_OnNewNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnNewNodes'
type MockNodeRegistry_OnNewNodes_Call struct {
	*mock.Call
}

// OnNewNodes is a helper method to define mock.On call
func (_e *MockNodeRegistry_Expecter) OnNewNodes() *MockNodeRegistry_OnNewNodes_Call {
	return &MockNodeRegistry_OnNewNodes_Call{Call: _e.mock.On("OnNewNodes")}
}

func (_c *MockNodeRegistry_OnNewNodes_Call) Run(run func()) *MockNodeRegistry_OnNewNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNodeRegistry_OnNewNodes_Call) Return(_a0 <-chan []registry.Node, _a1 registry.CancelSubscription) *MockNodeRegistry_OnNewNodes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNodeRegistry_OnNewNodes_Call) RunAndReturn(run func() (<-chan []registry.Node, registry.CancelSubscription)) *MockNodeRegistry_OnNewNodes_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with no fields
func (_m *MockNodeRegistry) Stop() {
	_m.Called()
}

// MockNodeRegistry_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type MockNodeRegistry_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *MockNodeRegistry_Expecter) Stop() *MockNodeRegistry_Stop_Call {
	return &MockNodeRegistry_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *MockNodeRegistry_Stop_Call) Run(run func()) *MockNodeRegistry_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNodeRegistry_Stop_Call) Return() *MockNodeRegistry_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockNodeRegistry_Stop_Call) RunAndReturn(run func()) *MockNodeRegistry_Stop_Call {
	_c.Run(run)
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
