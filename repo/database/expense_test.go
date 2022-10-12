package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

func (s *PostgresTestSuite) TestCreateExpense() {
	const (
		amount   float64 = 20
		category         = "test"
	)

	var (
		expenceInstance = NewExpense(s.c)

		idWithLimit1      = s.generatorID()
		idWithLimit1000   = s.generatorID()
		idWithLimitNotSet = s.generatorID()
	)

	tests := []struct {
		name    string
		req     repo.CreateExpenseReq
		want    *repo.CreateExpenseResp
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"user with no limit",
			repo.CreateExpenseReq{
				UserID:   idWithLimitNotSet,
				Amount:   amount,
				Category: category,
			},
			&repo.CreateExpenseResp{
				Amount:    amount,
				Category:  category,
				CreatedAt: time.Now(),
			},
			assert.NoError,
		},
		{
			"user with limit 1000",
			repo.CreateExpenseReq{
				UserID:   idWithLimit1000,
				Amount:   amount,
				Category: category,
			},
			&repo.CreateExpenseResp{
				Amount:    amount,
				Category:  category,
				CreatedAt: time.Now(),
			},
			assert.NoError,
		},
		{
			"user with limit 1",
			repo.CreateExpenseReq{
				UserID:   idWithLimit1,
				Amount:   amount,
				Category: category,
			},
			&repo.CreateExpenseResp{
				Amount:    amount,
				Category:  category,
				CreatedAt: time.Now(),
			},
			assert.Error,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.applyFixture(
				"fixtures_test/test_expeses_create.yml",
				map[string]interface{}{
					"idWithLimitNotSet": idWithLimitNotSet,
					"idWithLimit1":      idWithLimit1,
					"idWithLimit1000":   idWithLimit1000,
				},
			)

			resp, err := expenceInstance.CreateExpense(context.Background(), tt.req)
			tt.wantErr(t, err)
			if err != nil {
				return
			}

			require.NotNil(t, resp)
			assert.Equal(t, tt.want.Amount, resp.Amount)
			assert.Equal(t, tt.want.Category, resp.Category)
			assert.WithinDuration(t, tt.want.CreatedAt, resp.CreatedAt, 10*time.Second)
		})
	}
}
