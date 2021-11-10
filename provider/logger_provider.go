// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package provider

import (
	stdlog "log"
	"os"

	"github.com/euskadi31/go-application/config"
	"github.com/euskadi31/go-service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Logger Services keys
const (
	LoggerKey       = "logger"
	LoggerConfigKey = "logger.config"
)

// LoggerServiceProviderConfig struct
type LoggerServiceProviderConfig struct {
	Level   string                `hcl:"level"`
	Writers []*LoggerWriterSchema `hcl:"writer,block"`
}

// LoggerWriterSchema struct
type LoggerWriterSchema struct {
	Name  string       `hcl:"name,label"`
	Where *WhereSchema `hcl:"where,block"`
}

// WhereSchema struct
type WhereSchema struct {
	Env []string `hcl:"env,optional"`
}

// LoggerServiceProvider struct
type LoggerServiceProvider struct {
}

// NewLoggerServiceProvider constructor
func NewLoggerServiceProvider() *LoggerServiceProvider {
	return &LoggerServiceProvider{}
}

// Name implements application.ServiceProvider
func (p LoggerServiceProvider) Name() string {
	return "logger"
}

// Register implements application.ServiceProvider
func (p LoggerServiceProvider) Register(app service.Container) {
	app.Set(LoggerKey, func(c service.Container) interface{} {
		appCfg := c.Get("config").(*config.Config)
		loggerCfg := c.Get(LoggerConfigKey).(*LoggerServiceProviderConfig)

		zlevel, err := zerolog.ParseLevel(loggerCfg.Level)
		if err != nil {

		}

		zerolog.SetGlobalLevel(zlevel)

		logger := zerolog.New(os.Stdout).With().
			Timestamp().
			Str("role", appCfg.App.Name).
			//Str("version", version.Get().Version).
			Str("env", appCfg.App.Environment).
			//Caller().
			Logger()

		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal().Err(err).Msg("Stdin.Stat failed")
		}

		if (fi.Mode() & os.ModeCharDevice) != 0 {
			logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		stdlog.SetFlags(0)
		stdlog.SetOutput(logger)

		log.Logger = logger

		return logger
	})
}

// Config provider
func (p LoggerServiceProvider) Config(app service.Container, options *viper.Viper) error {
	cfg := &LoggerServiceProviderConfig{}

	options.SetDefault("level", "info")

	if err := options.Unmarshal(cfg); err != nil {
		return err
	}

	app.SetValue(LoggerConfigKey, cfg)

	return nil
}
