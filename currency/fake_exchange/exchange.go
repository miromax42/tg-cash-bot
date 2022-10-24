package fake_exchange

import (
	"context"

	"go.opentelemetry.io/otel"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Exchange struct{}

func (e Exchange) Convert(ctx context.Context, req currency.ConvertReq) (amount float64, err error) {
	_, span := otel.Tracer(util.RequestTrace).Start(ctx, "Exchange.Convert")
	defer span.End()

	return req.Amount, nil
}

func (e Exchange) Base() currency.Token {
	return currency.TokenRUB
}
