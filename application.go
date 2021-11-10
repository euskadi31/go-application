// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/v2"

	"github.com/euskadi31/go-application/config"
	"github.com/euskadi31/go-application/provider"
	"github.com/euskadi31/go-service"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
	options           *viper.Viper
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
		options: viper.New(),
		name:    name,
		configSearchPaths: []string{
			"./",
			"~/." + name + "/",
			"/etc/" + name + "/",
		},
		signal:    make(chan os.Signal, 1),
		container: service.New(),
	}

	app.Register(provider.NewLoggerServiceProvider())

	return app
}

// AddConfigPath to search paths
func (a *App) AddConfigPath(path string) {
	a.configSearchPaths = append(a.configSearchPaths, path)
}

func (a *App) loadConfig() hcl.Diagnostics {
	a.options.SetConfigName("config")

	for _, p := range a.configSearchPaths {
		a.options.AddConfigPath(p)
	}

	a.options.SetEnvPrefix(strings.Replace(a.name, "-", "_", -1))
	a.options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	a.options.AutomaticEnv()

	if err := a.options.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	schema := &config.ConfigSchema{}

	if err := a.options.Unmarshal(schema); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	cfg, diags := config.NewConfig([]*config.ConfigSchema{schema})
	if diags.HasErrors() {
		return diags
	}

	defer func() {
		a.cfg.App.Name = a.name
	}()

	a.cfg = cfg

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

	a.container.SetValue("config", a.cfg)

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

	if err := a.configure(configurables); err != nil {
		return err
	}

	// load logger provider
	a.container.Get(provider.LoggerKey)

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

func (a *App) configure(configurables map[string]ConfigurableProvider) error {
	for k, configurable := range configurables {
		sc := a.options.Sub("provider")

		b, _ := json.Marshal(sc.AllSettings())

		c := string(b)

		spew.Dump(c)
		spew.Dump(k)

		//sc := a.options.Sub("provider." + k)
		if sc == nil {
			sc = viper.New()
		}

		if err := configurable.Config(a.container, sc); err != nil {
			return err
		}
	}

	return nil
}
