// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// LoadConfigDir reads the .tales files in the given directory
// as config files (using LoadConfigFile) and then combines these files into
// a single Module.
//
// If this method returns nil, that indicates that the given directory does not
// exist at all or could not be opened for some reason. Callers may wish to
// detect this case and ignore the returned diagnostics so that they can
// produce a more context-aware error message in that case.
//
// If this method returns a non-nil module while error diagnostics are returned
// then the module may be incomplete but can be used carefully for static
// analysis.
//
// This file does not consider a directory with no files to be an error, and
// will simply return an empty module in that case. Callers should first call
// Parser.IsConfigDir if they wish to recognize that situation.
//
// .hcl files are parsed using the HCL native syntax
func (p *parser) LoadConfigDir(path string) (*Config, hcl.Diagnostics) {
	paths, diags := p.dirFiles(path)
	if diags.HasErrors() {
		return nil, diags
	}

	files, fDiags := p.loadFiles(paths)
	diags = append(diags, fDiags...)

	cfg, cfgDiags := NewConfig(files)
	diags = append(diags, cfgDiags...)

	cfg.SourceDir = path

	return cfg, diags
}

// ConfigDirFiles returns lists of the files configuration
// files in the given directory.
//
// If the given directory does not exist or cannot be read, error diagnostics
// are returned. If errors are returned, the resulting lists may be incomplete.
func (p parser) ConfigDirFiles(dir string) (primary []string, diags hcl.Diagnostics) {
	return p.dirFiles(dir)
}

// IsConfigDir determines whether the given path refers to a directory that
// exists and contains at least one config file (with a .hcl extension.)
func (p *parser) IsConfigDir(path string) bool {
	paths, _ := p.dirFiles(path)
	return len(paths) > 0
}

func (p *parser) loadFiles(paths []string) ([]*ConfigSchema, hcl.Diagnostics) {
	var files []*ConfigSchema
	var diags hcl.Diagnostics

	for _, path := range paths {
		f, d := p.LoadConfigFile(path)

		diags = append(diags, d...)
		if f != nil {
			files = append(files, f)
		}
	}

	return files, diags
}

func (p *parser) dirFiles(dir string) (files []string, diags hcl.Diagnostics) {
	infos, err := p.fs.ReadDir(dir)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to read config directory",
			Detail:   fmt.Sprintf("Config directory %s does not exist or cannot be read.", dir),
		})
		return
	}

	for _, info := range infos {
		if info.IsDir() {
			// We only care about files
			continue
		}

		name := info.Name()
		ext := fileExt(name)
		if ext == "" || IsIgnoredFile(name) {
			continue
		}

		fullPath := filepath.Join(dir, name)
		files = append(files, fullPath)
	}

	return
}

// fileExt returns the Terraform configuration extension of the given
// path, or a blank string if it is not a recognized extension.
func fileExt(path string) string {
	if strings.HasSuffix(path, ".hcl") {
		return ".hcl"
	}

	return ""
}

// IsIgnoredFile returns true if the given filename (which must not have a
// directory path ahead of it) should be ignored as e.g. an editor swap file.
func IsIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // Unix-like hidden files
		strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}

// IsEmptyDir returns true if the given filesystem path contains no Tales
// configuration files.
//
// Unlike the methods of the Parser type, this function always consults the
// real filesystem, and thus it isn't appropriate to use when working with
// configuration loaded from a plan file.
func IsEmptyDir(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return true, nil
	}

	p := NewParser(nil)
	files, diags := p.dirFiles(path)
	if diags.HasErrors() {
		return false, diags
	}

	return len(files) == 0, nil
}
