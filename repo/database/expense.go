package database

import (
	"context"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent/expense"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
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
	model, err := e.db.Expense.Create().
		SetAmount(req.Amount).
		SetCreatedBy(req.UserID).
		SetCategory(req.Category).
		Save(ctx)
	if err != nil {
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
