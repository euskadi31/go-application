// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/euskadi31/go-service"
	"github.com/rs/zerolog/log"
)

// Application interface.
type Application interface {
	Register(provider ServiceProvider)
	Run() error
	Close() error
}

var _ Application = (*App)(nil)

// App struct.
type App struct {
	signal    chan os.Signal
	container service.Container
	providers []*providerRegister
}

// New Application.
func New() Application {
	app := &App{
		signal:    make(chan os.Signal, 1),
		container: service.New(),
	}

	return app
}

// Register ServiceProvider.
func (a *App) Register(provider ServiceProvider) {
	a.providers = append(a.providers, &providerRegister{
		ServiceProvider: provider,
	})
}

// Run Application.
func (a *App) Run() (err error) {
	signal.Notify(a.signal, os.Interrupt, syscall.SIGTERM)

	bootables := []BootableProvider{}

	for _, provider := range a.providers {
		provider.Register(a.container)

		if bootable, ok := provider.ServiceProvider.(BootableProvider); ok {
			bootables = append(bootables, bootable)
		}
	}

	By(func(left, right BootableProvider) bool {
		return left.Priority() < right.Priority()
	}).Sort(bootables)

	log.Info().Msg("Starting...")

	for _, provider := range bootables {
		go func(p BootableProvider) {
			if err := p.Start(a.container); err != nil {
				log.Error().Err(err).Msg("Start failed")
			}
		}(provider)
	}

	<-a.signal

	// Reversing order for closing
	for i := len(bootables)/2 - 1; i >= 0; i-- {
		opp := len(bootables) - 1 - i
		bootables[i], bootables[opp] = bootables[opp], bootables[i]
	}

	log.Info().Msg("Shutdown...")

	for _, provider := range bootables {
		if err := provider.Stop(a.container); err != nil {
			log.Error().Err(err).Msg("Stop failed")
		}
	}

	return nil
}

// Close Application.
func (a *App) Close() error {
	a.signal <- syscall.SIGTERM

	return nil
}
