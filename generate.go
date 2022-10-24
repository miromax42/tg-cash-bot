package main

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema --feature sql/upsert --feature sql/versioned-migration
//go:generate go run -mod=mod github.com/vektra/mockery/v2 --all --keeptree --dir ./repo
//go:generate go run -mod=mod github.com/vektra/mockery/v2 --all --keeptree --dir ./currency
//go:generate go run -mod=mod github.com/vektra/mockery/v2 --all --keeptree --dir ./telegram/tools
//go:generate go run -mod=mod github.com/vektra/mockery/v2 --all --keeptree --dir ./util/logger
