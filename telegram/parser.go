package telegram

import (
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

var (
	ErrArgsCount  = errors.New("not supported args count")
	ErrValidation = errors.New("not valid")
)

func NewCreateExpenseReq(c tele.Context) (CreateExpenseReq, error) {
	const (
		awaitedArgs = 2

		minAmount      = 0
		maxAmount      = 10000
		minLenCategory = 0
		maxLenCategory = 100
	)

	if len(c.Args()) != awaitedArgs {
		return CreateExpenseReq{}, ErrArgsCount
	}

	amount, err := strconv.ParseFloat(c.Args()[0], 64) //nolint:gomnd
	if err != nil {
		return CreateExpenseReq{}, errors.Wrapf(err, "parse to float %q", c.Args()[0])
	}
	category := c.Args()[1]

	if !(minAmount < amount && amount < maxAmount) {
		return CreateExpenseReq{}, errors.WithHint(ErrValidation, "amount")
	}
	if !(minLenCategory < len(category) && len(category) <= maxLenCategory) {
		return CreateExpenseReq{}, errors.WithHint(ErrValidation, "category")
	}

	return CreateExpenseReq{
		UserID:   c.Sender().ID,
		Amount:   amount,
		Category: category,
	}, nil
}

func NewListUserExpenseReq(c tele.Context) (ListUserExpenseReq, error) {
	const (
		awaitedArgs = 1

		hoursInDay   = 24
		hoursInWeek  = hoursInDay * 7
		hoursInMonth = hoursInWeek * 30
		hoursInYear  = 8760
	)
	if len(c.Args()) != awaitedArgs {
		return ListUserExpenseReq{}, ErrArgsCount
	}

	durationToken := c.Args()[0]

	duration, err := time.ParseDuration(durationToken)
	if err != nil {
		switch durationToken {
		case "day":
			duration = hoursInDay * time.Hour
		case "week":
			duration = hoursInWeek * time.Hour
		case "month":
			duration = hoursInMonth * time.Hour
		case "year":
			duration = hoursInYear * time.Hour
		default:
			return ListUserExpenseReq{}, errors.WithHint(ErrValidation, "duration")
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
		return repo.PersonalSettingsReq{}, errors.Wrap(err, "parse to currency")
	}

	return repo.PersonalSettingsReq{
		UserID:   c.Sender().ID,
		Currency: &newCurrency,
	}, nil
}

func NewSetLimitRequest(c tele.Context) (SetLimitReq, error) {
	limit, err := strconv.ParseFloat(c.Data(), 64) //nolint:gomnd
	err = errors.Wrap(err, "arg must be float")

	return SetLimitReq{
		Limit: limit,
	}, err
}
