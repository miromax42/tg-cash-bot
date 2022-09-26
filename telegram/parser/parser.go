package parser

import (
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

func NewCreateExpenseReq(userID int64, text string) (repo.CreateExpenseReq, error) {
	tokens := strings.Split(text, " ")
	if len(tokens) != 3 {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}

	amount, err := strconv.Atoi(tokens[1])
	if err != nil {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}
	category := tokens[2]

	if !(0 < amount && amount < 10000) {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}
	if !(0 < len(category) && len(category) <= 100) {
		return repo.CreateExpenseReq{}, util.ErrBadFormat
	}

	return repo.CreateExpenseReq{
		UserID:   userID,
		Amount:   amount,
		Category: category,
	}, nil
}

func NewListUserExpenseReq(userID int64, text string) (repo.ListUserExpenseReq, error) {
	tokens := strings.Split(text, " ")
	if len(tokens) != 2 {
		return repo.ListUserExpenseReq{}, util.ErrBadFormat
	}

	duration, err := time.ParseDuration(tokens[1])
	if err != nil {
		switch tokens[1] {
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
		UserID:   userID,
		FromTime: time.Now().Add(-duration),
	}, nil
}
