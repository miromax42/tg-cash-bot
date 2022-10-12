package database

import (
	"context"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent/expense"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Expense struct {
	db *ent.Client
}

func NewExpense(client *ent.Client) *Expense {
	return &Expense{
		db: client,
	}
}

func (e Expense) CreateExpense(
	ctx context.Context,
	req repo.CreateExpenseReq,
) (*repo.CreateExpenseResp, error) {
	var model *ent.Expense
	if err := WithTx(ctx, e.db, func(tx *ent.Tx) error {
		settings := NewPersonalSettings(e.db)
		userSettings, err := settings.Get(ctx, req.UserID)
		if err != nil {
			return err
		}

		sum, err := e.allUserExpense(ctx, repo.ListUserExpenseReq{
			UserID:   req.UserID,
			FromTime: util.TimeMonthAgo(),
		})
		if err != nil {
			return err
		}

		if sum+req.Amount > userSettings.Limit && userSettings.Limit != 0 {
			return util.ErrLimitExceed
		}

		model, err = e.db.Expense.Create().
			SetAmount(req.Amount).
			SetCreatedBy(req.UserID).
			SetCategory(req.Category).
			Save(ctx)

		return err
	}); err != nil {
		return nil, err
	}

	return &repo.CreateExpenseResp{
		Amount:    model.Amount,
		Category:  model.Category,
		CreatedAt: model.CreateTime,
	}, nil
}

func (e Expense) ListUserExpense(ctx context.Context, req repo.ListUserExpenseReq) (repo.ListUserExpenseResp, error) {
	var expenses repo.ListUserExpenseResp
	if err := e.db.Expense.Query().
		Where(expense.CreatedBy(req.UserID)).
		Where(expense.CreateTimeGTE(req.FromTime)).
		GroupBy(expense.FieldCategory).
		Aggregate(ent.Sum(expense.FieldAmount)).
		Scan(ctx, &expenses); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (e Expense) allUserExpense(ctx context.Context, req repo.ListUserExpenseReq) (float64, error) {
	var result []struct {
		Sum       float64 `json:"sum"`
		CreatedBy string  `json:"created_by"`
	}

	if err := e.db.Expense.Query().
		Select(expense.FieldAmount).
		Where(expense.CreatedBy(req.UserID)).
		Where(expense.CreateTimeGTE(req.FromTime)).
		GroupBy(expense.FieldCreatedBy).
		Aggregate(ent.Sum(expense.FieldAmount)).
		Scan(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Sum, nil
}
