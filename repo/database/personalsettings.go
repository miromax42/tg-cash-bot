package database

import (
	"context"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
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
	settings, err := p.db.PersonalSettings.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return repo.DefaultPersonalSettingsReq(), nil
		}
		return nil, err
	}

	return &repo.PersonalSettingsResp{
		Currency: settings.Currency,
	}, nil
}

func (p *PersonalSettings) Set(ctx context.Context, req repo.PersonalSettingsReq) error {
	if err := p.db.PersonalSettings.
		Create().
		SetID(req.UserID).
		OnConflict().
		UpdateNewValues().
		Exec(ctx); err != nil {
		return err
	}

	return nil
}
