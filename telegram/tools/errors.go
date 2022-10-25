package tools

import (
	"strings"
)

var (
	ErrInternal = UserError{
		Title:         "unexpected error",
		Help:          "you have found a bug, report it plz",
		IsNotExpected: true,
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
		Title:         "invalid request",
		Help:          "use format [/limit <float>]",
		IsNotExpected: true,
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
	IsNotExpected bool

	internal error
}

func (e UserError) Error() string {
	b := strings.Builder{}

	b.WriteString("Error happened: " + e.Title + "\n")
	b.WriteString("message: " + e.internal.Error() + "\n\n")
	b.WriteString(e.Help)

	return b.String()
}

func (e UserError) With(internal error) UserError {
	e.internal = internal

	return e
}
