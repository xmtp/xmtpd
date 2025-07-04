// Code generated by mockery v2.52.2. DO NOT EDIT.

package payerreport

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	envelopes "github.com/xmtp/xmtpd/pkg/envelopes"

	mock "github.com/stretchr/testify/mock"

	payerreport "github.com/xmtp/xmtpd/pkg/payerreport"

	queries "github.com/xmtp/xmtpd/pkg/db/queries"
)

// MockIPayerReportStore is an autogenerated mock type for the IPayerReportStore type
type MockIPayerReportStore struct {
	mock.Mock
}

type MockIPayerReportStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIPayerReportStore) EXPECT() *MockIPayerReportStore_Expecter {
	return &MockIPayerReportStore_Expecter{mock: &_m.Mock}
}

// CreateAttestation provides a mock function with given fields: ctx, attestation, payerEnvelope
func (_m *MockIPayerReportStore) CreateAttestation(ctx context.Context, attestation *payerreport.PayerReportAttestation, payerEnvelope *envelopes.PayerEnvelope) error {
	ret := _m.Called(ctx, attestation, payerEnvelope)

	if len(ret) == 0 {
		panic("no return value specified for CreateAttestation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *payerreport.PayerReportAttestation, *envelopes.PayerEnvelope) error); ok {
		r0 = rf(ctx, attestation, payerEnvelope)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIPayerReportStore_CreateAttestation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAttestation'
type MockIPayerReportStore_CreateAttestation_Call struct {
	*mock.Call
}

// CreateAttestation is a helper method to define mock.On call
//   - ctx context.Context
//   - attestation *payerreport.PayerReportAttestation
//   - payerEnvelope *envelopes.PayerEnvelope
func (_e *MockIPayerReportStore_Expecter) CreateAttestation(ctx interface{}, attestation interface{}, payerEnvelope interface{}) *MockIPayerReportStore_CreateAttestation_Call {
	return &MockIPayerReportStore_CreateAttestation_Call{Call: _e.mock.On("CreateAttestation", ctx, attestation, payerEnvelope)}
}

func (_c *MockIPayerReportStore_CreateAttestation_Call) Run(run func(ctx context.Context, attestation *payerreport.PayerReportAttestation, payerEnvelope *envelopes.PayerEnvelope)) *MockIPayerReportStore_CreateAttestation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*payerreport.PayerReportAttestation), args[2].(*envelopes.PayerEnvelope))
	})
	return _c
}

func (_c *MockIPayerReportStore_CreateAttestation_Call) Return(_a0 error) *MockIPayerReportStore_CreateAttestation_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIPayerReportStore_CreateAttestation_Call) RunAndReturn(run func(context.Context, *payerreport.PayerReportAttestation, *envelopes.PayerEnvelope) error) *MockIPayerReportStore_CreateAttestation_Call {
	_c.Call.Return(run)
	return _c
}

// CreatePayerReport provides a mock function with given fields: ctx, report, payerEnvelope
func (_m *MockIPayerReportStore) CreatePayerReport(ctx context.Context, report *payerreport.PayerReport, payerEnvelope *envelopes.PayerEnvelope) (*payerreport.ReportID, error) {
	ret := _m.Called(ctx, report, payerEnvelope)

	if len(ret) == 0 {
		panic("no return value specified for CreatePayerReport")
	}

	var r0 *payerreport.ReportID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payerreport.PayerReport, *envelopes.PayerEnvelope) (*payerreport.ReportID, error)); ok {
		return rf(ctx, report, payerEnvelope)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payerreport.PayerReport, *envelopes.PayerEnvelope) *payerreport.ReportID); ok {
		r0 = rf(ctx, report, payerEnvelope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payerreport.ReportID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payerreport.PayerReport, *envelopes.PayerEnvelope) error); ok {
		r1 = rf(ctx, report, payerEnvelope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIPayerReportStore_CreatePayerReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePayerReport'
type MockIPayerReportStore_CreatePayerReport_Call struct {
	*mock.Call
}

// CreatePayerReport is a helper method to define mock.On call
//   - ctx context.Context
//   - report *payerreport.PayerReport
//   - payerEnvelope *envelopes.PayerEnvelope
func (_e *MockIPayerReportStore_Expecter) CreatePayerReport(ctx interface{}, report interface{}, payerEnvelope interface{}) *MockIPayerReportStore_CreatePayerReport_Call {
	return &MockIPayerReportStore_CreatePayerReport_Call{Call: _e.mock.On("CreatePayerReport", ctx, report, payerEnvelope)}
}

func (_c *MockIPayerReportStore_CreatePayerReport_Call) Run(run func(ctx context.Context, report *payerreport.PayerReport, payerEnvelope *envelopes.PayerEnvelope)) *MockIPayerReportStore_CreatePayerReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*payerreport.PayerReport), args[2].(*envelopes.PayerEnvelope))
	})
	return _c
}

