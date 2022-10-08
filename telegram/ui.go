package telegram

import (
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
)

func getCurrencySelector() (*tele.ReplyMarkup, *tele.Btn) {
	currencySelector := &tele.ReplyMarkup{ResizeKeyboard: true}
	rows := make([]tele.Row, len(currency.Supported))

	var button tele.Btn
	for i := range currency.Supported {
		button = currencySelector.Data(
			currency.Token(i).String(),
			"set-currency",
			currency.Token(i).String(),
		)
		rows[i] = currencySelector.Row(button)
	}

	currencySelector.Inline(rows...)

	return currencySelector, &button
}
