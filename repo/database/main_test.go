package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
)

func GetDB(t *testing.T) *ent.Client {
	db, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)

	err = db.Schema.Create(context.Background())
	require.NoError(t, err)

	return db
}

func generator() func() int64 {
	var inc int64
	return func() int64 {
		inc++
		return inc
	}
}
