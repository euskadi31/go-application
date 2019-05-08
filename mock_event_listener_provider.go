// Code generated by mockery v1.0.0. DO NOT EDIT.

package application

import eventemitter "github.com/euskadi31/go-eventemitter"
import mock "github.com/stretchr/testify/mock"
import service "github.com/euskadi31/go-service"

// MockEventListenerProvider is an autogenerated mock type for the EventListenerProvider type
type MockEventListenerProvider struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: app, dispatcher
func (_m *MockEventListenerProvider) Subscribe(app service.Container, dispatcher eventemitter.EventEmitter) {
	_m.Called(app, dispatcher)
}
