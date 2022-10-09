package exhange

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type ExchangeSuite struct {
	suite.Suite
	c currency.Converter
}

func (s *ExchangeSuite) SetupSuite() {
	cfg, err := util.NewConfig()
	require.NoError(s.T(), err)

	s.c, err = New(context.Background(), cfg.Exchange)
	require.NoError(s.T(), err)
}

func TestExchangeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExchangeSuite))
}

func (s *ExchangeSuite) TestConverter_Convert() {
	const testCount float64 = 1234

	RUBtoUSD, err := s.c.Convert(context.Background(), currency.ConvertReq{
		Amount: 1,
		From:   currency.TokenRUB,
		To:     currency.TokenUSD,
	})
	require.NoError(s.T(), err)

	RUBtoEUR, err := s.c.Convert(context.Background(), currency.ConvertReq{
		Amount: 1,
		From:   currency.TokenRUB,
		To:     currency.TokenEUR,
	})
	require.NoError(s.T(), err)

	tests := []struct {
		name       string
		arg        currency.ConvertReq
		wantAmount float64
		wantErr    bool
	}{
		{
			"eur to rub",
			currency.ConvertReq{
				Amount: testCount,
				From:   currency.TokenEUR,
				To:     currency.TokenRUB,
			},
			testCount / RUBtoEUR,
			false,
		},
		{
			"eur to usd",
			currency.ConvertReq{
				Amount: testCount,
				From:   currency.TokenEUR,
				To:     currency.TokenUSD,
			},
			testCount / RUBtoEUR * RUBtoUSD,
			false,
		},
		{
			"rub to usd",
			currency.ConvertReq{
				Amount: testCount,
				From:   currency.TokenRUB,
				To:     currency.TokenUSD,
			},
			testCount * RUBtoUSD,
			false,
		},
		{
			"rub to rub",
			currency.ConvertReq{
				Amount: testCount,
				From:   currency.TokenRUB,
				To:     currency.TokenRUB,
			},
			testCount,
			false,
		},
		{
			"usd to usd",
			currency.ConvertReq{
				Amount: testCount,
				From:   currency.TokenUSD,
				To:     currency.TokenUSD,
			},
			testCount,
			false,
		},
		{
			"unknown to usd",
			currency.ConvertReq{
				Amount: testCount,
				From:   currency.Token(1377),
				To:     currency.TokenUSD,
			},
			0,
			true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			gotAmount, err := s.c.Convert(context.Background(), tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAmount != tt.wantAmount {
				t.Errorf("Convert() gotAmount = %v, want %v", gotAmount, tt.wantAmount)
			}
		})
	}
}
