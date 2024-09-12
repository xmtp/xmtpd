// Code generated by mockery v2.44.1. DO NOT EDIT.

package mlsvalidate

import (
	associations "github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"

	context "context"

	mlsvalidate "github.com/xmtp/xmtpd/pkg/mlsvalidate"

	mock "github.com/stretchr/testify/mock"
)

// MockMLSValidationService is an autogenerated mock type for the MLSValidationService type
type MockMLSValidationService struct {
	mock.Mock
}

type MockMLSValidationService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMLSValidationService) EXPECT() *MockMLSValidationService_Expecter {
	return &MockMLSValidationService_Expecter{mock: &_m.Mock}
}

// GetAssociationState provides a mock function with given fields: ctx, oldUpdates, newUpdates
func (_m *MockMLSValidationService) GetAssociationState(ctx context.Context, oldUpdates []*associations.IdentityUpdate, newUpdates []*associations.IdentityUpdate) (*mlsvalidate.AssociationStateResult, error) {
	ret := _m.Called(ctx, oldUpdates, newUpdates)

	if len(ret) == 0 {
		panic("no return value specified for GetAssociationState")
	}

	var r0 *mlsvalidate.AssociationStateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []*associations.IdentityUpdate, []*associations.IdentityUpdate) (*mlsvalidate.AssociationStateResult, error)); ok {
		return rf(ctx, oldUpdates, newUpdates)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []*associations.IdentityUpdate, []*associations.IdentityUpdate) *mlsvalidate.AssociationStateResult); ok {
		r0 = rf(ctx, oldUpdates, newUpdates)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mlsvalidate.AssociationStateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []*associations.IdentityUpdate, []*associations.IdentityUpdate) error); ok {
		r1 = rf(ctx, oldUpdates, newUpdates)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMLSValidationService_GetAssociationState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAssociationState'
type MockMLSValidationService_GetAssociationState_Call struct {
	*mock.Call
}

// GetAssociationState is a helper method to define mock.On call
//   - ctx context.Context
//   - oldUpdates []*associations.IdentityUpdate
//   - newUpdates []*associations.IdentityUpdate
func (_e *MockMLSValidationService_Expecter) GetAssociationState(ctx interface{}, oldUpdates interface{}, newUpdates interface{}) *MockMLSValidationService_GetAssociationState_Call {
	return &MockMLSValidationService_GetAssociationState_Call{Call: _e.mock.On("GetAssociationState", ctx, oldUpdates, newUpdates)}
}

func (_c *MockMLSValidationService_GetAssociationState_Call) Run(run func(ctx context.Context, oldUpdates []*associations.IdentityUpdate, newUpdates []*associations.IdentityUpdate)) *MockMLSValidationService_GetAssociationState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*associations.IdentityUpdate), args[2].([]*associations.IdentityUpdate))
	})
	return _c
}

func (_c *MockMLSValidationService_GetAssociationState_Call) Return(_a0 *mlsvalidate.AssociationStateResult, _a1 error) *MockMLSValidationService_GetAssociationState_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMLSValidationService_GetAssociationState_Call) RunAndReturn(run func(context.Context, []*associations.IdentityUpdate, []*associations.IdentityUpdate) (*mlsvalidate.AssociationStateResult, error)) *MockMLSValidationService_GetAssociationState_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateGroupMessages provides a mock function with given fields: ctx, groupMessages
func (_m *MockMLSValidationService) ValidateGroupMessages(ctx context.Context, groupMessages []*apiv1.GroupMessageInput) ([]mlsvalidate.GroupMessageValidationResult, error) {
	ret := _m.Called(ctx, groupMessages)

	if len(ret) == 0 {
		panic("no return value specified for ValidateGroupMessages")
	}

	var r0 []mlsvalidate.GroupMessageValidationResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []*apiv1.GroupMessageInput) ([]mlsvalidate.GroupMessageValidationResult, error)); ok {
		return rf(ctx, groupMessages)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []*apiv1.GroupMessageInput) []mlsvalidate.GroupMessageValidationResult); ok {
		r0 = rf(ctx, groupMessages)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]mlsvalidate.GroupMessageValidationResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []*apiv1.GroupMessageInput) error); ok {
		r1 = rf(ctx, groupMessages)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMLSValidationService_ValidateGroupMessages_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateGroupMessages'
type MockMLSValidationService_ValidateGroupMessages_Call struct {
	*mock.Call
}

// ValidateGroupMessages is a helper method to define mock.On call
//   - ctx context.Context
//   - groupMessages []*apiv1.GroupMessageInput
func (_e *MockMLSValidationService_Expecter) ValidateGroupMessages(ctx interface{}, groupMessages interface{}) *MockMLSValidationService_ValidateGroupMessages_Call {
	return &MockMLSValidationService_ValidateGroupMessages_Call{Call: _e.mock.On("ValidateGroupMessages", ctx, groupMessages)}
}

func (_c *MockMLSValidationService_ValidateGroupMessages_Call) Run(run func(ctx context.Context, groupMessages []*apiv1.GroupMessageInput)) *MockMLSValidationService_ValidateGroupMessages_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*apiv1.GroupMessageInput))
	})
	return _c
}

