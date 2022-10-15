package tools

import (
	"strings"

	"github.com/cockroachdb/errors"
	tele "gopkg.in/telebot.v3"
)

var (
	ErrInternal = UserError{
		Title:         "unexpected error",
		Help:          "you have found a bug, report it plz",
		isNotExpected: true,
	}

	ErrInvalidCreateExpense = UserError{
		Title: "invalid request",
		Help:  "use format [/exp 100 Food]; where 100 is price and number, Food is any string",
	}
	ErrInvalidListExpense = UserError{
		Title: "invalid request",
		Help:  "use format [/all day]; you can use day, month, week, year or 2h30m",
	}
	ErrInvalidSetLimit = UserError{
		Title: "invalid request",
		Help:  "use format [/limit <float>]",
	}
	ErrLimitBlockExpense = UserError{
		Title: "inconsistent request",
		Help:  "cant do operation because of limit",
	}
	ErrSetLimitBlockedByExpenses = UserError{
		Title: "inconsistent request",
		Help:  "cant set limit less then sum",
	}
)

type UserError struct {
	Title         string
	Help          string
	isNotExpected bool

	internal error
}

func (e UserError) Error() string {
	b := strings.Builder{}

	b.WriteString("Error happened: " + e.Title + "\n")
	b.WriteString("message: " + e.internal.Error() + "\n\n")
	b.WriteString(e.Help)

	return b.String()
}

func (e UserError) with(internal error) UserError {
	e.internal = internal

	return e
}

func SendError(err error, c tele.Context, e UserError) error {
	terr := c.Send(e.with(err).Error())
	if terr != nil {
		err = errors.Wrapf(terr, "during handling: %v", err)
	} else if !e.isNotExpected {
		return nil
	}

	return err
}
