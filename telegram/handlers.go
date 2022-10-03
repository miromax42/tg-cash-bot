package telegram

import (
	"context"
	"strings"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/errors"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/parser"
)

func (s *Server) StartHelp(c tele.Context) error {
	b := strings.Builder{}

	b.WriteString("availible commands:\n")
	b.WriteString("/exp 99 Fun -- adds expense 99 of Fun category\n")
	b.WriteString("/all day -- show expenses for last day,\n\t examples of time modificators [day, week, month, year, 2h30m]\n")

	return c.Send(b.String())
}

func (s *Server) CreateExpense(c tele.Context) error {
	req, err := parser.NewCreateExpenseReq(c)
	if err != nil {
		errors.SendError(c, errors.ErrInvalidCreateExpense)
		return err
	}

	resp, err := s.expense.CreateExpense(context.TODO(), req)
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	return c.Send(resp.String())
}

func (s *Server) ListExpenses(c tele.Context) error {
	req, err := parser.NewListUserExpenseReq(c)
	if err != nil {
		errors.SendError(c, errors.ErrInvalidListExpense)
		return err
	}

	resp, err := s.expense.ListUserExpense(context.TODO(), req)
	if err != nil {
		errors.SendError(c, errors.ErrInternal)
		return err
	}

	return c.Send(resp.String())
}

func (s *Server) SelectCurrency(c tele.Context) error {
	return c.Send("Chose currency:", currencySelectorUI)
}

func (s *Server) SetCurrency(c tele.Context) error {
	defer func() {
		_ = c.Respond()
	}()

	req, err := parser.NewPersonalSettingsReq(c)
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
