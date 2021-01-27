// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/spf13/afero"
	"github.com/zclconf/go-cty/cty"
)

// Parser interface
//go:generate mockery -case=underscore -inpkg -name=Parser
type Parser interface {
	Parse(filename string, src []byte) (*ConfigSchema, hcl.Diagnostics)
	LoadConfigFile(path string) (*ConfigSchema, hcl.Diagnostics)
	LoadConfigDir(path string) (*Config, hcl.Diagnostics)
	ConfigDirFiles(dir string) (primary []string, diags hcl.Diagnostics)
	IsConfigDir(path string) bool
	Context() *hcl.EvalContext
}

type parser struct {
	fs  afero.Afero
	p   *hclparse.Parser
	ctx *hcl.EvalContext
}

// NewParser creates and returns a new Parser that reads files from the given
// filesystem. If a nil filesystem is passed then the system's "real" filesystem
// will be used, via afero.OsFs.
func NewParser(fs afero.Fs) *parser {
	if fs == nil {
		fs = afero.OsFs{}
	}

	return &parser{
		fs:  afero.Afero{Fs: fs},
		p:   hclparse.NewParser(),
		ctx: createEvalContext(),
	}
}

// Context return EvalContext
func (p *parser) Context() *hcl.EvalContext {
	return p.ctx
}

func (p *parser) Parse(filename string, src []byte) (*ConfigSchema, hcl.Diagnostics) {
	f, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})

	if diags.HasErrors() {
		return nil, diags
	}

	cfg := &ConfigSchema{}

	diags = gohcl.DecodeBody(f.Body, p.ctx, cfg)

	p.parseConfigVars(cfg)

	for _, p := range cfg.Providers {
		if p.Key == "" {
			p.Key = "default"
		}
	}

	return cfg, diags
}

func (p *parser) parseConfigVars(cfg *ConfigSchema) {
	vars := p.ctx.Variables["var"].AsValueMap()
	if vars == nil {
		vars = map[string]cty.Value{}
	}

	for _, v := range cfg.Variables {
		vars[v.Name] = v.Value
	}

	p.ctx.Variables["var"] = cty.ObjectVal(vars)
}
