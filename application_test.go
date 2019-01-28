// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"errors"
	"testing"
	"time"

	"github.com/euskadi31/go-service"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct{}

func (mockProvider) Register(app service.Container) {

}

func (mockProvider) Priority() int {
	return 1
}

func (mockProvider) Start(app service.Container) error {
	return errors.New("fail")
}

func (mockProvider) Stop(app service.Container) error {
	return errors.New("fail")
}

func TestApplication(t *testing.T) {
	app := New()

	app.Register(&mockProvider{})

	go func() {
		assert.NoError(t, app.Run())
	}()

	time.Sleep(200 * time.Millisecond)

	assert.NoError(t, app.Close())
}
