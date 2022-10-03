package parser

import (
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
)

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
