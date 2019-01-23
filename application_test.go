// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	app := New()

	go func() {
		assert.NoError(t, app.Run())
	}()

	time.Sleep(200 * time.Millisecond)

	assert.NoError(t, app.Close())
}
