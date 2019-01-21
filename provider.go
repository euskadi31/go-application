// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"sort"

	"github.com/euskadi31/go-eventemitter"
	"github.com/euskadi31/go-service"
)

// ServiceProvider interface
type ServiceProvider interface {
	Register(app service.Container)
	Priority() int
}

// BootableProvider interface
type BootableProvider interface {
	Start(app service.Container) error
	Stop(app service.Container) error
}

// EventListenerProvider interface
type EventListenerProvider interface {
	Subscribe(app service.Container, dispatcher eventemitter.EventEmitter)
}

type providerSorter struct {
	providers []ServiceProvider
	by        func(left, right ServiceProvider) bool
}

// Len is part of sort.Interface.
func (s *providerSorter) Len() int {
	return len(s.providers)
}

// Swap is part of sort.Interface.
func (s *providerSorter) Swap(i, j int) {
	s.providers[i], s.providers[j] = s.providers[j], s.providers[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *providerSorter) Less(i, j int) bool {
	return s.by(s.providers[i], s.providers[j])
}

// By sorter
type By func(left, right ServiceProvider) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(providers []ServiceProvider) {
	ps := &providerSorter{
		providers: providers,
		by:        by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}

	sort.Sort(ps)
}
