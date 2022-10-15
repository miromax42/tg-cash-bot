package tools

import (
	"strings"

	"github.com/cockroachdb/errors"
	tele "gopkg.in/telebot.v3"
)

var (
	ErrInvalidCreateExpense = UserError{
		Title:  "invalid request",
		Detail: "use format [/exp 100 Food]; where 100 is price and number, Food is any string",
	}
	ErrInternal = UserError{
		Title:  "internal error",
		Detail: "bot eaten too much fastfood",
	}
	ErrInvalidListExpense = UserError{
		Title:  "invalid request",
		Detail: "use format [/all day]; you can use day, month, week, year or 2h30m",
	}

	ErrInvalidSetLimit = UserError{
		Title:  "invalid request",
		Detail: "use format [/limit <float>]",
	}

	ErrLimitBlockExpense = UserError{
		Title:  "inconsistent request",
		Detail: "cant do operation because of limit",
	}

	ErrSetLimitBlockedByExpenses = UserError{
		Title:      "inconsistent request",
		Detail:     "cant set limit less then sum",
		isExpected: true,
	}
)

type UserError struct {
	Title      string
	Detail     string
	isExpected bool
}

func (e UserError) Error() string {
	b := strings.Builder{}

	b.WriteString("Error happened :(\n")
	b.WriteString(e.Title + ": " + e.Detail)

	return b.String()
}

func SendError(err error, c tele.Context, e UserError) error {
	terr := c.Send(e.Error())
	if terr != nil {
		err = errors.Wrapf(terr, "during handling: %v", err)
	} else if e.isExpected {
		return nil
	}

	return err
}
