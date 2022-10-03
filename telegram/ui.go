package telegram

import (
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
)

var (
	currencySelectorUI, currencyButtonsUI = getCurrencySelector()
)

func getCurrencySelector() (*tele.ReplyMarkup, []tele.Btn) {
	currencySelector := &tele.ReplyMarkup{ResizeKeyboard: true}

	rows := make([]tele.Row, len(currency.Supported))
	buttons := make([]tele.Btn, len(currency.Supported))

	for i := range currency.Supported {
		buttons[i] = currencySelector.Data(
			currency.Token(i).String(),
			"set-currency",
			currency.Token(i).String(),
		)
		rows[i] = currencySelector.Row(buttons[i])
	}

	currencySelector.Inline(rows...)

	return currencySelector, buttons
}
