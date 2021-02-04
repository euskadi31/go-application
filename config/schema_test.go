// Copyright 2021 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/stretchr/testify/assert"
)

type ProviderFoo struct {
	Foo string `hcl:"foo"`
}

func TestProviderSchemaParse(t *testing.T) {
	content := []byte(`provider "mysql" {
		key = "hi"
		foo = "bar"
	}`)

	s := &ConfigSchema{}

	err := hclsimple.Decode("config.hcl", content, nil, s)
	assert.NoError(t, err)

	p := &ProviderFoo{}

	diags := s.Providers[0].Parse(nil, p)
	assert.False(t, diags.HasErrors())
}
