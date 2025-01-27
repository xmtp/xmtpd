// Code generated by mockery v2.50.0. DO NOT EDIT.

package registry

import (
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	mock "github.com/stretchr/testify/mock"

	nodes "github.com/xmtp/xmtpd/contracts/pkg/nodes"
)

// MockNodesContract is an autogenerated mock type for the NodesContract type
type MockNodesContract struct {
	mock.Mock
}

type MockNodesContract_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNodesContract) EXPECT() *MockNodesContract_Expecter {
	return &MockNodesContract_Expecter{mock: &_m.Mock}
}

// AllNodes provides a mock function with given fields: opts
func (_m *MockNodesContract) AllNodes(opts *bind.CallOpts) ([]nodes.NodesNodeWithId, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for AllNodes")
	}

	var r0 []nodes.NodesNodeWithId
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) ([]nodes.NodesNodeWithId, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) []nodes.NodesNodeWithId); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]nodes.NodesNodeWithId)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNodesContract_AllNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AllNodes'
type MockNodesContract_AllNodes_Call struct {
	*mock.Call
}

// AllNodes is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *MockNodesContract_Expecter) AllNodes(opts interface{}) *MockNodesContract_AllNodes_Call {
	return &MockNodesContract_AllNodes_Call{Call: _e.mock.On("AllNodes", opts)}
}

func (_c *MockNodesContract_AllNodes_Call) Run(run func(opts *bind.CallOpts)) *MockNodesContract_AllNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *MockNodesContract_AllNodes_Call) Return(_a0 []nodes.NodesNodeWithId, _a1 error) *MockNodesContract_AllNodes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNodesContract_AllNodes_Call) RunAndReturn(run func(*bind.CallOpts) ([]nodes.NodesNodeWithId, error)) *MockNodesContract_AllNodes_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNodesContract creates a new instance of MockNodesContract. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNodesContract(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNodesContract {
	mock := &MockNodesContract{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
