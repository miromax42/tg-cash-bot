//go:build tools
// +build tools

package tools

import (
	_ "ariga.io/atlas/cmd/atlas"
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/vektra/mockery/v2"
)
