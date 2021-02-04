// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	p := NewParser(nil)

	cfg, diags := p.LoadConfigDir("./testdata")
	assert.Equal(t, "./testdata", cfg.SourceDir)
	if diags.HasErrors() {
		t.Error(diags.Error())
	}

	assert.False(t, diags.HasErrors())
	assert.Equal(t, 2, len(cfg.Providers))
	assert.Equal(t, 1, len(cfg.Variables))

}
