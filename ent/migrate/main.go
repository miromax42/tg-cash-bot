//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent/migrate"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/database/db_container_test"
)

func main() {
	ctx := context.Background()
	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := atlas.NewLocalDir("ent/migrate/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}
	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}

	db, err := db_container_test.NewTestDatabase()
	if err != nil {
		log.Fatal("cannot start container:", err)
	}
	defer db.Close()

	err = migrate.NamedDiff(ctx, db.ConnectionString(), os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}

	fmt.Println("Success!")
}
