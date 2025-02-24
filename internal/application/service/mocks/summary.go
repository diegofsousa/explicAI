// Code generated by mockery v2.52.3. DO NOT EDIT.

package servicemocks

import (
	context "context"

	service "github.com/diegofsousa/explicAI/internal/application/service"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// SummaryUseCase is an autogenerated mock type for the SummaryUseCase type
type SummaryUseCase struct {
	mock.Mock
}

type SummaryUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *SummaryUseCase) EXPECT() *SummaryUseCase_Expecter {
	return &SummaryUseCase_Expecter{mock: &_m.Mock}
}

// CreateSummaryAndTriggerAIProccess provides a mock function with given fields: ctx, audio
func (_m *SummaryUseCase) CreateSummaryAndTriggerAIProccess(ctx context.Context, audio []byte) (*service.SummarySimpleOutput, error) {
	ret := _m.Called(ctx, audio)

	if len(ret) == 0 {
		panic("no return value specified for CreateSummaryAndTriggerAIProccess")
	}

	var r0 *service.SummarySimpleOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) (*service.SummarySimpleOutput, error)); ok {
		return rf(ctx, audio)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte) *service.SummarySimpleOutput); ok {
		r0 = rf(ctx, audio)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.SummarySimpleOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctx, audio)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSummaryAndTriggerAIProccess'
type SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call struct {
	*mock.Call
}

// CreateSummaryAndTriggerAIProccess is a helper method to define mock.On call
//   - ctx context.Context
//   - audio []byte
func (_e *SummaryUseCase_Expecter) CreateSummaryAndTriggerAIProccess(ctx interface{}, audio interface{}) *SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call {
	return &SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call{Call: _e.mock.On("CreateSummaryAndTriggerAIProccess", ctx, audio)}
}

func (_c *SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call) Run(run func(ctx context.Context, audio []byte)) *SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte))
	})
	return _c
}

func (_c *SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call) Return(_a0 *service.SummarySimpleOutput, _a1 error) *SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call) RunAndReturn(run func(context.Context, []byte) (*service.SummarySimpleOutput, error)) *SummaryUseCase_CreateSummaryAndTriggerAIProccess_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteSummaryByExternalID provides a mock function with given fields: ctx, externalID
func (_m *SummaryUseCase) DeleteSummaryByExternalID(ctx context.Context, externalID uuid.UUID) error {
	ret := _m.Called(ctx, externalID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSummaryByExternalID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, externalID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SummaryUseCase_DeleteSummaryByExternalID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteSummaryByExternalID'
type SummaryUseCase_DeleteSummaryByExternalID_Call struct {
	*mock.Call
}

// DeleteSummaryByExternalID is a helper method to define mock.On call
//   - ctx context.Context
//   - externalID uuid.UUID
func (_e *SummaryUseCase_Expecter) DeleteSummaryByExternalID(ctx interface{}, externalID interface{}) *SummaryUseCase_DeleteSummaryByExternalID_Call {
	return &SummaryUseCase_DeleteSummaryByExternalID_Call{Call: _e.mock.On("DeleteSummaryByExternalID", ctx, externalID)}
}

func (_c *SummaryUseCase_DeleteSummaryByExternalID_Call) Run(run func(ctx context.Context, externalID uuid.UUID)) *SummaryUseCase_DeleteSummaryByExternalID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *SummaryUseCase_DeleteSummaryByExternalID_Call) Return(_a0 error) *SummaryUseCase_DeleteSummaryByExternalID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SummaryUseCase_DeleteSummaryByExternalID_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *SummaryUseCase_DeleteSummaryByExternalID_Call {
	_c.Call.Return(run)
	return _c
}

// GetSummaryByExternalID provides a mock function with given fields: ctx, externalID
func (_m *SummaryUseCase) GetSummaryByExternalID(ctx context.Context, externalID uuid.UUID) (*service.SummaryDetailedOutput, error) {
	ret := _m.Called(ctx, externalID)

	if len(ret) == 0 {
		panic("no return value specified for GetSummaryByExternalID")
	}

	var r0 *service.SummaryDetailedOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*service.SummaryDetailedOutput, error)); ok {
		return rf(ctx, externalID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *service.SummaryDetailedOutput); ok {
		r0 = rf(ctx, externalID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.SummaryDetailedOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, externalID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SummaryUseCase_GetSummaryByExternalID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSummaryByExternalID'
type SummaryUseCase_GetSummaryByExternalID_Call struct {
	*mock.Call
}

// GetSummaryByExternalID is a helper method to define mock.On call
//   - ctx context.Context
//   - externalID uuid.UUID
func (_e *SummaryUseCase_Expecter) GetSummaryByExternalID(ctx interface{}, externalID interface{}) *SummaryUseCase_GetSummaryByExternalID_Call {
	return &SummaryUseCase_GetSummaryByExternalID_Call{Call: _e.mock.On("GetSummaryByExternalID", ctx, externalID)}
}

func (_c *SummaryUseCase_GetSummaryByExternalID_Call) Run(run func(ctx context.Context, externalID uuid.UUID)) *SummaryUseCase_GetSummaryByExternalID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *SummaryUseCase_GetSummaryByExternalID_Call) Return(_a0 *service.SummaryDetailedOutput, _a1 error) *SummaryUseCase_GetSummaryByExternalID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SummaryUseCase_GetSummaryByExternalID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*service.SummaryDetailedOutput, error)) *SummaryUseCase_GetSummaryByExternalID_Call {
	_c.Call.Return(run)
	return _c
}

// ListSummaries provides a mock function with given fields: ctx
func (_m *SummaryUseCase) ListSummaries(ctx context.Context) (*service.SummaryListOutput, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListSummaries")
	}

	var r0 *service.SummaryListOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*service.SummaryListOutput, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *service.SummaryListOutput); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.SummaryListOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SummaryUseCase_ListSummaries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListSummaries'
type SummaryUseCase_ListSummaries_Call struct {
	*mock.Call
}

// ListSummaries is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SummaryUseCase_Expecter) ListSummaries(ctx interface{}) *SummaryUseCase_ListSummaries_Call {
	return &SummaryUseCase_ListSummaries_Call{Call: _e.mock.On("ListSummaries", ctx)}
}

func (_c *SummaryUseCase_ListSummaries_Call) Run(run func(ctx context.Context)) *SummaryUseCase_ListSummaries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SummaryUseCase_ListSummaries_Call) Return(_a0 *service.SummaryListOutput, _a1 error) *SummaryUseCase_ListSummaries_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SummaryUseCase_ListSummaries_Call) RunAndReturn(run func(context.Context) (*service.SummaryListOutput, error)) *SummaryUseCase_ListSummaries_Call {
	_c.Call.Return(run)
	return _c
}

// NewSummaryUseCase creates a new instance of SummaryUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSummaryUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *SummaryUseCase {
	mock := &SummaryUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
