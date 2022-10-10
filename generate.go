package main

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema --feature sql/upsert --feature sql/versioned-migration
