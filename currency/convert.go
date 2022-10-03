package currency

import "context"

type Converter interface {
	Convert(ctx context.Context, req ConvertReq) (amount float64, err error)
}

type ConvertReq struct {
	Amount float64
	From   Token
	To     Token
}
