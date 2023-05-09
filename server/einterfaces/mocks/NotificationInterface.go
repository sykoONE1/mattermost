// Code generated by mockery v2.23.2. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import (
	model "github.com/mattermost/mattermost-server/server/public/model"
	mock "github.com/stretchr/testify/mock"
)

// NotificationInterface is an autogenerated mock type for the NotificationInterface type
type NotificationInterface struct {
	mock.Mock
}

// CheckLicense provides a mock function with given fields:
func (_m *NotificationInterface) CheckLicense() *model.AppError {
	ret := _m.Called()

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func() *model.AppError); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// GetNotificationMessage provides a mock function with given fields: ack, userID
func (_m *NotificationInterface) GetNotificationMessage(ack *model.PushNotificationAck, userID string) (*model.PushNotification, *model.AppError) {
	ret := _m.Called(ack, userID)

	var r0 *model.PushNotification
	var r1 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.PushNotificationAck, string) (*model.PushNotification, *model.AppError)); ok {
		return rf(ack, userID)
	}
	if rf, ok := ret.Get(0).(func(*model.PushNotificationAck, string) *model.PushNotification); ok {
		r0 = rf(ack, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PushNotification)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.PushNotificationAck, string) *model.AppError); ok {
		r1 = rf(ack, userID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewNotificationInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewNotificationInterface creates a new instance of NotificationInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNotificationInterface(t mockConstructorTestingTNewNotificationInterface) *NotificationInterface {
	mock := &NotificationInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
