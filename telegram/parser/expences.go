package parser

import (
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

func NewCreateExpenseReq(c tele.Context) (repo.CreateExpenseReq, error) {
	if len(c.Args()) != 2 {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}

	amount, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}
	category := c.Args()[1]

	if !(0 < amount && amount < 10000) {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}
	if !(0 < len(category) && len(category) <= 100) {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}

	return repo.CreateExpenseReq{
		UserID:   c.Sender().ID,
		Amount:   amount,
		Category: category,
	}, nil
}

func NewListUserExpenseReq(c tele.Context) (repo.ListUserExpenseReq, error) {
	if len(c.Args()) != 1 {
		return repo.ListUserExpenseReq{}, util.ErrBadFormat
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
			return repo.ListUserExpenseReq{}, util.ErrBadFormat
		}
	}

	return repo.ListUserExpenseReq{
		UserID:   c.Sender().ID,
		FromTime: time.Now().Add(-duration),
	}, nil
}
