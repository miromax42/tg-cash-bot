package repo

import (
	"context"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
)

type PersonalSettings interface {
	Get(ctx context.Context, id int64) (*PersonalSettingsResp, error)
	Set(ctx context.Context, req PersonalSettingsReq) error
}

type PersonalSettingsResp struct {
	Currency currency.Token
}

type PersonalSettingsReq struct {
	UserID   int64
	Currency *currency.Token
}

func DefaultPersonalSettingsReq() *PersonalSettingsResp {
	return &PersonalSettingsResp{
		Currency: currency.TokenRUB,
	}
}
