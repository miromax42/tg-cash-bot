package cache

import (
	"context"

	"github.com/cockroachdb/errors"
)

var ErrMiss = errors.New("no value")

type Cache interface {
	Get(ctx context.Context, key Token, value interface{}) error
	Set(ctx context.Context, key Token, value interface{}) error
	Del(ctx context.Context, key Token) error
}