func (_c *MockIPayerReportStore_CreatePayerReport_Call) Return(_a0 *payerreport.ReportID, _a1 error) *MockIPayerReportStore_CreatePayerReport_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIPayerReportStore_CreatePayerReport_Call) RunAndReturn(run func(context.Context, *payerreport.PayerReport, *envelopes.PayerEnvelope) (*payerreport.ReportID, error)) *MockIPayerReportStore_CreatePayerReport_Call {
	_c.Call.Return(run)
	return _c
}

// FetchReport provides a mock function with given fields: ctx, id
func (_m *MockIPayerReportStore) FetchReport(ctx context.Context, id payerreport.ReportID) (*payerreport.PayerReportWithStatus, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FetchReport")
	}

	var r0 *payerreport.PayerReportWithStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, payerreport.ReportID) (*payerreport.PayerReportWithStatus, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, payerreport.ReportID) *payerreport.PayerReportWithStatus); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payerreport.PayerReportWithStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, payerreport.ReportID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIPayerReportStore_FetchReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchReport'
type MockIPayerReportStore_FetchReport_Call struct {
	*mock.Call
}

// FetchReport is a helper method to define mock.On call
//   - ctx context.Context
//   - id payerreport.ReportID
func (_e *MockIPayerReportStore_Expecter) FetchReport(ctx interface{}, id interface{}) *MockIPayerReportStore_FetchReport_Call {
	return &MockIPayerReportStore_FetchReport_Call{Call: _e.mock.On("FetchReport", ctx, id)}
}

func (_c *MockIPayerReportStore_FetchReport_Call) Run(run func(ctx context.Context, id payerreport.ReportID)) *MockIPayerReportStore_FetchReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(payerreport.ReportID))
	})
	return _c
}

func (_c *MockIPayerReportStore_FetchReport_Call) Return(_a0 *payerreport.PayerReportWithStatus, _a1 error) *MockIPayerReportStore_FetchReport_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIPayerReportStore_FetchReport_Call) RunAndReturn(run func(context.Context, payerreport.ReportID) (*payerreport.PayerReportWithStatus, error)) *MockIPayerReportStore_FetchReport_Call {
	_c.Call.Return(run)
	return _c
}

// FetchReports provides a mock function with given fields: ctx, query
func (_m *MockIPayerReportStore) FetchReports(ctx context.Context, query *payerreport.FetchReportsQuery) ([]*payerreport.PayerReportWithStatus, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for FetchReports")
	}

	var r0 []*payerreport.PayerReportWithStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payerreport.FetchReportsQuery) ([]*payerreport.PayerReportWithStatus, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payerreport.FetchReportsQuery) []*payerreport.PayerReportWithStatus); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*payerreport.PayerReportWithStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payerreport.FetchReportsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIPayerReportStore_FetchReports_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchReports'
type MockIPayerReportStore_FetchReports_Call struct {
	*mock.Call
}

// FetchReports is a helper method to define mock.On call
//   - ctx context.Context
//   - query *payerreport.FetchReportsQuery
func (_e *MockIPayerReportStore_Expecter) FetchReports(ctx interface{}, query interface{}) *MockIPayerReportStore_FetchReports_Call {
	return &MockIPayerReportStore_FetchReports_Call{Call: _e.mock.On("FetchReports", ctx, query)}
}

func (_c *MockIPayerReportStore_FetchReports_Call) Run(run func(ctx context.Context, query *payerreport.FetchReportsQuery)) *MockIPayerReportStore_FetchReports_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*payerreport.FetchReportsQuery))
	})
	return _c
}

