//go:build tools
// +build tools

package tools

import (
	_ "ariga.io/atlas/cmd/atlas"
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/kisielk/godepgraph"
	_ "github.com/rakyll/statik"
	_ "github.com/vektra/mockery/v2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
