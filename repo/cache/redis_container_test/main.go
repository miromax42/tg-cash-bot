package redis_container

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestRedis struct {
	instance testcontainers.Container
}

func NewTestRedis(t *testing.T) *TestRedis {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "redis:6-alpine",
		ExposedPorts: []string{"6379/tcp"},
		AutoRemove:   true,
		Env:          map[string]string{},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	redis, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	return &TestRedis{
		instance: redis,
	}
}

func (db *TestRedis) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "6379")
	require.NoError(t, err)

	return p.Int()
}

func (db *TestRedis) ConnectionSocketAddress(t *testing.T) string {
	return fmt.Sprintf("127.0.0.1:%d", db.Port(t))
}

func GetRedisConnectionString(t *testing.T) string {
	if _, ok := os.LookupEnv("CI_PROJECT_ID"); ok {
		connStr, ok := os.LookupEnv("REDIS_SOCKET_ADDRESS")
		require.True(t, ok)

		return connStr
	}

	c := NewTestRedis(t)

	return c.ConnectionSocketAddress(t)
}
