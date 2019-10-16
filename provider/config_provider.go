// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package provider

import (
	"github.com/euskadi31/go-service"
	"github.com/spf13/viper"
)

// Config Services keys
const (
	ConfigKey = "config"
)

// ConfigServiceProvider struct
type ConfigServiceProvider struct {
}

// NewConfigServiceProvider constructor
func NewConfigServiceProvider() *ConfigServiceProvider {
	return &ConfigServiceProvider{}
}

// Register implements application.ServiceProvider
func (p ConfigServiceProvider) Register(app service.Container) {
	app.Set(ConfigKey, func(c service.Container) interface{} {
		return viper.New() // *viper.Viper
	})
}

// Priority implements application.BootableProvider
func (p ConfigServiceProvider) Priority() int {
	return 0
}

// Start implements application.BootableProvider
func (p ConfigServiceProvider) Start(app service.Container) error {
	cfg := app.Get(ConfigKey).(*viper.Viper)

	err := cfg.ReadInConfig()
	if _, notFound := err.(viper.ConfigFileNotFoundError); err != nil && !notFound {
		return err
	}

	return nil
}

// Stop implements application.BootableProvider
func (p ConfigServiceProvider) Stop(app service.Container) error {
	return nil
}
