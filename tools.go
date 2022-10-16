//go:build tools
// +build tools

package tools

import (
	_ "ariga.io/atlas/cmd/atlas"
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/kisielk/godepgraph"
	_ "github.com/vektra/mockery/v2"
)
