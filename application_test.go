// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/euskadi31/go-application/config"
	"github.com/euskadi31/go-service"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

type mockProviderConfig struct {
	*config.ProviderSchema
	Foo string `hcl:"foo,optional"`
}

type mockProviderWithConfig struct {
	RegisterCalled uint64
	NameCalled     uint64
	ConfigCalled   uint64
	name           string
	mtx            sync.RWMutex
	cfg            *mockProviderConfig
}

func (p *mockProviderWithConfig) Register(app service.Container) {
	atomic.AddUint64(&p.RegisterCalled, 1)
}

func (p *mockProviderWithConfig) Name() string {
	atomic.AddUint64(&p.NameCalled, 1)

	return p.name
}

func (p *mockProviderWithConfig) GetConfig() *mockProviderConfig {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return p.cfg
}

func (p *mockProviderWithConfig) Config(ctx *hcl.EvalContext, schema *config.ProviderSchema) hcl.Diagnostics {
	atomic.AddUint64(&p.ConfigCalled, 1)

	p.mtx.Lock()
	defer p.mtx.Unlock()

	return schema.Parse(ctx, p.cfg)
}

type mockProvider struct {
	RegisterCalled uint64
	NameCalled     uint64
	StartCalled    uint64
	StopCalled     uint64
	priority       int
	name           string
}

func (p *mockProvider) Register(app service.Container) {
	atomic.AddUint64(&p.RegisterCalled, 1)
}

func (p *mockProvider) Name() string {
	atomic.AddUint64(&p.NameCalled, 1)

	return p.name
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
		name:     "mock-1",
	}
	providerMock2 := &mockProvider{
		priority: 1,
		name:     "mock-2",
	}

	app := New("acme")
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

func TestApplicationWithConfig(t *testing.T) {
	providerMock := &mockProviderWithConfig{
		name: "mock-1",
		cfg:  &mockProviderConfig{},
	}

	app := New("acme")
	app.AddConfigPath("./testdata/mock-1/")

	app.Register(providerMock)

	go func(app Application) {
		assert.NoError(t, app.Run())
	}(app)

	time.Sleep(200 * time.Millisecond)

	assert.NoError(t, app.Close())

	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock.RegisterCalled))
	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock.NameCalled))
	assert.Equal(t, uint64(1), atomic.LoadUint64(&providerMock.ConfigCalled))

	assert.Equal(t, "bar", providerMock.GetConfig().Foo)
}
