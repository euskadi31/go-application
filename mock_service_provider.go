// Code generated by mockery v1.0.0. DO NOT EDIT.

package application

import mock "github.com/stretchr/testify/mock"
import service "github.com/euskadi31/go-service"

// MockServiceProvider is an autogenerated mock type for the ServiceProvider type
type MockServiceProvider struct {
	mock.Mock
}

// Register provides a mock function with given fields: app
func (_m *MockServiceProvider) Register(app service.Container) {
	_m.Called(app)
}