func (_c *MockMLSValidationService_ValidateGroupMessages_Call) Return(_a0 []mlsvalidate.GroupMessageValidationResult, _a1 error) *MockMLSValidationService_ValidateGroupMessages_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMLSValidationService_ValidateGroupMessages_Call) RunAndReturn(run func(context.Context, []*apiv1.GroupMessageInput) ([]mlsvalidate.GroupMessageValidationResult, error)) *MockMLSValidationService_ValidateGroupMessages_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateKeyPackages provides a mock function with given fields: ctx, keyPackages
func (_m *MockMLSValidationService) ValidateKeyPackages(ctx context.Context, keyPackages [][]byte) ([]mlsvalidate.KeyPackageValidationResult, error) {
	ret := _m.Called(ctx, keyPackages)

	if len(ret) == 0 {
		panic("no return value specified for ValidateKeyPackages")
	}

	var r0 []mlsvalidate.KeyPackageValidationResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, [][]byte) ([]mlsvalidate.KeyPackageValidationResult, error)); ok {
		return rf(ctx, keyPackages)
	}
	if rf, ok := ret.Get(0).(func(context.Context, [][]byte) []mlsvalidate.KeyPackageValidationResult); ok {
		r0 = rf(ctx, keyPackages)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]mlsvalidate.KeyPackageValidationResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, [][]byte) error); ok {
		r1 = rf(ctx, keyPackages)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMLSValidationService_ValidateKeyPackages_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateKeyPackages'
type MockMLSValidationService_ValidateKeyPackages_Call struct {
	*mock.Call
}

// ValidateKeyPackages is a helper method to define mock.On call
//   - ctx context.Context
//   - keyPackages [][]byte
func (_e *MockMLSValidationService_Expecter) ValidateKeyPackages(ctx interface{}, keyPackages interface{}) *MockMLSValidationService_ValidateKeyPackages_Call {
	return &MockMLSValidationService_ValidateKeyPackages_Call{Call: _e.mock.On("ValidateKeyPackages", ctx, keyPackages)}
}

func (_c *MockMLSValidationService_ValidateKeyPackages_Call) Run(run func(ctx context.Context, keyPackages [][]byte)) *MockMLSValidationService_ValidateKeyPackages_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([][]byte))
	})
	return _c
}

func (_c *MockMLSValidationService_ValidateKeyPackages_Call) Return(_a0 []mlsvalidate.KeyPackageValidationResult, _a1 error) *MockMLSValidationService_ValidateKeyPackages_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMLSValidationService_ValidateKeyPackages_Call) RunAndReturn(run func(context.Context, [][]byte) ([]mlsvalidate.KeyPackageValidationResult, error)) *MockMLSValidationService_ValidateKeyPackages_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMLSValidationService creates a new instance of MockMLSValidationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMLSValidationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMLSValidationService {
	mock := &MockMLSValidationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
