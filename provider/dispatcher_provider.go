// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package provider

import (
	"github.com/euskadi31/go-eventemitter"
	"github.com/euskadi31/go-service"
)

// Dispatcher Services keys
const (
	EventDispatcherKey = "event.dispatcher"
)

// EventDispatcherServiceProvider struct
type EventDispatcherServiceProvider struct {
}

// NewEventDispatcherServiceProvider constructor
func NewEventDispatcherServiceProvider() *EventDispatcherServiceProvider {
	return &EventDispatcherServiceProvider{}
}

// Register implements application.ServiceProvider
func (p EventDispatcherServiceProvider) Register(app service.Container) {
	app.Set(EventDispatcherKey, func(c service.Container) interface{} {
		return eventemitter.New() // eventemitter.EventEmitter
	})
}

// Priority implements application.BootableProvider
func (p EventDispatcherServiceProvider) Priority() int {
	return 0
}

// Start implements application.BootableProvider
func (p EventDispatcherServiceProvider) Start(app service.Container) error {
	return nil
}

// Stop implements application.BootableProvider
func (p EventDispatcherServiceProvider) Stop(app service.Container) error {
	dispatcher := app.Get(EventDispatcherKey).(eventemitter.EventEmitter)

	dispatcher.Wait()

	return nil
}
