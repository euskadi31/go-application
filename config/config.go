// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Config struct
type Config struct {
	SourceDir string
	App       *ApplicationSchema
	Variables map[string]cty.Value
	Providers map[string]*ProviderSchema
}

// NewConfig constructor
func NewConfig(files []*ConfigSchema) (*Config, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	cfg := &Config{
		App:       &ApplicationSchema{},
		Variables: map[string]cty.Value{},
		Providers: map[string]*ProviderSchema{},
	}

	for _, file := range files {
		fileDiags := cfg.appendFile(file)
		diags = append(diags, fileDiags...)
	}

	return cfg, diags
}

func (c *Config) appendFile(file *ConfigSchema) hcl.Diagnostics {
	var diags hcl.Diagnostics

	c.App.merge(file.Application)

	for _, v := range file.Variables {
		key := v.key()

		if _, exists := c.Variables[key]; exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate variable %q configuration", key),
				Detail:   fmt.Sprintf("A var named %q was already declared. Variable names must be unique in config.", key),
				//Subject:  &r.DeclRange,
			})
			continue
		}

		c.Variables[key] = v.Value
	}

	for _, p := range file.Providers {
		key := p.key()

		if existing, exists := c.Providers[key]; exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate provider %q configuration", existing.Type),
				Detail:   fmt.Sprintf("A provider %q named %q was already declared. Provider names must be unique per type in config.", existing.Type, existing.Key),
				//Subject:  &r.DeclRange,
			})
			continue
		}

		c.Providers[key] = p
	}

	return diags
}
