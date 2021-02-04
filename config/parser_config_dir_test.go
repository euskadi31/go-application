// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserIsConfigDir(t *testing.T) {
	p := NewParser(nil)

	assert.True(t, p.IsConfigDir("./testdata"))
}