func (_c *MockIPayerReportStore_FetchReports_Call) Return(_a0 []*payerreport.PayerReportWithStatus, _a1 error) *MockIPayerReportStore_FetchReports_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIPayerReportStore_FetchReports_Call) RunAndReturn(run func(context.Context, *payerreport.FetchReportsQuery) ([]*payerreport.PayerReportWithStatus, error)) *MockIPayerReportStore_FetchReports_Call {
	_c.Call.Return(run)
	return _c
}

// Queries provides a mock function with no fields
func (_m *MockIPayerReportStore) Queries() *queries.Queries {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Queries")
	}

	var r0 *queries.Queries
	if rf, ok := ret.Get(0).(func() *queries.Queries); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*queries.Queries)
		}
	}

	return r0
}

// MockIPayerReportStore_Queries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Queries'
type MockIPayerReportStore_Queries_Call struct {
	*mock.Call
}

// Queries is a helper method to define mock.On call
func (_e *MockIPayerReportStore_Expecter) Queries() *MockIPayerReportStore_Queries_Call {
	return &MockIPayerReportStore_Queries_Call{Call: _e.mock.On("Queries")}
}

func (_c *MockIPayerReportStore_Queries_Call) Run(run func()) *MockIPayerReportStore_Queries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockIPayerReportStore_Queries_Call) Return(_a0 *queries.Queries) *MockIPayerReportStore_Queries_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIPayerReportStore_Queries_Call) RunAndReturn(run func() *queries.Queries) *MockIPayerReportStore_Queries_Call {
	_c.Call.Return(run)
	return _c
}

// SetReportAttestationStatus provides a mock function with given fields: ctx, id, fromStatus, toStatus
func (_m *MockIPayerReportStore) SetReportAttestationStatus(ctx context.Context, id payerreport.ReportID, fromStatus []payerreport.AttestationStatus, toStatus payerreport.AttestationStatus) error {
	ret := _m.Called(ctx, id, fromStatus, toStatus)

	if len(ret) == 0 {
		panic("no return value specified for SetReportAttestationStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, payerreport.ReportID, []payerreport.AttestationStatus, payerreport.AttestationStatus) error); ok {
		r0 = rf(ctx, id, fromStatus, toStatus)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIPayerReportStore_SetReportAttestationStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetReportAttestationStatus'
type MockIPayerReportStore_SetReportAttestationStatus_Call struct {
	*mock.Call
}

// SetReportAttestationStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - id payerreport.ReportID
//   - fromStatus []payerreport.AttestationStatus
//   - toStatus payerreport.AttestationStatus
func (_e *MockIPayerReportStore_Expecter) SetReportAttestationStatus(ctx interface{}, id interface{}, fromStatus interface{}, toStatus interface{}) *MockIPayerReportStore_SetReportAttestationStatus_Call {
	return &MockIPayerReportStore_SetReportAttestationStatus_Call{Call: _e.mock.On("SetReportAttestationStatus", ctx, id, fromStatus, toStatus)}
}

func (_c *MockIPayerReportStore_SetReportAttestationStatus_Call) Run(run func(ctx context.Context, id payerreport.ReportID, fromStatus []payerreport.AttestationStatus, toStatus payerreport.AttestationStatus)) *MockIPayerReportStore_SetReportAttestationStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(payerreport.ReportID), args[2].([]payerreport.AttestationStatus), args[3].(payerreport.AttestationStatus))
	})
	return _c
}

func (_c *MockIPayerReportStore_SetReportAttestationStatus_Call) Return(_a0 error) *MockIPayerReportStore_SetReportAttestationStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIPayerReportStore_SetReportAttestationStatus_Call) RunAndReturn(run func(context.Context, payerreport.ReportID, []payerreport.AttestationStatus, payerreport.AttestationStatus) error) *MockIPayerReportStore_SetReportAttestationStatus_Call {
	_c.Call.Return(run)
	return _c
}

