package telegram

import (
	"context"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/errors"
)

const oneCoin = 1

func (s *Server) StartHelp(c tele.Context) error {
	b := strings.Builder{}

	b.WriteString("availible commands:\n")
	b.WriteString("/exp 99 Fun -- adds expense 99 of Fun category\n")
	b.WriteString("/all day -- show expenses for last day,\n\t examples of time modificators [day, week, month, year, 2h30m]\n")

	return c.Send(b.String())
}

type CreateExpenseReq struct {
	UserID   int64
	Amount   float64
	Category string
}

func (s *Server) CreateExpense(c tele.Context) error {
	req, err := NewCreateExpenseReq(c)
	if err != nil {
		errors.SendError(c, errors.ErrInvalidCreateExpense)
		return err
	}

	settings, err := s.userSettings.Get(context.TODO(), req.UserID)
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	amount, err := s.exchange.Convert(context.TODO(), currency.ConvertReq{
		Amount: req.Amount,
		From:   settings.Currency,
		To:     s.exchange.Base(),
	})

	databaseReq := repo.CreateExpenseReq{
		UserID:   req.UserID,
		Amount:   amount,
		Category: req.Category,
	}

	resp, err := s.expense.CreateExpense(context.TODO(), databaseReq)
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	return c.Send(CreateExpenseAnswer(resp, req.Amount))
}

type ListUserExpenseReq struct {
	UserID   int64
	FromTime time.Time
}

func (s *Server) ListExpenses(c tele.Context) error {
	req, err := NewListUserExpenseReq(c)
	if err != nil {
		errors.SendError(c, errors.ErrInvalidListExpense)
		return err
	}

	settings, err := s.userSettings.Get(context.TODO(), req.UserID)
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	databaseReq := repo.ListUserExpenseReq{
		UserID:   req.UserID,
		FromTime: req.FromTime,
	}

	resp, err := s.expense.ListUserExpense(context.TODO(), databaseReq)
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	multiplier, err := s.exchange.Convert(context.TODO(), currency.ConvertReq{
		Amount: oneCoin,
		From:   s.exchange.Base(),
		To:     settings.Currency,
	})
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	return c.Send(ListExpensesAnswer(resp, multiplier))
}

func (s *Server) SelectCurrency(c tele.Context) error {
	return c.Send("Chose currency:", currencySelectorUI)
}

func (s *Server) SetCurrency(c tele.Context) error {
	defer func() {
		_ = c.Respond()
	}()

	req, err := NewPersonalSettingsReq(c)
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	if err := s.userSettings.Set(context.TODO(), req); err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	return c.Send("currency set to " + c.Data())
}
