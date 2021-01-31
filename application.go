// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/hcl/v2"

	"github.com/euskadi31/go-application/config"
	"github.com/euskadi31/go-service"
	"github.com/rs/zerolog/log"
)

// Application interface
type Application interface {
	Register(provider ServiceProvider)
	AddConfigPath(path string)
	Run() error
	Close() error
}

var _ Application = (*App)(nil)

// App struct
type App struct {
	parser            config.Parser
	name              string
	configSearchPaths []string
	cfg               *config.Config
	signal            chan os.Signal
	container         service.Container
	providers         []*providerRegister
}

// New Application
func New(name string) Application {
	app := &App{
		name: name,
		configSearchPaths: []string{
			"./",
			"~/." + name + "/",
			"/etc/" + name + "/",
		},
		signal:    make(chan os.Signal, 1),
		container: service.New(),
		parser:    config.NewParser(nil),
	}

	return app
}

// AddConfigPath to search paths
func (a *App) AddConfigPath(path string) {
	a.configSearchPaths = append(a.configSearchPaths, path)
}

func (a *App) loadConfig() hcl.Diagnostics {

	for _, p := range a.configSearchPaths {
		if a.parser.IsConfigDir(p) {
			cfg, diags := a.parser.LoadConfigDir(p)
			if diags.HasErrors() {
				return diags
			}

			a.cfg = cfg

			return nil
		}
	}

	a.cfg, _ = config.NewConfig(nil)

	return nil
}

// Register ServiceProvider
func (a *App) Register(provider ServiceProvider) {
	//provider.Register(a.container)

	a.providers = append(a.providers, &providerRegister{
		ServiceProvider: provider,
	})
}

// Run Application
func (a *App) Run() (err error) {
	if err := a.loadConfig(); err != nil {
		return err
	}

	signal.Notify(a.signal, os.Interrupt, syscall.SIGTERM)

	bootables := []BootableProvider{}
	configurables := map[string]ConfigurableProvider{}

	for _, provider := range a.providers {
		provider.Register(a.container)

		if bootable, ok := provider.ServiceProvider.(BootableProvider); ok {
			bootables = append(bootables, bootable)
		}

		if configurable, ok := provider.ServiceProvider.(ConfigurableProvider); ok {
			configurables[provider.Name()] = configurable
		}
	}

	if diags := a.configure(configurables); diags.HasErrors() {
		return diags
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

// Close Application
func (a *App) Close() error {
	a.signal <- syscall.SIGTERM

	return nil
}

func (a App) configure(configurables map[string]ConfigurableProvider) hcl.Diagnostics {
	log.Info().Msg("Configuring...")

	for _, p := range a.cfg.Providers {
		if configurable, ok := configurables[p.Type]; ok {
			diags := configurable.Config(a.parser.Context(), p)
			if diags.HasErrors() {
				return diags
			}
		} else {
			return hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Provider not registered",
					Detail:   fmt.Sprintf("The provider %q is not registered in app.", p.Type),
				},
			}
		}
	}

	return hcl.Diagnostics{}
}
