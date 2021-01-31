// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

// ConfigSchema struct
type ConfigSchema struct {
	Variables   []*VariableSchema  `hcl:"var,block"`
	Application *ApplicationSchema `hcl:"app,block"`
	Providers   []*ProviderSchema  `hcl:"provider,block"`
}

// VariableSchema struct
type VariableSchema struct {
	Name  string    `hcl:"name,label"`
	Value cty.Value `hcl:"value"`
}

func (p VariableSchema) key() string {
	return p.Name
}

// ApplicationSchema struct
type ApplicationSchema struct {
}

func (left *ApplicationSchema) merge(right *ApplicationSchema) {

}

// ProviderSchema struct
type ProviderSchema struct {
	Type string   `hcl:"type,label"`
	Key  string   `hcl:"key,optional"`
	HCL  hcl.Body `hcl:",remain"`
}

func (p ProviderSchema) key() string {
	return p.Type + "." + p.Key
}

// Parse ProviderSchema into sub type
func (p ProviderSchema) Parse(ctx *hcl.EvalContext, v interface{}) hcl.Diagnostics {
	return gohcl.DecodeBody(p.HCL, ctx, v)
}
