// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/euskadi31/go-service"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	RegisterCalled uint64
	StartCalled    uint64
	StopCalled     uint64
	priority       int
}

func (p *mockProvider) Register(app service.Container) {
	atomic.AddUint64(&p.RegisterCalled, 1)
}

func (p mockProvider) Priority() int {
	return p.priority
}

func (p *mockProvider) Start(app service.Container) error {
	atomic.AddUint64(&p.StartCalled, 1)

	return errors.New("fail")
}

func (p *mockProvider) Stop(app service.Container) error {
	atomic.AddUint64(&p.StopCalled, 1)

	return errors.New("fail")
}

func TestApplication(t *testing.T) {
	providerMock1 := &mockProvider{
		priority: 2,
	}
	providerMock2 := &mockProvider{
		priority: 1,
	}

	app := New()
	app.Register(providerMock1)
	app.Register(providerMock2)

	go func(app Application) {
		assert.NoError(t, app.Run())
	}(app)

	time.Sleep(200 * time.Millisecond)

	assert.NoError(t, app.Close())

	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock1.RegisterCalled))
	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock2.RegisterCalled))
	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock1.StartCalled))
	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock2.StartCalled))
	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock1.StopCalled))
	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock2.StopCalled))
}
