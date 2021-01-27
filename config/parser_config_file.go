// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
)

// LoadConfigFile config
func (p *parser) LoadConfigFile(path string) (*ConfigSchema, hcl.Diagnostics) {
	return p.loadConfigFile(path)
}

func (p *parser) loadConfigFile(path string) (*ConfigSchema, hcl.Diagnostics) {
	src, err := p.fs.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Configuration file not found",
					Detail:   fmt.Sprintf("The configuration file %s does not exist.", path),
				},
			}
		}

		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read configuration",
				Detail:   fmt.Sprintf("Can't read %s: %s.", path, err),
			},
		}
	}

	return p.Parse(path, src)
}
