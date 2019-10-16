// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

// Configuration interface
//go:generate mockery -case=underscore -inpkg -name=Configuration
type Configuration interface {
	SetDefault(key string, value interface{})
}
