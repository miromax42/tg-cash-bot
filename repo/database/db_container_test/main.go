package db_container_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	instance testcontainers.Container
}

func NewTestDatabase() (*TestDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &TestDatabase{
		instance: postgres,
	}, nil
}

func (db *TestDatabase) Port() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal("cannot get port:", err)
	}

	return p.Int()
}

func (db *TestDatabase) ConnectionString() string {
	return fmt.Sprintf("postgresql://postgres:postgres@127.0.0.1:%d/postgres?sslmode=disable", db.Port())
}

func (db *TestDatabase) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	db.instance.Terminate(ctx)
}
