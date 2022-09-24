package repo

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Expense interface {
	CreateExpense(ctx context.Context, req CreateExpenseReq) (*CreateExpenseResp, error)
	ListUserExpense(ctx context.Context, req ListUserExpenseReq) (ListUserExpenseResp, error)
}

type CreateExpenseReq struct {
	UserID   int64
	Amount   int
	Category string
}

func NewCreateExpenseReq(userID int64, text string) (CreateExpenseReq, error) {
	tokens := strings.Split(text, " ")
	if len(tokens) != 3 {
		return CreateExpenseReq{}, util.ErrBadFormat
	}

	amount, err := strconv.Atoi(tokens[1])
	if err != nil {
		return CreateExpenseReq{}, util.ErrBadFormat
	}
	category := tokens[2]

	return CreateExpenseReq{
		UserID:   userID,
		Amount:   amount,
		Category: category,
	}, nil
}

type CreateExpenseResp struct {
	Amount    int
	Category  string
	CreatedAt time.Time
}

func (r *CreateExpenseResp) String() string {
	b := strings.Builder{}

	b.WriteString(strconv.Itoa(r.Amount))
	b.WriteString(" on ")
	b.WriteString(r.Category + "\n")
	b.WriteString("(" + r.CreatedAt.Format(time.Kitchen) + ")")

	return b.String()
}

type ListUserExpenseReq struct {
	UserID   int64
	FromTime time.Time
}

func NewListUserExpenseReq(userID int64, text string) (ListUserExpenseReq, error) {
	tokens := strings.Split(text, " ")
	if len(tokens) != 2 {
		return ListUserExpenseReq{}, util.ErrBadFormat
	}

	duration, err := time.ParseDuration(tokens[1])
	if err != nil {
		return ListUserExpenseReq{}, err
	}

	return ListUserExpenseReq{
		UserID:   userID,
		FromTime: time.Now().Add(-duration),
	}, nil
}

type ListUserExpenseResp []struct {
	Category string `json:"category"`
	Sum      int    `json:"sum"`
}

func (r ListUserExpenseResp) String() string {
	b := strings.Builder{}

	for i := range r {
		b.WriteString(r[i].Category)
		b.WriteString(": ")
		b.WriteString(strconv.Itoa(r[i].Sum))

		if i != len(r)-1 {
			b.WriteString("\n")
		}
	}

	if b.Len() == 0 {
		return "no expenses"
	}

	return b.String()
}
