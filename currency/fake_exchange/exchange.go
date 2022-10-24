package fake_exchange

import (
	"context"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
)

type Exchange struct{}

func (e Exchange) Convert(_ context.Context, req currency.ConvertReq) (amount float64, err error) {
	return req.Amount, nil
}

func (e Exchange) Base() currency.Token {
	return currency.TokenRUB
}
