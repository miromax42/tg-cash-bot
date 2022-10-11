package tools

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

var (
	ErrInvalidCreateExpense = ErrObj{
		Title:  "invalid request",
		Detail: "use format [/exp 100 Food]; where 100 is price and number, Food is any string",
	}
	ErrInternal = ErrObj{
		Title:  "internal error",
		Detail: "bot eaten too much fastfood",
	}
	ErrInvalidListExpense = ErrObj{
		Title:  "invalid request",
		Detail: "use format [/all day]; you can use day, month, week, year or 2h30m",
	}

	ErrInvalidSetLimit = ErrObj{
		Title:  "invalid request",
		Detail: "use format [/limit <float>]",
	}

	ErrLimitBlockExpense = ErrObj{
		Title:  "inconsistent request",
		Detail: "cant do operation because of limit",
	}

	ErrSetLimitBlockedByExpenses = ErrObj{
		Title:  "inconsistent request",
		Detail: "cant set limit less then sum",
	}
)

type ErrObj struct {
	Title  string
	Detail string
}

func (e ErrObj) Error() string {
	b := strings.Builder{}
	b.WriteString("Error happened :(\n")
	b.WriteString(e.Title + ": " + e.Detail)
	return b.String()
}

func SendError(c tele.Context, e ErrObj) {
	_ = c.Send(e.Error())
}
