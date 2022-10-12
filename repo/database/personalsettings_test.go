package database

import (
	"context"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

func (s *PostgresTestSuite) TestGet() {
	var (
		settings = NewPersonalSettings(s.c)

		idNotExists = s.generatorID()
		idExist     = s.generatorID()

		newCurrency = currency.MustParse(currency.TokenCNY.String())
	)

	s.applyFixture(
		"fixtures_test/test_personal_settings_get.yml",
		map[string]interface{}{
			"idExist":     idExist,
			"currencyCNY": newCurrency.String(),
		},
	)

	tests := []struct {
		name    string
		arg     int64
		want    *repo.PersonalSettingsResp
		wantErr bool
	}{
		{
			"default",
			idNotExists,
			repo.DefaultPersonalSettingsResp(),
			false,
		},
		{
			"id with CNY",
			idExist,
			&repo.PersonalSettingsResp{Currency: newCurrency},
			false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := settings.Get(context.Background(), tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *PostgresTestSuite) TestSet() {
	var (
		settings = NewPersonalSettings(s.c)

		CurrencyCNY = currency.MustParse(currency.TokenCNY.String())
		CurrencyEUR = currency.MustParse(currency.TokenEUR.String())
	)

	createWithCNY := func() int64 {
		newID := s.generatorID()

		err := settings.Set(context.Background(), repo.PersonalSettingsReq{
			UserID:   newID,
			Currency: &CurrencyCNY,
		})
		require.NoError(s.T(), err)

		return newID
	}

	get := func(id int64) *repo.PersonalSettingsResp {
		resp, err := settings.Get(context.Background(), id)
		require.NoError(s.T(), err)

		return resp
	}

	tests := []struct {
		name    string
		prepare func() repo.PersonalSettingsReq
		want    *repo.PersonalSettingsResp
		wantErr bool
	}{
		{
			"change existing",
			func() repo.PersonalSettingsReq {
				id := createWithCNY()

				return repo.PersonalSettingsReq{
					UserID:   id,
					Currency: &CurrencyEUR,
				}
			},
			&repo.PersonalSettingsResp{Currency: currency.TokenEUR},
			false,
		},
		{
			"set existing without currency",
			func() repo.PersonalSettingsReq {
				id := createWithCNY()

				return repo.PersonalSettingsReq{
					UserID:   id,
					Currency: nil,
				}
			},
			&repo.PersonalSettingsResp{
				Currency: currency.TokenCNY,
			},
			false,
		},
		{
			"set new",
			func() repo.PersonalSettingsReq {
				return repo.PersonalSettingsReq{
					UserID:   s.generatorID(),
					Currency: &CurrencyEUR,
				}
			},
			&repo.PersonalSettingsResp{Currency: currency.TokenEUR},
			false,
		},
		{
			"set new with default currency",
			func() repo.PersonalSettingsReq {
				return repo.PersonalSettingsReq{
					UserID:   s.generatorID(),
					Currency: nil,
				}
			},
			&repo.PersonalSettingsResp{Currency: currency.TokenRUB},
			false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			arg := tt.prepare()
			err := settings.Set(context.Background(), arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(get(arg.UserID), tt.want) {
				t.Errorf("Get() got = %v, want %v", get(arg.UserID), tt.want)
			}
		})
	}
}
