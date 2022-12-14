package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/pb"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/cache"
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

	amount, err := s.exchange.Convert(RequestContext(c), currency.ConvertReq{
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

	resp, err := s.expense.CreateExpense(RequestContext(c), databaseReq)
	if err != nil {
		if errors.Is(err, repo.ErrLimitExceed) {
			return s.SendError(err, c, tools.ErrLimitBlockExpense)
		}

		return s.SendError(err, c, tools.ErrInternal)
	}

	return s.Send(c, CreateExpenseAnswer(resp, req.Amount))
}

type ListUserExpenseReq struct {
	UserID   int64
	FromTime time.Time
	ToTime   *time.Time
}

func (s *Server) ListExpenses(c tele.Context) error {
	req, err := NewListUserExpenseReq(c)
	if err != nil {
		return s.SendError(err, c, tools.ErrInvalidListExpense)
	}

	databaseReq := req.ToDB()

	var resp repo.ListUserExpenseResp
	err = s.dbCache.Get(RequestContext(c), cache.UserReportToken(databaseReq), &resp)
	if err != nil {
		if errors.Is(err, cache.ErrMiss) {
			resp, err = s.expense.ListUserExpense(RequestContext(c), databaseReq)
			if err != nil {
				return s.SendError(err, c, tools.ErrInternal)
			}

			if err = s.dbCache.Set(RequestContext(c), cache.UserReportToken(databaseReq), resp); err != nil {
				return s.SendError(err, c, tools.ErrInternal)
			}
		} else {
			_ = s.SendError(err, c, tools.ErrInternal)

			return errors.Wrapf(err, "cache")
		}
	}

	multiplier, err := s.exchange.Convert(RequestContext(c), currency.ConvertReq{
		Amount: oneCoin,
		From:   s.exchange.Base(),
		To:     c.Get(SettingsKey.String()).(*repo.PersonalSettingsResp).Currency,
	})
	if err != nil {
		return s.SendError(err, c, tools.ErrInternal)
	}

	if err := s.reportSender.ReportSend(RequestContext(c), &pb.ReportRequest{
		UserId:     c.Sender().ID,
		StartTime:  timestamppb.New(req.FromTime),
		EndTime:    timestamppb.Now(),
		Multiplier: multiplier,
	}); err != nil {
		return errors.Wrapf(err, "report send")
	}

	return nil
}

func (s *Server) SelectCurrency(reply *tele.ReplyMarkup) func(c tele.Context) error {
	return func(c tele.Context) error {
		return s.Send(c, "Chose currency:", reply)
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

	if err := s.userSettings.Set(RequestContext(c), req); err != nil {
		return s.SendError(err, c, tools.ErrInternal)
	}

	return s.Send(c, "currency set to "+c.Data())
}

type SetLimitReq struct {
	Limit float64
}

func (s *Server) SetLimit(c tele.Context) error {
	req, err := NewSetLimitRequest(c)
	if err != nil {
		return s.SendError(err, c, tools.ErrInvalidSetLimit)
	}

	amount, err := s.exchange.Convert(RequestContext(c), currency.ConvertReq{
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

	if err = s.userSettings.Set(RequestContext(c), repoReq); err != nil {
		if errors.Is(err, repo.ErrLimitExceed) {
			return s.SendError(err, c, tools.ErrSetLimitBlockedByExpenses)
		}

		return s.SendError(err, c, tools.ErrInternal)
	}

	return s.Send(c, "limit set to "+fmt.Sprintf("%.2f", req.Limit))
}
