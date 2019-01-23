// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package provider

import (
	"testing"

	"github.com/euskadi31/go-service"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServiceProvider(t *testing.T) {
	container := service.New()

	p := NewHTTPServiceProvider()

	assert.Equal(t, 0, p.Priority())

	p.Register(container)

	assert.True(t, container.Has(HTTPServerConfigKey))
	assert.True(t, container.Has(HTTPServerKey))

	assert.EqualError(t, p.Start(container), "http or https server is not configured")

	assert.NoError(t, p.Stop(container))
}
