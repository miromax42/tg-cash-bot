package telegram

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/tools"
)

const oneCoin = 1

func StartHelp(c tele.Context) error {
	b := strings.Builder{}

	b.WriteString("available commands:\n")
	b.WriteString("/exp 99 Fun -- adds expense 99 of Fun category\n")
	b.WriteString("/all day -- show expenses for last day,\n")
	b.WriteString("\t examples of time modificators [day, week, month, year, 2h30m]\n")

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
		return s.SendError(err, c, tools.ErrInvalidCreateExpense)
	}

	amount, err := s.exchange.Convert(requestContext(c), currency.ConvertReq{
		Amount: req.Amount,
		From:   c.Get(SettingsKey.String()).(*repo.PersonalSettingsResp).Currency,
		To:     s.exchange.Base(),
	})
	if err != nil {
		return s.SendError(err, c, tools.ErrInternal)
	}

	databaseReq := repo.CreateExpenseReq{
		UserID:   req.UserID,
		Amount:   amount,
		Category: req.Category,
	}

	resp, err := s.expense.CreateExpense(requestContext(c), databaseReq)
	if err != nil {
		if errors.Is(err, repo.ErrLimitExceed) {
			return s.SendError(err, c, tools.ErrLimitBlockExpense)
		}

		return s.SendError(err, c, tools.ErrInternal)
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
		return s.SendError(err, c, tools.ErrInvalidListExpense)
	}

	databaseReq := repo.ListUserExpenseReq{
		UserID:   req.UserID,
		FromTime: req.FromTime,
	}

	resp, err := s.expense.ListUserExpense(requestContext(c), databaseReq)
	if err != nil {
		return s.SendError(err, c, tools.ErrInternal)
	}

	multiplier, err := s.exchange.Convert(requestContext(c), currency.ConvertReq{
		Amount: oneCoin,
		From:   s.exchange.Base(),
		To:     c.Get(SettingsKey.String()).(*repo.PersonalSettingsResp).Currency,
	})
	if err != nil {
		return s.SendError(err, c, tools.ErrInternal)
	}

	return c.Send(ListExpensesAnswer(resp, multiplier))
}

func (s *Server) SelectCurrency(reply *tele.ReplyMarkup) func(c tele.Context) error {
	return func(c tele.Context) error {
		return c.Send("Chose currency:", reply)
	}
}

func (s *Server) SetCurrency(c tele.Context) error {
	defer func() {
		_ = c.Respond()
	}()

	req, err := NewPersonalSettingsReq(c)
	if err != nil {
		return s.SendError(err, c, tools.ErrInternal)
	}

	if err := s.userSettings.Set(requestContext(c), req); err != nil {
		return s.SendError(err, c, tools.ErrInternal)
	}

	return c.Send("currency set to " + c.Data())
}

type SetLimitReq struct {
	Limit float64
}

func (s *Server) SetLimit(c tele.Context) error {
	req, err := NewSetLimitRequest(c)
	if err != nil {
		return s.SendError(err, c, tools.ErrInvalidSetLimit)
	}

	amount, err := s.exchange.Convert(requestContext(c), currency.ConvertReq{
		Amount: req.Limit,
		From:   c.Get(SettingsKey.String()).(*repo.PersonalSettingsResp).Currency,
		To:     s.exchange.Base(),
	})
	if err != nil {
		return s.SendError(err, c, tools.ErrInvalidSetLimit)
	}

	repoReq := repo.PersonalSettingsReq{
		UserID: c.Sender().ID,
		Limit:  &amount,
	}

	if err = s.userSettings.Set(requestContext(c), repoReq); err != nil {
		if errors.Is(err, repo.ErrLimitExceed) {
			return s.SendError(err, c, tools.ErrSetLimitBlockedByExpenses)
		}

		return s.SendError(err, c, tools.ErrInternal)
	}

	return c.Send("limit set to " + fmt.Sprintf("%.2f", req.Limit))
}
