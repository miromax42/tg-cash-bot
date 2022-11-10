package database

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/cockroachdb/errors"
	"go.opentelemetry.io/otel"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type PersonalSettings struct {
	db *ent.Client
}

func NewPersonalSettings(client *ent.Client) *PersonalSettings {
	return &PersonalSettings{
		db: client,
	}
}

func (p *PersonalSettings) Get(ctx context.Context, id int64) (*repo.PersonalSettingsResp, error) {
	ctx, span := otel.Tracer(util.RequestTrace).Start(ctx, "Settings.Get")
	defer span.End()

	settings, err := p.db.PersonalSettings.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return repo.DefaultPersonalSettingsResp(), nil
		}

		return nil, errors.Wrap(err, "get")
	}

	return &repo.PersonalSettingsResp{
		Currency: settings.Currency,
		Limit:    settings.Limit,
	}, nil
}

func (p *PersonalSettings) Set(ctx context.Context, req repo.PersonalSettingsReq) error {
	ctx, span := otel.Tracer(util.RequestTrace).Start(ctx, "Settings.Set")
	defer span.End()

	return WithTx(ctx, p.db, func(tx *ent.Tx) error {
		expenses := NewExpense(p.db)
		sum, err := expenses.allUserExpense(ctx, repo.ListUserExpenseReq{
			UserID:   req.UserID,
			FromTime: util.TimeMonthAgo(),
			ToTime:   time.Now(),
		})
		if err != nil {
			return errors.Wrap(err, "get sum expenses")
		}
		if req.Limit != nil && sum > *req.Limit {
			return errors.WithHint(repo.ErrLimitExceed, "current sum is bigger than chosen limit")
		}

		return errors.Wrap(p.db.PersonalSettings.Create().
			SetID(req.UserID).
			SetNillableCurrency(req.Currency).
			SetNillableLimit(req.Limit).
			OnConflict(
				sql.ConflictColumns("id"),
			).
			UpdateNewValues().
			Exec(ctx), "upsert settings")
	})
}
