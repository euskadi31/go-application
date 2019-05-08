// Code generated by mockery v1.0.0. DO NOT EDIT.

package application

import mock "github.com/stretchr/testify/mock"
import service "github.com/euskadi31/go-service"

// MockBootableProvider is an autogenerated mock type for the BootableProvider type
type MockBootableProvider struct {
	mock.Mock
}

// Priority provides a mock function with given fields:
func (_m *MockBootableProvider) Priority() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Start provides a mock function with given fields: app
func (_m *MockBootableProvider) Start(app service.Container) error {
	ret := _m.Called(app)

	var r0 error
	if rf, ok := ret.Get(0).(func(service.Container) error); ok {
		r0 = rf(app)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields: app
func (_m *MockBootableProvider) Stop(app service.Container) error {
	ret := _m.Called(app)

	var r0 error
	if rf, ok := ret.Get(0).(func(service.Container) error); ok {
		r0 = rf(app)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
