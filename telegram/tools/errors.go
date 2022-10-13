package tools

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

var (
	ErrInvalidCreateExpense = TelegramError{
		Title:  "invalid request",
		Detail: "use format [/exp 100 Food]; where 100 is price and number, Food is any string",
	}
	ErrInternal = TelegramError{
		Title:  "internal error",
		Detail: "bot eaten too much fastfood",
	}
	ErrInvalidListExpense = TelegramError{
		Title:  "invalid request",
		Detail: "use format [/all day]; you can use day, month, week, year or 2h30m",
	}

	ErrInvalidSetLimit = TelegramError{
		Title:  "invalid request",
		Detail: "use format [/limit <float>]",
	}

	ErrLimitBlockExpense = TelegramError{
		Title:  "inconsistent request",
		Detail: "cant do operation because of limit",
	}

	ErrSetLimitBlockedByExpenses = TelegramError{
		Title:      "inconsistent request",
		Detail:     "cant set limit less then sum",
		isExpected: true,
	}
)

type TelegramError struct {
	Title      string
	Detail     string
	isExpected bool
}

func (e TelegramError) Error() string {
	b := strings.Builder{}

	b.WriteString("Error happened :(\n")
	b.WriteString(e.Title + ": " + e.Detail)

	return b.String()
}

// SendError default value of returnErr=true
func SendError(c tele.Context, e TelegramError) error {
	_ = c.Send(e.Error())

	if e.isExpected {
		return nil
	}

	return e
}
