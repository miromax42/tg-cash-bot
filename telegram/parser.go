package telegram

import (
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

func NewCreateExpenseReq(c tele.Context) (CreateExpenseReq, error) {
	if len(c.Args()) != 2 {
		return CreateExpenseReq{}, util.ErrBadFormat
	}

	amount, err := strconv.ParseFloat(c.Args()[0], 64)
	if err != nil {
		return CreateExpenseReq{}, util.ErrBadFormat
	}
	category := c.Args()[1]

	if !(0 < amount && amount < 10000) {
		return CreateExpenseReq{}, util.ErrBadFormat
	}
	if !(0 < len(category) && len(category) <= 100) {
		return CreateExpenseReq{}, util.ErrBadFormat
	}

	return CreateExpenseReq{
		UserID:   c.Sender().ID,
		Amount:   amount,
		Category: category,
	}, nil
}

func NewListUserExpenseReq(c tele.Context) (ListUserExpenseReq, error) {
	if len(c.Args()) != 1 {
		return ListUserExpenseReq{}, util.ErrBadFormat
	}

	durationToken := c.Args()[0]

	duration, err := time.ParseDuration(durationToken)
	if err != nil {
		switch durationToken {
		case "day":
			duration = 24 * time.Hour
		case "week":
			duration = 7 * 24 * time.Hour
		case "month":
			duration = 30 * 24 * time.Hour
		case "year":
			duration = 8760 * time.Hour
		default:
			return ListUserExpenseReq{}, util.ErrBadFormat
		}
	}

	return ListUserExpenseReq{
		UserID:   c.Sender().ID,
		FromTime: time.Now().Add(-duration),
	}, nil
}

func NewPersonalSettingsReq(c tele.Context) (repo.PersonalSettingsReq, error) {
	newCurrency, err := currency.Parse(c.Data())
	if err != nil {
		return repo.PersonalSettingsReq{}, err
	}

	return repo.PersonalSettingsReq{
		UserID:   c.Sender().ID,
		Currency: &newCurrency,
	}, nil
}

func NewSetLimitRequest(c tele.Context) (SetLimitReq, error) {
	limit, err := strconv.ParseFloat(c.Data(), 64)
	return SetLimitReq{
		Limit: limit,
	}, err
}
