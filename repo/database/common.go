package database

import (
	"context"

	"github.com/cockroachdb/errors"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
)

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return errors.Wrap(err, "create tx")
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.WithHintf(err, "rolling back tx: %v", rerr)
		}

		return errors.Wrapf(err, "during tx")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "committing tx")
	}

	return nil
}

func firstOrZero[T any](arr []T) T {
	var zero T
	if len(arr) == 0 {
		return zero
	}

	return arr[0]
}
