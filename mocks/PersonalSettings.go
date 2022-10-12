// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	repo "gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

// PersonalSettings is an autogenerated mock type for the PersonalSettings type
type PersonalSettings struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, userID
func (_m *PersonalSettings) Get(ctx context.Context, userID int64) (*repo.PersonalSettingsResp, error) {
	ret := _m.Called(ctx, userID)

	var r0 *repo.PersonalSettingsResp
	if rf, ok := ret.Get(0).(func(context.Context, int64) *repo.PersonalSettingsResp); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repo.PersonalSettingsResp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: ctx, req
func (_m *PersonalSettings) Set(ctx context.Context, req repo.PersonalSettingsReq) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.PersonalSettingsReq) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPersonalSettings interface {
	mock.TestingT
	Cleanup(func())
}

// NewPersonalSettings creates a new instance of PersonalSettings. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPersonalSettings(t mockConstructorTestingTNewPersonalSettings) *PersonalSettings {
	mock := &PersonalSettings{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}