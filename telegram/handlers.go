package telegram

import (
	"context"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

func (s Server) CreateExpense(c tele.Context) error {
	req, err := repo.NewCreateExpenseReq(c.Sender().ID, c.Text())
	if err != nil {
		return err
	}

	resp, err := s.expense.CreateExpense(context.TODO(), req)
	if err != nil {
		return err
	}

	return c.Send(resp.String())
}

func (s Server) ListExpenses(c tele.Context) error {
	req, err := repo.NewListUserExpenseReq(c.Sender().ID, c.Text())
	if err != nil {
		return err
	}

	resp, err := s.expense.ListUserExpense(context.TODO(), req)
	if err != nil {
		return err
	}

	return c.Send(resp.String())
}
