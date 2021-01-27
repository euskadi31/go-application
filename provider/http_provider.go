// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package provider

import (
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-service"
)

// @TODO: extract this provider to github.com/hyperscale-stack/router-service-provider

// HTTP Services keys
const (
	HTTPServerConfigKey = "http.server.config"
	HTTPServerKey       = "http.server"
)

// HTTPServiceProvider struct
type HTTPServiceProvider struct {
}

// NewHTTPServiceProvider constructor
func NewHTTPServiceProvider() *HTTPServiceProvider {
	return &HTTPServiceProvider{}
}

// Register implements application.ServiceProvider
func (p HTTPServiceProvider) Register(app service.Container) {
	app.Set(HTTPServerConfigKey, func(c service.Container) interface{} {
		return &server.Configuration{}
	})

	app.Set(HTTPServerKey, func(c service.Container) interface{} {
		cfg := c.Get(HTTPServerConfigKey).(*server.Configuration)

		router := server.New(cfg)

		return router // *server.Server
	})
}

// Priority implements application.BootableProvider
func (p HTTPServiceProvider) Priority() int {
	return 255
}

// Start implements application.BootableProvider
func (p HTTPServiceProvider) Start(app service.Container) error {
	serv := app.Get(HTTPServerKey).(*server.Server)

	return serv.Run()
}

// Stop implements application.BootableProvider
func (p HTTPServiceProvider) Stop(app service.Container) error {
	serv := app.Get(HTTPServerKey).(*server.Server)

	return serv.Shutdown()
}
