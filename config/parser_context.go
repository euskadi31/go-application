// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func createEvalContext() *hcl.EvalContext {
	variables := map[string]cty.Value{
		"var": cty.ObjectVal(map[string]cty.Value{}),
		"env": cty.ObjectVal(map[string]cty.Value{}),
	}

	functions := map[string]function.Function{
		"abs":          stdlib.AbsoluteFunc,
		"ceil":         stdlib.CeilFunc,
		"chomp":        stdlib.ChompFunc,
		"chunklist":    stdlib.ChunklistFunc,
		"coalescelist": stdlib.CoalesceListFunc,
		"compact":      stdlib.CompactFunc,
		"concat":       stdlib.ConcatFunc,
		"contains":     stdlib.ContainsFunc,
		"csvdecode":    stdlib.CSVDecodeFunc,
		"distinct":     stdlib.DistinctFunc,
		"element":      stdlib.ElementFunc,
		"flatten":      stdlib.FlattenFunc,
		"floor":        stdlib.FloorFunc,
		"format":       stdlib.FormatFunc,
		"formatdate":   stdlib.FormatDateFunc,
		"formatlist":   stdlib.FormatListFunc,
		"indent":       stdlib.IndentFunc,
		"join":         stdlib.JoinFunc,
		"jsondecode":   stdlib.JSONDecodeFunc,
		"jsonencode":   stdlib.JSONEncodeFunc,
		"keys":         stdlib.KeysFunc,
		"lower":        stdlib.LowerFunc,
		"max":          stdlib.MaxFunc,
		"merge":        stdlib.MergeFunc,
		"min":          stdlib.MinFunc,
		"parseint":     stdlib.ParseIntFunc,
		"pow":          stdlib.PowFunc,
		"range":        stdlib.RangeFunc,
		"regex":        stdlib.RegexFunc,
		"regexall":     stdlib.RegexAllFunc,
		"signum":       stdlib.SignumFunc,
		"slice":        stdlib.SliceFunc,
		"sort":         stdlib.SortFunc,
		"split":        stdlib.SplitFunc,
		"strrev":       stdlib.ReverseFunc,
		"substr":       stdlib.SubstrFunc,
		"trim":         stdlib.TrimFunc,
		"trimprefix":   stdlib.TrimPrefixFunc,
		"trimspace":    stdlib.TrimSpaceFunc,
		"trimsuffix":   stdlib.TrimSuffixFunc,
		"upper":        stdlib.UpperFunc,
		"values":       stdlib.ValuesFunc,
	}

	return &hcl.EvalContext{
		Variables: variables,
		Functions: functions,
	}
}
