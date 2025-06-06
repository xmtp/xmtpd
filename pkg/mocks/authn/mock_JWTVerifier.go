// Code generated by mockery v2.52.2. DO NOT EDIT.

package authn

import (
	mock "github.com/stretchr/testify/mock"
	authn "github.com/xmtp/xmtpd/pkg/authn"
)

// MockJWTVerifier is an autogenerated mock type for the JWTVerifier type
type MockJWTVerifier struct {
	mock.Mock
}

type MockJWTVerifier_Expecter struct {
	mock *mock.Mock
}

func (_m *MockJWTVerifier) EXPECT() *MockJWTVerifier_Expecter {
	return &MockJWTVerifier_Expecter{mock: &_m.Mock}
}

// Verify provides a mock function with given fields: tokenString
func (_m *MockJWTVerifier) Verify(tokenString string) (uint32, authn.CloseFunc, error) {
	ret := _m.Called(tokenString)

	if len(ret) == 0 {
		panic("no return value specified for Verify")
	}

	var r0 uint32
	var r1 authn.CloseFunc
	var r2 error
	if rf, ok := ret.Get(0).(func(string) (uint32, authn.CloseFunc, error)); ok {
		return rf(tokenString)
	}
	if rf, ok := ret.Get(0).(func(string) uint32); ok {
		r0 = rf(tokenString)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(string) authn.CloseFunc); ok {
		r1 = rf(tokenString)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(authn.CloseFunc)
		}
	}

	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(tokenString)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockJWTVerifier_Verify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Verify'
type MockJWTVerifier_Verify_Call struct {
	*mock.Call
}

// Verify is a helper method to define mock.On call
//   - tokenString string
func (_e *MockJWTVerifier_Expecter) Verify(tokenString interface{}) *MockJWTVerifier_Verify_Call {
	return &MockJWTVerifier_Verify_Call{Call: _e.mock.On("Verify", tokenString)}
}

func (_c *MockJWTVerifier_Verify_Call) Run(run func(tokenString string)) *MockJWTVerifier_Verify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockJWTVerifier_Verify_Call) Return(_a0 uint32, _a1 authn.CloseFunc, _a2 error) *MockJWTVerifier_Verify_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockJWTVerifier_Verify_Call) RunAndReturn(run func(string) (uint32, authn.CloseFunc, error)) *MockJWTVerifier_Verify_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockJWTVerifier creates a new instance of MockJWTVerifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockJWTVerifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJWTVerifier {
	mock := &MockJWTVerifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