// StoreSyncedAttestation provides a mock function with given fields: ctx, envelope, payerID
func (_m *MockIPayerReportStore) StoreSyncedAttestation(ctx context.Context, envelope *envelopes.OriginatorEnvelope, payerID int32) error {
	ret := _m.Called(ctx, envelope, payerID)

	if len(ret) == 0 {
		panic("no return value specified for StoreSyncedAttestation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *envelopes.OriginatorEnvelope, int32) error); ok {
		r0 = rf(ctx, envelope, payerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIPayerReportStore_StoreSyncedAttestation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreSyncedAttestation'
type MockIPayerReportStore_StoreSyncedAttestation_Call struct {
	*mock.Call
}

// StoreSyncedAttestation is a helper method to define mock.On call
//   - ctx context.Context
//   - envelope *envelopes.OriginatorEnvelope
//   - payerID int32
func (_e *MockIPayerReportStore_Expecter) StoreSyncedAttestation(ctx interface{}, envelope interface{}, payerID interface{}) *MockIPayerReportStore_StoreSyncedAttestation_Call {
	return &MockIPayerReportStore_StoreSyncedAttestation_Call{Call: _e.mock.On("StoreSyncedAttestation", ctx, envelope, payerID)}
}

func (_c *MockIPayerReportStore_StoreSyncedAttestation_Call) Run(run func(ctx context.Context, envelope *envelopes.OriginatorEnvelope, payerID int32)) *MockIPayerReportStore_StoreSyncedAttestation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*envelopes.OriginatorEnvelope), args[2].(int32))
	})
	return _c
}

func (_c *MockIPayerReportStore_StoreSyncedAttestation_Call) Return(_a0 error) *MockIPayerReportStore_StoreSyncedAttestation_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIPayerReportStore_StoreSyncedAttestation_Call) RunAndReturn(run func(context.Context, *envelopes.OriginatorEnvelope, int32) error) *MockIPayerReportStore_StoreSyncedAttestation_Call {
	_c.Call.Return(run)
	return _c
}

// StoreSyncedReport provides a mock function with given fields: ctx, envelope, payerID, domainSeparator
func (_m *MockIPayerReportStore) StoreSyncedReport(ctx context.Context, envelope *envelopes.OriginatorEnvelope, payerID int32, domainSeparator common.Hash) error {
	ret := _m.Called(ctx, envelope, payerID, domainSeparator)

	if len(ret) == 0 {
		panic("no return value specified for StoreSyncedReport")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *envelopes.OriginatorEnvelope, int32, common.Hash) error); ok {
		r0 = rf(ctx, envelope, payerID, domainSeparator)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIPayerReportStore_StoreSyncedReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreSyncedReport'
type MockIPayerReportStore_StoreSyncedReport_Call struct {
	*mock.Call
}

// StoreSyncedReport is a helper method to define mock.On call
//   - ctx context.Context
//   - envelope *envelopes.OriginatorEnvelope
//   - payerID int32
//   - domainSeparator common.Hash
func (_e *MockIPayerReportStore_Expecter) StoreSyncedReport(ctx interface{}, envelope interface{}, payerID interface{}, domainSeparator interface{}) *MockIPayerReportStore_StoreSyncedReport_Call {
	return &MockIPayerReportStore_StoreSyncedReport_Call{Call: _e.mock.On("StoreSyncedReport", ctx, envelope, payerID, domainSeparator)}
}

func (_c *MockIPayerReportStore_StoreSyncedReport_Call) Run(run func(ctx context.Context, envelope *envelopes.OriginatorEnvelope, payerID int32, domainSeparator common.Hash)) *MockIPayerReportStore_StoreSyncedReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*envelopes.OriginatorEnvelope), args[2].(int32), args[3].(common.Hash))
	})
	return _c
}

func (_c *MockIPayerReportStore_StoreSyncedReport_Call) Return(_a0 error) *MockIPayerReportStore_StoreSyncedReport_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIPayerReportStore_StoreSyncedReport_Call) RunAndReturn(run func(context.Context, *envelopes.OriginatorEnvelope, int32, common.Hash) error) *MockIPayerReportStore_StoreSyncedReport_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIPayerReportStore creates a new instance of MockIPayerReportStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIPayerReportStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIPayerReportStore {
	mock := &MockIPayerReportStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
