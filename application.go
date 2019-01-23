// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/euskadi31/go-application/provider"
	"github.com/euskadi31/go-service"
	"github.com/rs/zerolog/log"
)

// Application interface
type Application interface {
	Register(provider ServiceProvider)
	Run() error
	Close() error
}

var _ Application = (*App)(nil)

// App struct
type App struct {
	container service.Container
	providers []ServiceProvider
}

// New Application
func New() Application {
	app := &App{
		container: service.New(),
	}

	app.Register(provider.NewEventDispatcherServiceProvider())
	app.Register(provider.NewHTTPServiceProvider())

	return app
}

// Register ServiceProvider
func (a *App) Register(provider ServiceProvider) {
	provider.Register(a.container)

	a.providers = append(a.providers, provider)
}

// Run Application
func (a *App) Run() (err error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	By(func(left, right ServiceProvider) bool {
		return left.Priority() < right.Priority()
	}).Sort(a.providers)

	log.Info().Msg("Starting...")

	for _, provider := range a.providers {
		provider.Register(a.container)

		if bootable, ok := provider.(BootableProvider); ok {
			go func() {
				if err := bootable.Start(a.container); err != nil {
					log.Error().Err(err).Msg("Start failed")
				}
			}()
		}
	}

	log.Info().Msg("Started")

	<-sig

	log.Info().Msg("Shutdown")

	for _, provider := range a.providers {
		if bootable, ok := provider.(BootableProvider); ok {
			if err := bootable.Stop(a.container); err != nil {
				log.Error().Err(err).Msg("Stop failed")
			}
		}
	}

	return nil
}

// Close Application
func (a *App) Close() error {
	return nil
}
